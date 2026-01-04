package permissionservicelogic

import (
	"context"

	"github.com/cy77cc/go-microstack/common/pkg/xcode"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionsLogic {
	return &ListPermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPermissionsLogic) ListPermissions(in *pb.ListPermissionsReq) (*pb.ListPermissionsResp, error) {
	permissions, total, err := l.svcCtx.PermissionsModel.FindAll(l.ctx, int64(in.Page), int64(in.PageSize))
	if err != nil {
		return nil, xcode.NewErrCodeMsg(xcode.DatabaseError, "list permissions failed")
	}

	var pbPermissions []*pb.Permission
	for _, permission := range permissions {
		pbPermissions = append(pbPermissions, &pb.Permission{
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
		})
	}

	return &pb.ListPermissionsResp{
		Permissions: pbPermissions,
		Total:       total,
	}, nil
}
