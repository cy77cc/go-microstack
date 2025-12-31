package config

import (
	"github.com/cy77cc/go-microstack/common/register"
	"github.com/cy77cc/go-microstack/common/types"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Register   register.Config `json:",optional"`
	CacheRedis cache.CacheConf
	Mysql      types.Mysql `json:"mysql" yaml:"mysql"`
	Minio      types.Minio `json:"minio" yaml:"minio"`

	Upload Upload `json:"upload" yaml:"upload"`
	Local  Local  `json:"local" yaml:"local"`
}


type Local struct {
	BaseDir string `yaml:"base-dir" json:"base-dir"`
}

type Upload struct {
	MaxFileSize int64    `yaml:"max-file-size" json:"max-file-size"` // Max file size in bytes, 0 for unlimited
	AllowedExts []string `yaml:"allowed-exts" json:"allowed-exts"`   // Allowed extensions, empty for all
}
