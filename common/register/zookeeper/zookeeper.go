package zookeeper

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cy77cc/go-microstack/common/register/types"
	"github.com/go-zookeeper/zk"
)

type Instance struct {
	conn *zk.Conn
}

func NewZookeeperInstance(opts *types.Options) (*Instance, error) {
	conn, _, err := zk.Connect(opts.Addrs, opts.Timeout)
	if err != nil {
		return nil, err
	}
	return &Instance{conn: conn}, nil
}

func (ins *Instance) Register(ctx context.Context, instance *types.ServiceInstance) error {
	// Create path: /services/serviceName/instanceID
	basePath := "/services"
	if exists, _, _ := ins.conn.Exists(basePath); !exists {
		ins.conn.Create(basePath, []byte{}, 0, zk.WorldACL(zk.PermAll))
	}
	servicePath := basePath + "/" + instance.ServiceName
	if exists, _, _ := ins.conn.Exists(servicePath); !exists {
		ins.conn.Create(servicePath, []byte{}, 0, zk.WorldACL(zk.PermAll))
	}

	nodePath := servicePath + "/" + instance.ID
	data, _ := json.Marshal(instance)

	_, err := ins.conn.Create(nodePath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

func (ins *Instance) Deregister(ctx context.Context, instanceID string) error {
	// Need serviceName again.
	// Assuming structure /services/serviceName/instanceID
	// Without serviceName, difficult.
	return nil
}

func (ins *Instance) Heartbeat(ctx context.Context, instanceID string) error {
	return nil // ZK uses ephemeral nodes tied to session
}

func (ins *Instance) GetService(ctx context.Context, serviceName string, groupName string) ([]*types.ServiceInstance, error) {
	path := "/services/" + serviceName
	children, _, err := ins.conn.Children(path)
	if err != nil {
		return nil, err
	}

	var instances []*types.ServiceInstance
	for _, child := range children {
		data, _, err := ins.conn.Get(path + "/" + child)
		if err != nil {
			continue
		}
		var inst types.ServiceInstance
		if err := json.Unmarshal(data, &inst); err == nil {
			instances = append(instances, &inst)
		}
	}
	return instances, nil
}

func (ins *Instance) WatchService(ctx context.Context, serviceName string, groupName string) (<-chan []*types.ServiceInstance, error) {
	ch := make(chan []*types.ServiceInstance)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				path := "/services/" + serviceName
				_, _, events, err := ins.conn.ChildrenW(path)
				if err != nil {
					time.Sleep(time.Second)
					continue
				}

				// Fetch current state
				insts, _ := ins.GetService(ctx, serviceName, groupName)
				ch <- insts

				// Wait for event
				select {
				case <-events:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return ch, nil
}

func (ins *Instance) ListServices(ctx context.Context, groupName string) ([]string, error) {
	children, _, err := ins.conn.Children("/services")
	return children, err
}

func (ins *Instance) GetConfig(ctx context.Context, key, group string) (*types.ConfigItem, error) {
	data, _, err := ins.conn.Get(key)
	if err != nil {
		return nil, err
	}
	return &types.ConfigItem{
		Key:   key,
		Value: string(data),
		Group: group,
	}, nil
}

func (ins *Instance) WatchConfig(ctx context.Context, key, group string) (<-chan *types.ConfigItem, error) {
	ch := make(chan *types.ConfigItem)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data, _, events, err := ins.conn.GetW(key)
				if err != nil {
					time.Sleep(time.Second)
					continue
				}

				ch <- &types.ConfigItem{
					Key:   key,
					Value: string(data),
					Group: group,
				}

				select {
				case <-events:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return ch, nil
}

func (ins *Instance) PublishConfig(ctx context.Context, config *types.ConfigItem) error {
	// Check if exists
	exists, _, _ := ins.conn.Exists(config.Key)
	if !exists {
		// Create parent paths if needed... simplified here
		_, err := ins.conn.Create(config.Key, []byte(config.Value), 0, zk.WorldACL(zk.PermAll))
		return err
	}
	_, err := ins.conn.Set(config.Key, []byte(config.Value), -1)
	return err
}

func (ins *Instance) RemoveConfig(ctx context.Context, key, group string) error {
	return ins.conn.Delete(key, -1)
}

func (ins *Instance) Close() error {
	ins.conn.Close()
	return nil
}
