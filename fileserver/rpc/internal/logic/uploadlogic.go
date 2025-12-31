package logic

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/cy77cc/go-microstack/fileserver/model"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"
	"github.com/google/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadLogic) Upload(in *pb.UploadReq) (*pb.UploadResp, error) {
	// 校验扩展名
	if !l.svcCtx.Tools.CheckExtension(in.ObjectName) {
		return nil, errors.New("extension not allowed")
	}
	// 校验大小
	if !l.svcCtx.Tools.CheckFileSize(int64(len(in.Data))) {
		return nil, errors.New("file size exceed limit")
	}
	// 推断 Content-Type
	if in.ContentType == "" {
		in.ContentType = l.svcCtx.Tools.GetContentType(in.ObjectName)
	}

	// 秒传检查
	if in.Hash != "" {
		exist, err := l.svcCtx.FileModel.FindOneByHash(l.ctx, in.Hash)
		if err == nil && exist != nil {
			// 如果需要权限隔离，可以在这里检查 Uploader
			// if exist.Uploader != in.Uid { return error }
			// 简单实现：直接返回已存在文件
			return &pb.UploadResp{
				FileId:     exist.FileId,
				Etag:       exist.Hash, // 假设 Hash 等同 ETag
				Bucket:     exist.Bucket,
				ObjectName: exist.ObjectName,
				Size:       exist.Size,
			}, nil
		}
	}

	stor, err := l.svcCtx.Storage.Select(l.ctx, in.Bucket)
	if err != nil {
		return nil, err
	}
	etag, err := stor.PutObject(l.ctx, in.Bucket, in.ObjectName, in.Data, in.ContentType)
	if err != nil {
		return nil, err
	}
	fileId := uuid.NewString()
	// 如果前端没传 Hash，这里可以考虑计算，或者暂存空
	hash := in.Hash
	if hash == "" {
		hash = etag // 简单用 etag 代替，实际上应该算 md5
	}

	_, err = l.svcCtx.FileModel.Insert(l.ctx, &model.FileInfo{
		FileId:      fileId,
		FileName:    filepath.Base(in.ObjectName),
		Bucket:      in.Bucket,
		ObjectName:  in.ObjectName,
		Size:        in.Size,
		ContentType: in.ContentType,
		Uploader:    in.Uid,
		UploadTime:  time.Now().Unix(),
		Hash:        hash,
		Description: "",
		DeletedTime: 0,
		Status:      1,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UploadResp{
		FileId:     fileId,
		Etag:       etag,
		Bucket:     in.Bucket,
		ObjectName: in.ObjectName,
		Size:       in.Size,
	}, nil
}
