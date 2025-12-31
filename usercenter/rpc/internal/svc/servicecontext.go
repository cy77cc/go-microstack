package svc

import (
	"github.com/cy77cc/go-microstack/usercenter/model"
	"github.com/cy77cc/go-microstack/usercenter/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                 config.Config
	UsersModel             model.UsersModel
	RolesModel             model.RolesModel
	PermissionsModel       model.PermissionsModel
	UserRolesModel         model.UserRolesModel
	RolePermissionsModel   model.RolePermissionsModel
	AuthRefreshTokensModel model.AuthRefreshTokensModel
	Rdb                    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	var rdb *redis.Redis
	if c.CacheRedis != nil {
		rdb = redis.MustNewRedis(c.CacheRedis[0].RedisConf)
	}
	return &ServiceContext{
		Config:                 c,
		UsersModel:             model.NewUsersModel(conn, c.CacheRedis),
		RolesModel:             model.NewRolesModel(conn, c.CacheRedis),
		PermissionsModel:       model.NewPermissionsModel(conn, c.CacheRedis),
		UserRolesModel:         model.NewUserRolesModel(conn, c.CacheRedis),
		RolePermissionsModel:   model.NewRolePermissionsModel(conn, c.CacheRedis),
		AuthRefreshTokensModel: model.NewAuthRefreshTokensModel(conn, c.CacheRedis),
		Rdb:                    rdb,
	}
}
