package files

import (
	"context"

	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDownloadUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDownloadUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDownloadUrlLogic {
	return &GetDownloadUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDownloadUrlLogic) GetDownloadUrl() (resp *types.PresignResp, err error) {
	// todo: add your logic here and delete this line

	return
}
