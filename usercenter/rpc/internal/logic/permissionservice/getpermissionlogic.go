package permissionservicelogic

import (
	"context"
	"errors"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionLogic {
	return &GetPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPermissionLogic) GetPermission(in *pb.GetPermissionReq) (*pb.PermissionResp, error) {
	permission, err := l.svcCtx.PermissionsModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xcode.NewErrCodeMsg(xcode.NotFound, "permission not found")
		}
		return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
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
