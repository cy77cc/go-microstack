package uploads

import (
	"context"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadPartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadPartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadPartLogic {
	return &UploadPartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadPartLogic) UploadPart(req *types.UploadPartReq) (resp *types.UploadPartResp, err error) {
	if req.UploadId == "" || req.PartNumber == 0 {
		return nil, xcode.NewErrCodeMsg(xcode.ErrInvalidParam, "uploadId or partNumber empty")
	}
	data, _ := l.ctx.Value("partData").([]byte)
	if len(data) == 0 {
		return nil, xcode.NewErrCodeMsg(xcode.ErrInvalidParam, "empty part data")
	}
	rpcResp, err := l.svcCtx.FilesRpc.UploadPart(l.ctx, &pb.UploadPartReq{
		UploadId:   req.UploadId,
		PartNumber: req.PartNumber,
		Data:       data,
	})
	if err != nil {
		return nil, err
	}

	return &types.UploadPartResp{
		ETag:       rpcResp.Etag,
		PartNumber: rpcResp.PartNumber,
	}, nil
}
