package roles

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/roleservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokePermissionFromRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRevokePermissionFromRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokePermissionFromRoleLogic {
	return &RevokePermissionFromRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokePermissionFromRoleLogic) RevokePermissionFromRole(roleId uint64, req *types.RevokePermissionBody) error {
	_, err := l.svcCtx.RoleService.RevokePermissionFromRole(l.ctx, &roleservice.RevokePermissionReq{
		RoleId:        roleId,
		PermissionIds: req.PermissionIds,
	})
	return err
}
