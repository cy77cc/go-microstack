package users

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRoleToUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssignRoleToUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRoleToUserLogic {
	return &AssignRoleToUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignRoleToUserLogic) AssignRoleToUser(userId uint64, req *types.AssignRoleBody) error {
	_, err := l.svcCtx.UserService.AssignRoleToUser(l.ctx, &userservice.AssignRoleReq{
		Uid:     userId,
		RoleIds: req.RoleIds,
	})
	return err
}
