package roles

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/roleservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.RoleCreateReq) (resp *types.RoleResp, err error) {
	res, err := l.svcCtx.RoleService.CreateRole(l.ctx, &roleservice.CreateRoleReq{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
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
