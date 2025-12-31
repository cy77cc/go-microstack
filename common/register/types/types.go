package types

import (
	"context"
	"errors"
	"time"
)

// Option 函数类型，用于配置选项
type Option func(*Options)

// 常用错误定义
var (
	ErrServiceNotFound = errors.New("service not found")
	ErrConfigNotFound  = errors.New("config not found")
	ErrInvalidInstance = errors.New("invalid service instance")
	ErrConnectFailed   = errors.New("connect to register failed")
)

// ServiceInstance 服务实例信息
type ServiceInstance struct {
	ID          string            // 实例唯一ID
	ServiceName string            // 服务名称
	Host        string            // 主机地址
	Port        int               // 端口
	Endpoints   []string          // 服务端点
	Metadata    map[string]string // 元数据
	HealthCheck *HealthCheck      // 健康检查配置
	Timestamp   time.Time         // 注册时间
	Weight      float64           // 权重
	GroupName   string            // 服务组名
	ClusterName string            // 集群名称
}

// HealthCheck 健康检查配置
type HealthCheck struct {
	Type     string        // 检查类型: http, tcp, grpc
	URL      string        // 检查URL
	Interval time.Duration // 检查间隔
	Timeout  time.Duration // 超时时间
}

// ConfigItem 配置项
type ConfigItem struct {
	Key    string // 配置键
	Value  string // 配置值
	Group  string // 配置组
	Format string // 配置格式(json, yaml, properties等)
}

// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
	// GetService 获取指定服务的所有实例
	GetService(ctx context.Context, serviceName string, groupName string) ([]*ServiceInstance, error)

	// WatchService 监听服务变化
	WatchService(ctx context.Context, serviceName string, groupName string) (<-chan []*ServiceInstance, error)

	// ListServices 获取所有服务名称
	ListServices(ctx context.Context, groupName string) ([]string, error)
}

// ServiceRegistry 服务注册接口
type ServiceRegistry interface {
	// Register 注册服务实例
	Register(ctx context.Context, instance *ServiceInstance) error

	// Deregister 注销服务实例
	Deregister(ctx context.Context, instanceID string) error

	// Heartbeat 服务心跳/健康检查
	Heartbeat(ctx context.Context, instanceID string) error
}

// ConfigManager 配置管理接口
type ConfigManager interface {
	// GetConfig 获取配置
	GetConfig(ctx context.Context, key, group string) (*ConfigItem, error)

	// WatchConfig 监听配置变更
	WatchConfig(ctx context.Context, key, group string) (<-chan *ConfigItem, error)

	// PublishConfig 发布配置
	PublishConfig(ctx context.Context, config *ConfigItem) error

	// RemoveConfig 删除配置
	RemoveConfig(ctx context.Context, key, group string) error
}

// Register 通用注册中心接口，整合服务注册发现和配置管理
type Register interface {
	ServiceRegistry
	ServiceDiscovery
	ConfigManager

	// Close 关闭注册中心连接
	Close() error
}

// Options 注册中心配置选项
type Options struct {
	Addrs       []string      // 注册中心地址
	Timeout     time.Duration // 超时时间
	Username    string        // 用户名
	Password    string        // 密码
	Namespace   string        // 命名空间
	ClusterName string        // 集群名称
}
