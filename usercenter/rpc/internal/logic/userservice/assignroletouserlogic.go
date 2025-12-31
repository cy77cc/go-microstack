package userservicelogic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRoleToUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignRoleToUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRoleToUserLogic {
	return &AssignRoleToUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignRoleToUserLogic) AssignRoleToUser(in *pb.AssignRoleReq) (*pb.LogoutResp, error) {
	for _, roleId := range in.RoleIds {
		_, err := l.svcCtx.UserRolesModel.FindOneByUserIdRoleId(l.ctx, in.Uid, roleId)
		if errors.Is(err, model.ErrNotFound) {
			_, err = l.svcCtx.UserRolesModel.Insert(l.ctx, &model.UserRoles{
				UserId: in.Uid,
				RoleId: roleId,
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
