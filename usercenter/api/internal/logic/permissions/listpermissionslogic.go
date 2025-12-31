package permissions

import (
	"context"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/permissionservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionsLogic {
	return &ListPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPermissionsLogic) ListPermissions(req *types.PermissionListReq) (resp *types.PermissionListResp, err error) {
	res, err := l.svcCtx.PermissionService.ListPermissions(l.ctx, &permissionservice.ListPermissionsReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.PermissionResp, 0, len(res.Permissions))
	for _, perm := range res.Permissions {
		list = append(list, types.PermissionResp{
			Id:          perm.Id,
			Name:        perm.Name,
			Code:        perm.Code,
			Type:        int32(perm.Type),
			Resource:    perm.Resource,
			Action:      perm.Action,
			Description: perm.Description,
			Status:      perm.Status,
			CreateTime:  perm.CreateTime,
			UpdateTime:  perm.UpdateTime,
		})
	}

	return &types.PermissionListResp{
		Total: res.Total,
		List:  list,
	}, nil
}
