package userservicelogic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeRoleFromUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeRoleFromUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeRoleFromUserLogic {
	return &RevokeRoleFromUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RevokeRoleFromUserLogic) RevokeRoleFromUser(in *pb.RevokeRoleReq) (*pb.LogoutResp, error) {
	for _, roleId := range in.RoleIds {
		userRole, err := l.svcCtx.UserRolesModel.FindOneByUserIdRoleId(l.ctx, in.Uid, roleId)
		if err == nil {
			err = l.svcCtx.UserRolesModel.Delete(l.ctx, userRole.Id)
			if err != nil {
				return nil, err
			}
		} else if !errors.Is(err, model.ErrNotFound) {
			return nil, err
		}
	}

	return &pb.LogoutResp{Success: true}, nil
}
