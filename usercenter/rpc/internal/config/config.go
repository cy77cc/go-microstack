package config

import (
	"github.com/cy77cc/go-microstack/common/register"
	"github.com/cy77cc/go-microstack/common/types"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Register register.Config `json:",optional"`

	Mysql      types.Mysql     `json:",optional"`
	CacheRedis cache.CacheConf `json:",optional"`
	Salt       string          `json:",optional"`
	JwtAuth    struct {
		AccessSecret string
		AccessExpire int64
	} `json:",optional"`
}
