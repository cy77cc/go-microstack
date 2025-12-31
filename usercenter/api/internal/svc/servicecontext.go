package svc

import (
	"github.com/cy77cc/go-microstack/usercenter/api/internal/config"
	"github.com/cy77cc/go-microstack/usercenter/api/internal/middleware"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/authservice"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/permissionservice"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/roleservice"
	"github.com/cy77cc/go-microstack/usercenter/rpc/client/userservice"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	SignatureMiddleware rest.Middleware
	AuthService         authservice.AuthService
	UserService         userservice.UserService
	RoleService         roleservice.RoleService
	PermissionService   permissionservice.PermissionService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		SignatureMiddleware: middleware.NewSignatureMiddleware().Handle,
		AuthService:         authservice.NewAuthService(zrpc.MustNewClient(c.UserCenterRpc)),
		UserService:         userservice.NewUserService(zrpc.MustNewClient(c.UserCenterRpc)),
		RoleService:         roleservice.NewRoleService(zrpc.MustNewClient(c.UserCenterRpc)),
		PermissionService:   permissionservice.NewPermissionService(zrpc.MustNewClient(c.UserCenterRpc)),
	}
}
