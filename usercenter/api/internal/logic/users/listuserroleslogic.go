package users

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserRolesLogic {
	return &ListUserRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserRolesLogic) ListUserRoles(userId uint64) (resp *types.RoleListResp, err error) {
	res, err := l.svcCtx.UserService.ListUserRoles(l.ctx, &userservice.ListUserRolesReq{
		Uid: userId,
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
		Total: int64(len(list)),
		List:  list,
	}, nil
}
