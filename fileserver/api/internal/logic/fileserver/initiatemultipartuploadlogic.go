package fileserver

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitiateMultipartUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitiateMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitiateMultipartUploadLogic {
	return &InitiateMultipartUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitiateMultipartUploadLogic) InitiateMultipartUpload(req *types.InitiateMultipartReq) (resp *types.InitiateMultipartResp, err error) {
	if req.Bucket == "" || req.Key == "" {
		return nil, errors.New("bucket or key empty")
	}
	uid, _ := l.ctx.Value("uid").(uint64)
	rpcResp, err := l.svcCtx.FilesRpc.InitiateMultipartUpload(l.ctx, &pb.InitiateMultipartUploadReq{
		Bucket:      req.Bucket,
		ObjectName:  req.Key,
		Size:        req.Size,
		ContentType: req.ContentType,
		Hash:        req.Hash,
		Uid:         uid,
	})
	if err != nil {
		return nil, err
	}

	return &types.InitiateMultipartResp{
		UploadId: rpcResp.UploadId,
	}, nil
}
