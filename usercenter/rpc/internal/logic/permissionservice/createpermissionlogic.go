package permissionservicelogic

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

type CreatePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePermissionLogic {
	return &CreatePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePermissionLogic) CreatePermission(in *pb.CreatePermissionReq) (*pb.PermissionResp, error) {
	_, err := l.svcCtx.PermissionsModel.FindOneByCode(l.ctx, in.Code)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "permission code already exists")
	}
	if !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}

	newPermission := &model.Permissions{
		Name:        in.Name,
		Code:        in.Code,
		Type:        int64(in.Type),
		Resource:    in.Resource,
		Action:      in.Action,
		Description: in.Description,
		Status:      1,
	}

	res, err := l.svcCtx.PermissionsModel.Insert(l.ctx, newPermission)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	newPermission.Id = uint64(id)

	return &pb.PermissionResp{
		Permission: &pb.Permission{
			Id:          newPermission.Id,
			Name:        newPermission.Name,
			Code:        newPermission.Code,
			Type:        pb.PermissionType(newPermission.Type),
			Resource:    newPermission.Resource,
			Action:      newPermission.Action,
			Description: newPermission.Description,
			Status:      int32(newPermission.Status),
			CreateTime:  newPermission.CreateTime,
			UpdateTime:  newPermission.UpdateTime,
		},
	}, nil
}
