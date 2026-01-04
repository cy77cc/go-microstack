package uploads

import (
	"context"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"
	"github.com/cy77cc/go-microstack/fileserver/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AbortMultipartUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAbortMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbortMultipartUploadLogic {
	return &AbortMultipartUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AbortMultipartUploadLogic) AbortMultipartUpload(req *types.AbortMultipartReq) error {
	if req.UploadId == "" {
		return xcode.NewErrCodeMsg(xcode.ErrInvalidParam, "uploadId empty")
	}
	uid, _ := l.ctx.Value("uid").(uint64)
	_, err := l.svcCtx.FilesRpc.AbortMultipartUpload(l.ctx, &pb.AbortMultipartUploadReq{
		UploadId: req.UploadId,
		Uid:      uid,
	})
	if err != nil {
		return err
	}

	return nil
}
