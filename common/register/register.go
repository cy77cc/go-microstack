package register

import (
	"context"
	"errors"
	"time"

	"github.com/cy77cc/go-microstack/common/register/consul"
	"github.com/cy77cc/go-microstack/common/register/etcd"
	"github.com/cy77cc/go-microstack/common/register/nacos"
	"github.com/cy77cc/go-microstack/common/register/types"
	"github.com/cy77cc/go-microstack/common/register/zookeeper"
)

// WithEndpoints 设置注册中心地址
func WithEndpoints(endpoints ...string) types.Option {
	return func(opts *types.Options) {
		opts.Addrs = endpoints
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) types.Option {
	return func(opts *types.Options) {
		opts.Timeout = timeout
	}
}

// WithAuth 设置认证信息
func WithAuth(username, password string) types.Option {
	return func(opts *types.Options) {
		opts.Username = username
		opts.Password = password
	}
}

// WithNamespace 设置命名空间
func WithNamespace(namespace string) types.Option {
	return func(opts *types.Options) {
		opts.Namespace = namespace
	}
}

// NewRegister 创建注册中心实例
func NewRegister(ctx context.Context, regType string, opts ...types.Option) (types.Register, error) {
	options := &types.Options{
		Timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	switch regType {
	case "nacos":
		return newNacosRegister(ctx, options)
	case "etcd":
		return etcd.NewEtcdInstance(options)
	case "consul":
		return consul.NewConsulInstance(options)
	case "zookeeper":
		return zookeeper.NewZookeeperInstance(options)
	default:
		return nil, errors.New("unsupported register type: " + regType)
	}
}

func newNacosRegister(ctx context.Context, options *types.Options) (types.Register, error) {
	cfg := nacos.NewNacosConfig()
	cfg.Endpoints = options.Addrs
	cfg.TimeoutMs = uint64(options.Timeout.Milliseconds())
	cfg.Username = options.Username
	cfg.Password = options.Password
	cfg.Namespace = options.Namespace
	cfg.Group = options.ClusterName
	cfg.ContextPath = "/nacos"
	cfg.IdentityKey = "identityKey"
	cfg.IdentityVal = "identityVal"
	return nacos.NewNacosInstance(cfg)
}
