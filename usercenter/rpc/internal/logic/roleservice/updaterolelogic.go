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

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *pb.UpdateRoleReq) (*pb.RoleResp, error) {
	role, err := l.svcCtx.RolesModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "role not found")
		}
		return nil, err
	}

	// Check if code is being updated and if it exists
	// if in.Code != "" && in.Code != role.Code {
	// 	_, err := l.svcCtx.RolesModel.FindOneByCode(l.ctx, in.Code)
	// 	if err == nil {
	// 		return nil, status.Error(codes.AlreadyExists, "role code already exists")
	// 	}
	// 	if err != model.ErrNotFound {
	// 		return nil, err
	// 	}
	// 	role.Code = in.Code
	// }

	if in.Name != "" {
		role.Name = in.Name
	}
	if in.Description != "" {
		role.Description = in.Description
	}
	if in.Status != 0 {
		role.Status = int64(in.Status)
	}

	err = l.svcCtx.RolesModel.Update(l.ctx, role)
	if err != nil {
		return nil, err
	}

	return &pb.RoleResp{
		Role: &pb.Role{
			Id:          role.Id,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Status:      int32(role.Status),
			CreateTime:  role.CreateTime,
			UpdateTime:  role.UpdateTime,
		},
	}, nil
}
