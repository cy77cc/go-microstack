package roleservicelogic

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolePermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolePermissionsLogic {
	return &ListRolePermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListRolePermissionsLogic) ListRolePermissions(in *pb.ListRolePermissionsReq) (*pb.ListRolePermissionsResp, error) {
	rolePermissions, err := l.svcCtx.RolePermissionsModel.FindAllByRoleId(l.ctx, in.RoleId)
	if err != nil {
		return nil, err
	}

	var permissions []*pb.Permission
	for _, rp := range rolePermissions {
		permission, err := l.svcCtx.PermissionsModel.FindOne(l.ctx, rp.PermissionId)
		if err != nil {
			continue
		}
		permissions = append(permissions, &pb.Permission{
			Id:          permission.Id,
			Code:        permission.Code,
			Name:        permission.Name,
			Description: permission.Description,
			Type:        pb.PermissionType(permission.Type),
			Resource:    permission.Resource,
			Action:      permission.Action,
			Status:      int32(permission.Status),
			CreateTime:  permission.CreateTime,
			UpdateTime:  permission.UpdateTime,
		})
	}

	return &pb.ListRolePermissionsResp{
		Permissions: permissions,
	}, nil
}
