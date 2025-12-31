package nacos

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cy77cc/go-microstack/common/logx"
	"github.com/cy77cc/go-microstack/common/register/types"
	"github.com/cy77cc/go-microstack/common/utils"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// NewNacosInstance 创建Nacos实例
func NewNacosInstance(cfg *Config) (*Instance, error) {
	sc := make([]constant.ServerConfig, 0, 1)
	for _, endpoint := range cfg.Endpoints {

		logx.Infof("Connecting to Nacos: %s", endpoint)
		ip := strings.Split(endpoint, ":")[0]
		port, err := strconv.ParseUint(strings.Split(endpoint, ":")[1], 10, 64)

		if err != nil {
			logx.Errorf("failed to parse endpoint: %s", endpoint)
			continue
		}

		sc = append(sc, constant.ServerConfig{
			IpAddr:      ip,
			Port:        port,
			ContextPath: cfg.ContextPath,
		})
	}

	cc := &constant.ClientConfig{
		NamespaceId: cfg.Namespace,
		LogLevel:    "warn",
		Username:    cfg.Username,
		Password:    cfg.Password,
		ContextPath: cfg.ContextPath,
		AccessKey:   cfg.IdentityKey,
		SecretKey:   cfg.IdentityVal,
		LogDir:      "log",
		TimeoutMs:   cfg.TimeoutMs,
	}

	namingClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ServerConfigs: sc,
		ClientConfig:  cc,
	})
	if err != nil {
		return nil, err
	}

	configClient, err := clients.NewConfigClient(vo.NacosClientParam{
		ServerConfigs: sc,
		ClientConfig:  cc,
	})

	if err != nil {
		return nil, err
	}

	return &Instance{
		clientConfig: cc,
		serverConfig: sc,
		ConfigClient: configClient,
		NamingClient: namingClient,
	}, nil
}

func NewNacosConfig() *Config {
	return &Config{
		Endpoints: make([]string, 1),
	}
}

// LoadNacosEnv loads Nacos config from environment variables
func (c *Config) LoadNacosEnv() {
	c.Endpoints[0] = os.Getenv("NACOS_ADDR")

	if port := os.Getenv("NACOS_PORT"); port != "" {
		c.Port, _ = strconv.ParseUint(port, 10, 64)
	} else {
		c.Port = 8848
	}

	c.Namespace = os.Getenv("NACOS_NAMESPACE")
	c.ContextPath = os.Getenv("NACOS_CONTEXT_PATH")
	if c.ContextPath == "" {
		c.ContextPath = "/nacos"
	}
	c.Username = os.Getenv("NACOS_USERNAME")
	c.Password = os.Getenv("NACOS_PASSWORD")
	c.IdentityKey = os.Getenv("NACOS_AUTH_IDENTITY_KEY")
	c.IdentityVal = os.Getenv("NACOS_AUTH_IDENTITY_VALUE")
	c.Token = os.Getenv("NACOS_AUTH_TOKEN")
}

// Register 实现服务注册
func (ins *Instance) Register(ctx context.Context, instance *types.ServiceInstance) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 参数验证
	if instance == nil {
		return fmt.Errorf("service instance cannot be nil")
	}
	if instance.Host == "" {
		return fmt.Errorf("service instance host cannot be empty")
	}
	if instance.Port <= 0 || instance.Port > 65535 {
		return fmt.Errorf("service instance port must be between 1 and 65535, got %d", instance.Port)
	}
	if instance.ServiceName == "" {
		return fmt.Errorf("service instance service name cannot be empty")
	}

	// 注册实例到Nacos
	_, err := ins.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          utils.GetMachineIP(),
		Port:        uint64(instance.Port),
		ServiceName: instance.ServiceName,
		Weight:      instance.Weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    instance.Metadata,
		GroupName:   instance.GroupName,
	})

	if err != nil {
		return fmt.Errorf("failed to register instance to Nacos: %w", err)
	}

	return err
}

// Deregister 实现服务注销
func (ins *Instance) Deregister(ctx context.Context, instanceID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 从Nacos注销实例
	_, err := ins.NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          instanceID,
		Port:        0,          // 需要从实例ID中解析IP和端口
		ServiceName: instanceID, // 实际上需要正确的服务名
		Ephemeral:   true,
	})

	return err
}

