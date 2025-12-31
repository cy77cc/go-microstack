package roles

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/roleservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GrantPermissionToRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGrantPermissionToRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantPermissionToRoleLogic {
	return &GrantPermissionToRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GrantPermissionToRoleLogic) GrantPermissionToRole(roleId uint64, req *types.GrantPermissionBody) error {
	_, err := l.svcCtx.RoleService.GrantPermissionToRole(l.ctx, &roleservice.GrantPermissionReq{
		RoleId:        roleId,
		PermissionIds: req.PermissionIds,
	})
	return err
}
