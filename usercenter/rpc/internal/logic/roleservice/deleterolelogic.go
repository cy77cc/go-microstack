package roleservicelogic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRoleLogic) DeleteRole(in *pb.DeleteRoleReq) (*pb.LogoutResp, error) {
	_, err := l.svcCtx.RolesModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "role not found")
		}
		return nil, err
	}

	// Delete role permissions
	err = l.svcCtx.RolePermissionsModel.DeleteByRoleId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// Delete user roles
	err = l.svcCtx.UserRolesModel.DeleteByRoleId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// Delete role
	err = l.svcCtx.RolesModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.LogoutResp{Success: true}, nil
}
