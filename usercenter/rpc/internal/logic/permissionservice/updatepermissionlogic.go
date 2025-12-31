package permissionservicelogic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePermissionLogic {
	return &UpdatePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePermissionLogic) UpdatePermission(in *pb.UpdatePermissionReq) (*pb.PermissionResp, error) {
	permission, err := l.svcCtx.PermissionsModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "permission not found")
		}
		return nil, err
	}

	// if in.Code != "" && in.Code != permission.Code {
	// 	_, err := l.svcCtx.PermissionsModel.FindOneByCode(l.ctx, in.Code)
	// 	if err == nil {
	// 		return nil, status.Error(codes.AlreadyExists, "permission code already exists")
	// 	}
	// 	if err != model.ErrNotFound {
	// 		return nil, err
	// 	}
	// 	permission.Code = in.Code
	// }

	if in.Name != "" {
		permission.Name = in.Name
	}
	if in.Type != 0 {
		permission.Type = int64(in.Type)
	}
	if in.Resource != "" {
		permission.Resource = in.Resource
	}
	if in.Action != "" {
		permission.Action = in.Action
	}
	if in.Description != "" {
		permission.Description = in.Description
	}
	if in.Status != 0 {
		permission.Status = int64(in.Status)
	}

	err = l.svcCtx.PermissionsModel.Update(l.ctx, permission)
	if err != nil {
		return nil, err
	}

	return &pb.PermissionResp{
		Permission: &pb.Permission{
			Id:          permission.Id,
			Name:        permission.Name,
			Code:        permission.Code,
			Type:        pb.PermissionType(permission.Type),
			Resource:    permission.Resource,
			Action:      permission.Action,
			Description: permission.Description,
			Status:      int32(permission.Status),
			CreateTime:  permission.CreateTime,
			UpdateTime:  permission.UpdateTime,
		},
	}, nil
}
