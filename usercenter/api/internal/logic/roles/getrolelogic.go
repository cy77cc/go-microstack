package roles

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/roleservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleLogic) GetRole(id uint64) (resp *types.RoleResp, err error) {
	res, err := l.svcCtx.RoleService.GetRole(l.ctx, &roleservice.GetRoleReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &types.RoleResp{
		Id:          res.Role.Id,
		Name:        res.Role.Name,
		Code:        res.Role.Code,
		Description: res.Role.Description,
		Status:      res.Role.Status,
		CreateTime:  res.Role.CreateTime,
		UpdateTime:  res.Role.UpdateTime,
	}, nil
}
