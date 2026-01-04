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

type DeletePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePermissionLogic) DeletePermission(in *pb.DeletePermissionReq) (*pb.LogoutResp, error) {
	_, err := l.svcCtx.PermissionsModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xcode.NewErrCodeMsg(xcode.NotFound, "permission not found")
		}
		return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "database error")
	}

	// Delete permission
	err = l.svcCtx.PermissionsModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "delete permission failed")
	}

	return &pb.LogoutResp{Success: true}, nil
}
