package logic

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"

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

func (l *GetHealthLogic) GetHealth() (resp *types.GetHealthResp, err error) {
	return &types.GetHealthResp{
		Status: "ok",
	}, nil
}
