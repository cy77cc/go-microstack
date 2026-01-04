package config

import (
	"github.com/cy77cc/go-microstack/common/pkg/register"
	"github.com/cy77cc/go-microstack/common/pkg/types"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Register   register.Config `json:",optional"`
	CacheRedis cache.CacheConf
	Mysql      types.Mysql `json:"mysql" yaml:"mysql"`
	Minio      types.Minio `json:"minio" yaml:"minio"`

	Upload Upload
	Local  Local
}

type Local struct {
	BaseDir string
}

type Upload struct {
	MaxFileSize int64    // Max file size in bytes, 0 for unlimited
	AllowedExts []string // Allowed extensions, empty for all
}
