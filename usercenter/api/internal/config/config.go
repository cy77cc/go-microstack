package config

import (
	"github.com/cy77cc/go-microstack/common/register"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Register register.Config `json:",optional"`

	Auth struct {
		AccessSecret string
		AccessExpire int64
	} `json:",optional"`
	UserCenterRpc zrpc.RpcClientConf `json:",optional"`
}
