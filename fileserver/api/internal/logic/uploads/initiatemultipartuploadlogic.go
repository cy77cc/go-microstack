package uploads

import (
	"context"

	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
