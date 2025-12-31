package svc

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/cy77cc/go-microstack/fileserver/model"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/config"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage interface {
	CreateBucket(ctx context.Context, bucket string) error
	PutObject(ctx context.Context, bucket, objectName string, data []byte, Hash string, contentType string) (etag string, err error)
	GetObject(ctx context.Context, bucket, objectName string) (payload []byte, contentType string, err error)
	InitiateMultipart(ctx context.Context, bucket, objectName string, contentType string) (uploadID string, err error)
	UploadPart(ctx context.Context, bucket, objectName, uploadID string, partNumber int, data []byte) (etag string, err error)
	CompleteMultipart(ctx context.Context, bucket, objectName, uploadID string, parts []CompletedPart) (etag string, err error)
	AbortMultipart(ctx context.Context, bucket, objectName, uploadID string) error
	Presign(ctx context.Context, bucket, objectName string, expire int64) (string, error)
}

type CompletedPart struct {
	PartNumber int
	ETag       string
}

type StorageRouter interface {
	Select(ctx context.Context, bucket string) (Storage, error)
}

type storageRouter struct {
	c         config.Config
	local     *localStorage
	MinioCore *minio.Core
	MinioCli  *minio.Client
	bm        model.BucketConfigModel
}

func NewStorageRouter(c config.Config, baseDir string) (StorageRouter, error) {
	var core *minio.Core
	var cli *minio.Client
	var err error
	if c.Minio.Endpoint != "" {
		core, err = minio.NewCore(c.Minio.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(c.Minio.Username, c.Minio.Password, ""),
			Secure: c.Minio.SSL,
		})
		if err != nil {
			return nil, err
		}
		cli, err = minio.New(c.Minio.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(c.Minio.Username, c.Minio.Password, ""),
			Secure: c.Minio.SSL,
		})
		if err != nil {
			return nil, err
		}
	}
	return &storageRouter{
		c:         c,
		local:     &localStorage{baseDir: baseDir},
		MinioCore: core,
		MinioCli:  cli,
	}, nil
}

func (r *storageRouter) Select(ctx context.Context, bucket string) (Storage, error) {
	if r.bm == nil {
		return r.local, nil
	}
	cfg, err := r.bm.FindOneByBucket(ctx, bucket)
	if err != nil {
		// default to local
		return r.local, nil
	}
	switch cfg.StorageType {
	case 1:
		if r.MinioCore == nil {
			return nil, fmt.Errorf("minio core not initialized")
		}
		return &minioStorage{core: r.MinioCore, cli: r.MinioCli}, nil
	default:
		return r.local, nil
	}
}

type localStorage struct {
	baseDir string
}

func (l *localStorage) ensureBucket(bucket string) error {
	p := filepath.Join(l.baseDir, bucket)
	return os.MkdirAll(p, 0o755)
}

func (l *localStorage) CreateBucket(ctx context.Context, bucket string) error {
	return l.ensureBucket(bucket)
}

func (l *localStorage) PutObject(ctx context.Context, bucket, objectName string, data []byte, hash string, contentType string) (string, error) {
	if err := l.ensureBucket(bucket); err != nil {
		return "", err
	}
	path := filepath.Join(l.baseDir, bucket, objectName)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:]), nil
}

func (l *localStorage) GetObject(ctx context.Context, bucket, objectName string) ([]byte, string, error) {
	path := filepath.Join(l.baseDir, bucket, objectName)
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	return b, "", nil
}

func (l *localStorage) InitiateMultipart(ctx context.Context, bucket, objectName string, contentType string) (string, error) {
	if err := l.ensureBucket(bucket); err != nil {
		return "", err
	}
	uploadID := uuid.NewString()
	dir := filepath.Join(l.baseDir, bucket, ".multipart", uploadID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return uploadID, nil
}

func (l *localStorage) UploadPart(ctx context.Context, bucket, objectName, uploadID string, partNumber int, data []byte) (string, error) {
	dir := filepath.Join(l.baseDir, bucket, ".multipart", uploadID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	partPath := filepath.Join(dir, fmt.Sprintf("%06d.part", partNumber))
	if err := os.WriteFile(partPath, data, 0o644); err != nil {
		return "", err
	}
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:]), nil
}

