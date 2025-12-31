package permissions

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/permissionservice"
	"github.com/cy77cc/go-microstack/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePermissionLogic {
	return &CreatePermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePermissionLogic) CreatePermission(req *types.PermissionCreateReq) (resp *types.PermissionResp, err error) {
	res, err := l.svcCtx.PermissionService.CreatePermission(l.ctx, &permissionservice.CreatePermissionReq{
		Name:        req.Name,
		Code:        req.Code,
		Type:        pb.PermissionType(req.Type),
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
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
