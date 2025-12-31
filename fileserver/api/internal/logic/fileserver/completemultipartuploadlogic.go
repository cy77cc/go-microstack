package fileserver

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteMultipartUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteMultipartUploadLogic {
	return &CompleteMultipartUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteMultipartUploadLogic) CompleteMultipartUpload(req *types.CompleteMultipartReq) (resp *types.CompleteMultipartResp, err error) {
	if req.UploadId == "" || len(req.Parts) == 0 {
		return nil, errors.New("uploadId or parts empty")
	}
	uid, _ := l.ctx.Value("uid").(uint64)
	var parts []*pb.CompletedPart
	for _, p := range req.Parts {
		parts = append(parts, &pb.CompletedPart{
			PartNumber: p.PartNumber,
			Etag:       p.ETag,
		})
	}
	rpcResp, err := l.svcCtx.FilesRpc.CompleteMultipartUpload(l.ctx, &pb.CompleteMultipartUploadReq{
		UploadId: req.UploadId,
		Parts:    parts,
		Uid:      uid,
	})
	if err != nil {
		return nil, err
	}

	return &types.CompleteMultipartResp{
		FileId: rpcResp.FileId,
		Bucket: rpcResp.Bucket,
		Key:    rpcResp.ObjectName,
		Size:   rpcResp.Size,
	}, nil
}
