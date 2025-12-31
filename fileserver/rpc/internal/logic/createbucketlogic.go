package logic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/fileserver/model"
	"github.com/cy77cc/go-microstack/fileserver/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBucketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBucketLogic {
	return &CreateBucketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateBucketLogic) CreateBucket(in *pb.CreateBucketReq) (*pb.CreateBucketResp, error) {
	if in.Bucket == "" {
		return nil, errors.New("bucket empty")
	}
	var st int64
	switch in.StorageType {
	case pb.StorageType_STORAGE_MINIO:
		st = 1
	default:
		st = 2
	}
	_, err := l.svcCtx.BucketModel.Insert(l.ctx, &model.BucketConfig{
		Bucket:      in.Bucket,
		StorageType: st,
		Endpoint:    l.svcCtx.Config.Minio.Endpoint,
		Region:      "",
		IsPublic:    0,
	})
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		logx.Errorf("insert bucket config error: %v", err)
	}
	stor, err := l.svcCtx.Storage.Select(l.ctx, in.Bucket)
	if err != nil {
		return nil, err
	}
	if err = stor.CreateBucket(l.ctx, in.Bucket); err != nil {
		return nil, err
	}

	return &pb.CreateBucketResp{Bucket: in.Bucket}, nil
}