// GetService 实现服务发现
func (ins *Instance) GetService(ctx context.Context, serviceName string, groupName string) ([]*types.ServiceInstance, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}

	instances, err := ins.NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		HealthyOnly: true,
		GroupName:   groupName,
	})
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	var result []*types.ServiceInstance
	for _, inst := range instances {
		result = append(result, &types.ServiceInstance{
			ID:          inst.InstanceId,
			ServiceName: inst.ServiceName,
			Host:        inst.Ip,
			Port:        int(inst.Port),
			Metadata:    inst.Metadata,
			Weight:      inst.Weight,
		})
	}
	return result, nil
}

// WatchService 监听服务变化
func (ins *Instance) WatchService(ctx context.Context, serviceName string, groupName string) (<-chan []*types.ServiceInstance, error) {
	resultChan := make(chan []*types.ServiceInstance, 1)

	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}

	// 注册服务监听器
	err := ins.NamingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		SubscribeCallback: func(services []model.Instance, err error) {
			if err != nil {
				logx.Errorf("Error in service subscribe: %v", err)
				return
			}

			var instances []*types.ServiceInstance
			for _, service := range services {
				instances = append(instances, &types.ServiceInstance{
					ID:          service.InstanceId,
					ServiceName: service.ServiceName,
					Host:        service.Ip,
					Port:        int(service.Port),
					Metadata:    service.Metadata,
					Weight:      service.Weight,
				})
			}

			select {
			case resultChan <- instances:
			default:
				// 如果通道已满，跳过本次更新
			}
		},
	})

	if err != nil {
		return nil, err
	}

	return resultChan, nil
}

// ListServices 获取所有服务名称
func (ins *Instance) ListServices(ctx context.Context, groupName string) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}

	services, err := ins.NamingClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: ins.clientConfig.NamespaceId,
		GroupName: groupName,
		PageNo:    1,
		PageSize:  100,
	})
	if err != nil {
		return nil, err
	}

	return services.Doms, nil
}

// GetConfig 获取配置
func (ins *Instance) GetConfig(ctx context.Context, key, group string) (*types.ConfigItem, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if group == "" {
		group = "DEFAULT_GROUP"
	}

	content, err := ins.ConfigClient.GetConfig(vo.ConfigParam{
		DataId: key,
		Group:  group,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get config from nacos: %w", err)
	}

	return &types.ConfigItem{
		Key:   key,
		Value: content,
		Group: group,
	}, nil
}

// WatchConfig 监听配置变更
func (ins *Instance) WatchConfig(ctx context.Context, key, group string) (<-chan *types.ConfigItem, error) {
	configChan := make(chan *types.ConfigItem, 1)

	if group == "" {
		group = "DEFAULT_GROUP"
	}

	onChangeCallback := func(namespace, group, dataId, data string) {
		config := &types.ConfigItem{
			Key:   dataId,
			Value: data,
			Group: group,
		}

		select {
		case configChan <- config:
		default:
			// 如果通道已满，跳过本次更新
		}
	}

	err := ins.ConfigClient.ListenConfig(vo.ConfigParam{
		DataId:   key,
		Group:    group,
		OnChange: onChangeCallback,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to listen config changes: %w", err)
	}

	return configChan, nil
}

// PublishConfig 发布配置
func (ins *Instance) PublishConfig(ctx context.Context, config *types.ConfigItem) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	success, err := ins.ConfigClient.PublishConfig(vo.ConfigParam{
		DataId:  config.Key,
		Group:   config.Group,
		Content: config.Value,
	})
	if err != nil {
		return fmt.Errorf("failed to publish config: %w", err)
	}

	if !success {
		return fmt.Errorf("publish config failed, dataId: %s, group: %s", config.Key, config.Group)
	}

	return nil
}

// RemoveConfig 删除配置
func (ins *Instance) RemoveConfig(ctx context.Context, key, group string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if group == "" {
		group = "DEFAULT_GROUP"
	}

	success, err := ins.ConfigClient.DeleteConfig(vo.ConfigParam{
		DataId: key,
		Group:  group,
	})
	if err != nil {
		return fmt.Errorf("failed to remove config: %w", err)
	}

	if !success {
		return fmt.Errorf("remove config failed, dataId: %s, group: %s", key, group)
	}

	return nil
}

// Close 关闭连接
func (ins *Instance) Close() error {
	// 关闭Nacos客户端连接
	if ins.NamingClient != nil {
		// Nacos客户端没有直接的关闭方法，但可以取消所有订阅
	}
	if ins.ConfigClient != nil {
		// 清理配置监听器
	}

	return nil
}

// Heartbeat 发送心跳
func (ins *Instance) Heartbeat(ctx context.Context, instanceID string) error {
	// Nacos使用临时实例，自动处理心跳
	// 这里可以实现额外的心跳逻辑
	return nil
}
