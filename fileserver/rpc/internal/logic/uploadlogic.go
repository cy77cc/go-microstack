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
	// 1. 校验文件扩展名是否允许
	if !l.svcCtx.Tools.CheckExtension(in.ObjectName) {
		return nil, errors.New("extension not allowed")
	}
	// 2. 校验文件大小是否超过限制
	if !l.svcCtx.Tools.CheckFileSize(int64(len(in.Data))) {
		return nil, errors.New("file size exceed limit")
	}
	// 3. 推断 Content-Type (如果未提供)
	if in.ContentType == "" {
		in.ContentType = l.svcCtx.Tools.GetContentType(in.ObjectName)
	}

	// 4. 秒传（重复文件检测）检查
	// 如果提供了 Hash，则检查数据库中是否已存在相同 Hash 的文件
	if in.Hash != "" {
		exist, err := l.svcCtx.FileModel.FindOneByHash(l.ctx, in.Hash)
		if err == nil && exist != nil {
			// TODO: 如果需要权限隔离，可以在这里检查 Uploader 是否匹配，或者增加引用计数
			// if exist.Uploader != in.Uid { return error }
			
			// 简单实现：直接返回已存在文件的信息，实现秒传效果
			return &pb.UploadResp{
				FileId:     exist.FileId,
				Etag:       exist.Hash, // 假设 Hash 等同 ETag
				Bucket:     exist.Bucket,
				ObjectName: exist.ObjectName,
				Size:       exist.Size,
			}, nil
		}
	}

	// 5. 选择存储后端（如 MinIO）并上传文件
	stor, err := l.svcCtx.Storage.Select(l.ctx, in.Bucket)
	if err != nil {
		return nil, err
	}
	etag, err := stor.PutObject(l.ctx, in.Bucket, in.ObjectName, in.Data, in.ContentType)
	if err != nil {
		return nil, err
	}

	// 6. 生成唯一 FileID 并保存元数据到数据库
	fileId := uuid.NewString()
	// 如果前端没传 Hash，这里可以考虑计算，或者暂存空。建议前端计算 Hash 以支持秒传。
	hash := in.Hash
	if hash == "" {
		hash = etag // 简单用 etag 代替，注意：ETag 不一定等于 MD5
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
