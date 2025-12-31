package logic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/fileserver/model"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadPartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadPartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadPartLogic {
	return &UploadPartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadPartLogic) UploadPart(in *pb.UploadPartReq) (*pb.UploadPartResp, error) {
	if in.UploadId == "" || in.PartNumber == 0 {
		return nil, errors.New("uploadId or partNumber empty")
	}
	rec, err := l.svcCtx.UploadModel.FindOneByUploadId(l.ctx, in.UploadId)
	if err != nil {
		return nil, err
	}
	// idempotency
	if part, err := l.svcCtx.UploadPartMod.FindOneByUploadIdPartNumber(l.ctx, in.UploadId, int64(in.PartNumber)); err == nil && part != nil {
		return &pb.UploadPartResp{Etag: part.Etag, PartNumber: uint32(part.PartNumber)}, nil
	}
	stor, err := l.svcCtx.Storage.Select(l.ctx, rec.Bucket)
	if err != nil {
		return nil, err
	}
	etag, err := stor.UploadPart(l.ctx, rec.Bucket, rec.ObjectName, in.UploadId, int(in.PartNumber), in.Data)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.UploadPartMod.Insert(l.ctx, &model.MultipartPart{
		UploadId:   in.UploadId,
		PartNumber: int64(in.PartNumber),
		Etag:       etag,
		Size:       int64(len(in.Data)),
	})
	if err != nil {
		return nil, err
	}

	return &pb.UploadPartResp{Etag: etag, PartNumber: in.PartNumber}, nil
}
