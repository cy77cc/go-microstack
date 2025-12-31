package roles

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/roleservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRolesLogic) ListRoles(req *types.RoleListReq) (resp *types.RoleListResp, err error) {
	res, err := l.svcCtx.RoleService.ListRoles(l.ctx, &roleservice.ListRolesReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.RoleResp, 0, len(res.Roles))
	for _, role := range res.Roles {
		list = append(list, types.RoleResp{
			Id:          role.Id,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Status:      role.Status,
			CreateTime:  role.CreateTime,
			UpdateTime:  role.UpdateTime,
		})
	}

	return &types.RoleListResp{
		Total: res.Total,
		List:  list,
	}, nil
}
