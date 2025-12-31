package roles

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/roleservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.RoleUpdateReq) (resp *types.RoleResp, err error) {
	res, err := l.svcCtx.RoleService.UpdateRole(l.ctx, &roleservice.UpdateRoleReq{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
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
