package roleservicelogic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokePermissionFromRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokePermissionFromRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokePermissionFromRoleLogic {
	return &RevokePermissionFromRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RevokePermissionFromRoleLogic) RevokePermissionFromRole(in *pb.RevokePermissionReq) (*pb.LogoutResp, error) {
	for _, permissionId := range in.PermissionIds {
		rolePermission, err := l.svcCtx.RolePermissionsModel.FindOneByRoleIdPermissionId(l.ctx, in.RoleId, permissionId)
		if err == nil {
			err = l.svcCtx.RolePermissionsModel.Delete(l.ctx, rolePermission.Id)
			if err != nil {
				return nil, err
			}
		} else if !errors.Is(err, model.ErrNotFound) {
			return nil, err
		}
	}

	return &pb.LogoutResp{Success: true}, nil
}
