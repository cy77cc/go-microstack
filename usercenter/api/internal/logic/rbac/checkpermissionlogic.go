package rbac

import (
	"context"
	"encoding/json"

	"github.com/cy77cc/go-microstack/usercenter/api/internal/svc"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/types"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/permissionservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPermissionLogic {
	return &CheckPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckPermissionLogic) CheckPermission(req *types.CheckPermissionReq) (resp *types.CheckPermissionResp, err error) {
	// 从 context 中获取 JWT 解析出来的 userId
	userId := l.ctx.Value("userId")
	var uid uint64
	if userId != nil {
		// 根据 userId 的实际类型进行断言转换
		switch v := userId.(type) {
		case uint64:
			uid = v
		case int64:
			uid = uint64(v)
		case float64:
			uid = uint64(v)
		case json.Number:
			if id, err := v.Int64(); err == nil {
				uid = uint64(id)
			}
		}
	}

	// 调用 RPC 服务检查权限
	res, err := l.svcCtx.PermissionService.CheckPermission(l.ctx, &permissionservice.CheckPermissionReq{
		Uid:      uid,
		Resource: req.Resource,
		Action:   req.Action,
	})
	if err != nil {
		return nil, err
	}

	return &types.CheckPermissionResp{
		Allowed: res.Allowed,
	}, nil
}
