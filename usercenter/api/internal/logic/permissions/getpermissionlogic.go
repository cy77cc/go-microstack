package permissions

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/permissionservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionLogic {
	return &GetPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPermissionLogic) GetPermission(id uint64) (resp *types.PermissionResp, err error) {
	res, err := l.svcCtx.PermissionService.GetPermission(l.ctx, &permissionservice.GetPermissionReq{
		Id: id,
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
