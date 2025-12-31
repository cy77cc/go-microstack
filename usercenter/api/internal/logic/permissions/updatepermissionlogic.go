package permissions

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/permissionservice"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePermissionLogic {
	return &UpdatePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePermissionLogic) UpdatePermission(req *types.PermissionUpdateReq) (resp *types.PermissionResp, err error) {
	res, err := l.svcCtx.PermissionService.UpdatePermission(l.ctx, &permissionservice.UpdatePermissionReq{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Resource:    req.Resource,
		Action:      req.Action,
		Type:        pb.PermissionType(req.Type),
	})
	if err != nil {
		return nil, err
	}

	return &types.PermissionResp{
		Id:          res.Permission.Id,
		Name:        res.Permission.Name,
		Code:        res.Permission.Code,
		Type:        int32(res.Permission.Type),
		Resource:    res.Permission.Resource,
		Action:      res.Permission.Action,
		Description: res.Permission.Description,
		Status:      res.Permission.Status,
		CreateTime:  res.Permission.CreateTime,
		UpdateTime:  res.Permission.UpdateTime,
	}, nil
}
