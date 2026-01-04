package permissionservicelogic

import (
	"context"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPermissionLogic {
	return &CheckPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckPermissionLogic) CheckPermission(in *pb.CheckPermissionReq) (*pb.CheckPermissionResp, error) {
	// 1. Get user roles
	userRoles, err := l.svcCtx.UserRolesModel.FindAllByUserId(l.ctx, in.Uid)
	if err != nil {
		return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
	}

	// 2. Check if any role has the permission
	for _, ur := range userRoles {
		rolePermissions, err := l.svcCtx.RolePermissionsModel.FindAllByRoleId(l.ctx, ur.RoleId)
		if err != nil {
			return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
		}

		for _, rp := range rolePermissions {
			permission, err := l.svcCtx.PermissionsModel.FindOne(l.ctx, rp.PermissionId)
			if err != nil {
				continue
			}
			if permission.Resource == in.Resource && permission.Action == in.Action {
				return &pb.CheckPermissionResp{Allowed: true}, nil
			}
		}
	}

	return &pb.CheckPermissionResp{Allowed: false}, nil
}
