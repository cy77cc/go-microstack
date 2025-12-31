package logic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/fileserver/model"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitiateMultipartUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitiateMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitiateMultipartUploadLogic {
	return &InitiateMultipartUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitiateMultipartUploadLogic) InitiateMultipartUpload(in *pb.InitiateMultipartUploadReq) (*pb.InitiateMultipartUploadResp, error) {
	if in.Bucket == "" || in.ObjectName == "" {
		return nil, errors.New("bucket or object empty")
	}

	// 1. 校验扩展名
	if !l.svcCtx.Tools.CheckExtension(in.ObjectName) {
		return nil, errors.New("extension not allowed")
	}
	// 2. 校验大小
	if in.Size > 0 && !l.svcCtx.Tools.CheckFileSize(in.Size) {
		return nil, errors.New("file size exceed limit")
	}
	// 3. 推断 Content-Type
	if in.ContentType == "" {
		in.ContentType = l.svcCtx.Tools.GetContentType(in.ObjectName)
	}

	// 4. 选择存储后端
	stor, err := l.svcCtx.Storage.Select(l.ctx, in.Bucket)
	if err != nil {
		return nil, err
	}
	// 5. 调用存储后端（MinIO）初始化分片上传，获取 UploadID
	uploadID, err := stor.InitiateMultipart(l.ctx, in.Bucket, in.ObjectName, in.ContentType)
	if err != nil {
		return nil, err
	}
	// 6. 在数据库中记录分片上传任务信息
	_, err = l.svcCtx.UploadModel.Insert(l.ctx, &model.MultipartUpload{
		UploadId:     uploadID,
		Bucket:       in.Bucket,
		ObjectName:   in.ObjectName,
		Size:         in.Size,
		ContentType:  in.ContentType,
		Uploader:     in.Uid,
		Hash:         in.Hash,
		Status:       0, // 0: 初始化/上传中
		CompleteTime: 0,
	})
	if err != nil {
		return nil, err
	}

	return &pb.InitiateMultipartUploadResp{UploadId: uploadID}, nil
}
