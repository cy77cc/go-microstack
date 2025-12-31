package roles

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/roleservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolePermissionsLogic {
	return &ListRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRolePermissionsLogic) ListRolePermissions(roleId uint64) (resp *types.PermissionListResp, err error) {
	res, err := l.svcCtx.RoleService.ListRolePermissions(l.ctx, &roleservice.ListRolePermissionsReq{
		RoleId: roleId,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.PermissionResp, 0, len(res.Permissions))
	for _, perm := range res.Permissions {
		list = append(list, types.PermissionResp{
			Id:          perm.Id,
			Name:        perm.Name,
			Code:        perm.Code,
			Type:        int32(perm.Type),
			Resource:    perm.Resource,
			Action:      perm.Action,
			Description: perm.Description,
			Status:      perm.Status,
			CreateTime:  perm.CreateTime,
			UpdateTime:  perm.UpdateTime,
		})
	}

	return &types.PermissionListResp{
		Total: int64(len(list)),
		List:  list,
	}, nil
}