func (l *localStorage) CompleteMultipart(ctx context.Context, bucket, objectName, uploadID string, parts []CompletedPart) (string, error) {
	dir := filepath.Join(l.baseDir, bucket, ".multipart", uploadID)
	sort.Slice(parts, func(i, j int) bool { return parts[i].PartNumber < parts[j].PartNumber })
	var buf bytes.Buffer
	for _, item := range parts {
		partPath := filepath.Join(dir, fmt.Sprintf("%06s.part", fmtPartNumber(item.PartNumber)))
		b, err := os.ReadFile(partPath)
		if err != nil {
			return "", err
		}
		if _, err = buf.Write(b); err != nil {
			return "", err
		}
	}
	if err := l.PutCombined(bucket, objectName, buf.Bytes()); err != nil {
		return "", err
	}
	sum := md5.Sum(buf.Bytes())
	_ = os.RemoveAll(dir)
	return hex.EncodeToString(sum[:]), nil
}

func fmtPartNumber(n int) string {
	s := strconv.Itoa(n)
	for len(s) < 6 {
		s = "0" + s
	}
	return s
}

func (l *localStorage) PutCombined(bucket, objectName string, data []byte) error {
	path := filepath.Join(l.baseDir, bucket, objectName)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (l *localStorage) AbortMultipart(ctx context.Context, bucket, objectName, uploadID string) error {
	dir := filepath.Join(l.baseDir, bucket, ".multipart", uploadID)
	return os.RemoveAll(dir)
}

func (l *localStorage) Presign(ctx context.Context, bucket, objectName string, expire int64) (string, error) {
	// Local storage returns relative path or direct link pattern
	// API layer should prepend host
	return fmt.Sprintf("/%s/%s", bucket, objectName), nil
}

type minioStorage struct {
	core *minio.Core
	cli  *minio.Client
}

func (m *minioStorage) CreateBucket(ctx context.Context, bucket string) error {
	exists, err := m.core.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return m.core.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}

func (m *minioStorage) PutObject(ctx context.Context, bucket, objectName string, data []byte, hash string, contentType string) (string, error) {
	r := bytes.NewReader(data)
	info, err := m.cli.PutObject(ctx, bucket, objectName, r, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	return info.ETag, nil
}

func (m *minioStorage) GetObject(ctx context.Context, bucket, objectName string) ([]byte, string, error) {
	obj, err := m.cli.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", err
	}
	defer obj.Close()
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, obj); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "", nil
}

func (m *minioStorage) InitiateMultipart(ctx context.Context, bucket, objectName string, contentType string) (string, error) {
	UploadID, err := m.core.NewMultipartUpload(ctx, bucket, objectName, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return UploadID, nil
}

func (m *minioStorage) UploadPart(ctx context.Context, bucket, objectName, uploadID string, partNumber int, data []byte) (string, error) {
	reader := bytes.NewReader(data)
	objPart, err := m.core.PutObjectPart(ctx, bucket, objectName, uploadID, partNumber, reader, int64(len(data)), minio.PutObjectPartOptions{})
	if err != nil {
		return "", err
	}
	return objPart.ETag, nil
}

func (m *minioStorage) CompleteMultipart(ctx context.Context, bucket, objectName, uploadID string, parts []CompletedPart) (string, error) {
	var objParts []minio.CompletePart
	for _, p := range parts {
		objParts = append(objParts, minio.CompletePart{ETag: p.ETag, PartNumber: p.PartNumber})
	}
	info, err := m.core.CompleteMultipartUpload(ctx, bucket, objectName, uploadID, objParts, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return info.ETag, nil
}

func (m *minioStorage) AbortMultipart(ctx context.Context, bucket, objectName, uploadID string) error {
	return m.core.AbortMultipartUpload(ctx, bucket, objectName, uploadID)
}

func (m *minioStorage) Presign(ctx context.Context, bucket, objectName string, expire int64) (string, error) {
	if expire <= 0 {
		expire = 3600
	}
	url, err := m.cli.PresignedGetObject(ctx, bucket, objectName, time.Duration(expire)*time.Second, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
