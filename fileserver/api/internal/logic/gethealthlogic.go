package logic

import (
	"context"

	"github.com/cy77cc/go-microstack/fileserver/api/internal/svc"
	"github.com/cy77cc/go-microstack/fileserver/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHealthLogic {
	return &GetHealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHealthLogic) GetHealth() (resp *types.HealthResp, err error) {
	return &types.HealthResp{
		Status: "ok",
	}, nil
}
