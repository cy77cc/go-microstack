package buckets

import (
	"context"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateBucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBucketLogic {
	return &CreateBucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBucketLogic) CreateBucket(req *types.BucketReq) (resp *types.BucketResp, err error) {
	if req.Bucket == "" {
		return nil, xcode.NewErrCodeMsg(xcode.ErrInvalidParam, "bucket empty")
	}
	var st pb.StorageType
	switch req.StorageType {
	case 1:
		st = pb.StorageType_STORAGE_MINIO
	case 2:
		st = pb.StorageType_STORAGE_LOCAL
	default:
		st = pb.StorageType_STORAGE_UNSPECIFIED
	}
	rpcResp, err := l.svcCtx.FilesRpc.CreateBucket(l.ctx, &pb.CreateBucketReq{
		Bucket:      req.Bucket,
		StorageType: st,
	})
	if err != nil {
		return nil, err
	}

	return &types.BucketResp{
		Bucket: rpcResp.Bucket,
	}, nil
}
