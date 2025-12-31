package consul

import (
	"context"
	"fmt"

	"github.com/cy77cc/go-microstack/common/register/types"
	"github.com/hashicorp/consul/api"
)

type Instance struct {
	client *api.Client
}

func NewConsulInstance(opts *types.Options) (*Instance, error) {
	config := api.DefaultConfig()
	if len(opts.Addrs) > 0 {
		config.Address = opts.Addrs[0]
	}
	if opts.Username != "" && opts.Password != "" {
		config.HttpAuth = &api.HttpBasicAuth{
			Username: opts.Username,
			Password: opts.Password,
		}
	}
	// Token handling can be added here if needed

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Instance{client: client}, nil
}

func (ins *Instance) Register(ctx context.Context, instance *types.ServiceInstance) error {
	registration := &api.AgentServiceRegistration{
		ID:      instance.ID,
		Name:    instance.ServiceName,
		Port:    instance.Port,
		Address: instance.Host,
		Meta:    instance.Metadata,
		Tags:    []string{instance.GroupName}, // Use GroupName as Tag
	}

	if instance.HealthCheck != nil {
		check := &api.AgentServiceCheck{
			Interval: instance.HealthCheck.Interval.String(),
			Timeout:  instance.HealthCheck.Timeout.String(),
		}
		if instance.HealthCheck.Type == "http" {
			check.HTTP = instance.HealthCheck.URL
		} else if instance.HealthCheck.Type == "tcp" {
			check.TCP = instance.HealthCheck.URL
		}
		registration.Check = check
	}

	return ins.client.Agent().ServiceRegister(registration)
}

func (ins *Instance) Deregister(ctx context.Context, instanceID string) error {
	return ins.client.Agent().ServiceDeregister(instanceID)
}

func (ins *Instance) Heartbeat(ctx context.Context, instanceID string) error {
	return ins.client.Agent().PassTTL("service:"+instanceID, "")
}

func (ins *Instance) GetService(ctx context.Context, serviceName string, groupName string) ([]*types.ServiceInstance, error) {
	services, _, err := ins.client.Health().Service(serviceName, groupName, true, nil)
	if err != nil {
		return nil, err
	}

	var instances []*types.ServiceInstance
	for _, entry := range services {
		instances = append(instances, &types.ServiceInstance{
			ID:          entry.Service.ID,
			ServiceName: entry.Service.Service,
			Host:        entry.Service.Address,
			Port:        entry.Service.Port,
			Metadata:    entry.Service.Meta,
			GroupName:   groupName,
		})
	}
	return instances, nil
}

func (ins *Instance) WatchService(ctx context.Context, serviceName string, groupName string) (<-chan []*types.ServiceInstance, error) {
	ch := make(chan []*types.ServiceInstance)
	go func() {
		var lastIndex uint64
		for {
			select {
			case <-ctx.Done():
				return
			default:
				services, meta, err := ins.client.Health().Service(serviceName, groupName, true, &api.QueryOptions{
					WaitIndex: lastIndex,
				})
				if err != nil {
					// Handle error, maybe backoff
					continue
				}
				lastIndex = meta.LastIndex

				var instances []*types.ServiceInstance
				for _, entry := range services {
					instances = append(instances, &types.ServiceInstance{
						ID:          entry.Service.ID,
						ServiceName: entry.Service.Service,
						Host:        entry.Service.Address,
						Port:        entry.Service.Port,
						Metadata:    entry.Service.Meta,
						GroupName:   groupName,
					})
				}
				ch <- instances
			}
		}
	}()
	return ch, nil
}

func (ins *Instance) ListServices(ctx context.Context, groupName string) ([]string, error) {
	services, _, err := ins.client.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}
	var names []string
	for name := range services {
		names = append(names, name)
	}
	return names, nil
}

func (ins *Instance) GetConfig(ctx context.Context, key, group string) (*types.ConfigItem, error) {
	kv, _, err := ins.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, fmt.Errorf("config not found: %s", key)
	}
	return &types.ConfigItem{
		Key:   kv.Key,
		Value: string(kv.Value),
		Group: group,
	}, nil
}

func (ins *Instance) WatchConfig(ctx context.Context, key, group string) (<-chan *types.ConfigItem, error) {
	ch := make(chan *types.ConfigItem)
	go func() {
		var lastIndex uint64
		for {
			select {
			case <-ctx.Done():
				return
			default:
				kv, meta, err := ins.client.KV().Get(key, &api.QueryOptions{
					WaitIndex: lastIndex,
				})
				if err != nil {
					continue
				}
				if kv == nil {
					continue
				}
				lastIndex = meta.LastIndex
				ch <- &types.ConfigItem{
					Key:   kv.Key,
					Value: string(kv.Value),
					Group: group,
				}
			}
		}
	}()
	return ch, nil
}

func (ins *Instance) PublishConfig(ctx context.Context, config *types.ConfigItem) error {
	p := &api.KVPair{Key: config.Key, Value: []byte(config.Value)}
	_, err := ins.client.KV().Put(p, nil)
	return err
}

func (ins *Instance) RemoveConfig(ctx context.Context, key, group string) error {
	_, err := ins.client.KV().Delete(key, nil)
	return err
}

func (ins *Instance) Close() error {
	return nil
}
