package users

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeRoleFromUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRevokeRoleFromUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeRoleFromUserLogic {
	return &RevokeRoleFromUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeRoleFromUserLogic) RevokeRoleFromUser(userId uint64, req *types.RevokeRoleBody) error {
	_, err := l.svcCtx.UserService.RevokeRoleFromUser(l.ctx, &userservice.RevokeRoleReq{
		Uid:     userId,
		RoleIds: req.RoleIds,
	})
	return err
}
