package svc

import (
	"fmt"

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
	conn := sqlx.NewMysql(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=10s",
		c.Mysql.Username, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Database))

	c.Mysql.Endpoint = fmt.Sprintf("%s:%d", c.Mysql.Host, c.Mysql.Port)

	rdb := redis.MustNewRedis(c.Redis.RedisConf)
	
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
