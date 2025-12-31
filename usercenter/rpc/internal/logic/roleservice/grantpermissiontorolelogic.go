package roleservicelogic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GrantPermissionToRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGrantPermissionToRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantPermissionToRoleLogic {
	return &GrantPermissionToRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GrantPermissionToRoleLogic) GrantPermissionToRole(in *pb.GrantPermissionReq) (*pb.LogoutResp, error) {
	for _, permissionId := range in.PermissionIds {
		_, err := l.svcCtx.RolePermissionsModel.FindOneByRoleIdPermissionId(l.ctx, in.RoleId, permissionId)
		if errors.Is(err, model.ErrNotFound) {
			_, err = l.svcCtx.RolePermissionsModel.Insert(l.ctx, &model.RolePermissions{
				RoleId:       in.RoleId,
				PermissionId: permissionId,
			})
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
	}

	return &pb.LogoutResp{Success: true}, nil
}
