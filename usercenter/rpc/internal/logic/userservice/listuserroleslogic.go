package userservicelogic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserRolesLogic {
	return &ListUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUserRolesLogic) ListUserRoles(in *pb.ListUserRolesReq) (*pb.ListUserRolesResp, error) {
	userRoles, err := l.svcCtx.UserRolesModel.FindAllByUserId(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}

	var roles []*pb.Role
	for _, ur := range userRoles {
		role, err := l.svcCtx.RolesModel.FindOne(l.ctx, ur.RoleId)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				continue
			}
			return nil, err
		}
		roles = append(roles, &pb.Role{
			Id:          role.Id,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Status:      int32(role.Status),
			CreateTime:  role.CreateTime,
			UpdateTime:  role.UpdateTime,
		})
	}

	return &pb.ListUserRolesResp{
		Roles: roles,
	}, nil
}
