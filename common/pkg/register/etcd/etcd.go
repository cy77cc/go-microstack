package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/cy77cc/go-microstack/common/pkg/register/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Instance struct {
	client    *clientv3.Client
	leaseIDs  sync.Map // map[string]clientv3.LeaseID
	namespace string
}

func NewEtcdInstance(opts *types.Options) (*Instance, error) {
	cfg := clientv3.Config{
		Endpoints:   opts.Addrs,
		DialTimeout: opts.Timeout,
		Username:    opts.Username,
		Password:    opts.Password,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	return &Instance{
		client:    client,
		namespace: opts.Namespace,
	}, nil
}

// buildServiceKey 构造服务注册的 Key
// 格式: /<namespace>/services/<group>/<serviceName>/<instanceID>
func (ins *Instance) buildServiceKey(serviceName, groupName, instanceID string) string {
	var sb strings.Builder
	sb.WriteString("/")
	if ins.namespace != "" {
		sb.WriteString(ins.namespace)
		sb.WriteString("/")
	}
	sb.WriteString("services/")
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}
	sb.WriteString(groupName)
	sb.WriteString("/")
	sb.WriteString(serviceName)
	sb.WriteString("/")
	sb.WriteString(instanceID)
	return sb.String()
}

// buildServicePrefix 构造服务发现的前缀
func (ins *Instance) buildServicePrefix(serviceName, groupName string) string {
	var sb strings.Builder
	sb.WriteString("/")
	if ins.namespace != "" {
		sb.WriteString(ins.namespace)
		sb.WriteString("/")
	}
	sb.WriteString("services/")
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}
	sb.WriteString(groupName)
	sb.WriteString("/")
	sb.WriteString(serviceName)
	sb.WriteString("/")
	return sb.String()
}

// buildConfigKey 构造配置的 Key
// 格式: /<namespace>/config/<group>/<key>
func (ins *Instance) buildConfigKey(key, group string) string {
	var sb strings.Builder
	sb.WriteString("/")
	if ins.namespace != "" {
		sb.WriteString(ins.namespace)
		sb.WriteString("/")
	}
	sb.WriteString("config/")
	if group == "" {
		group = "DEFAULT_GROUP"
	}
	sb.WriteString(group)
	sb.WriteString("/")
	sb.WriteString(key)
	return sb.String()
}

func (ins *Instance) Register(ctx context.Context, instance *types.ServiceInstance) error {
	key := ins.buildServiceKey(instance.ServiceName, instance.GroupName, instance.ID)
	val, err := json.Marshal(instance)
	if err != nil {
		return err
	}

	// 创建 Lease
	leaseResp, err := ins.client.Grant(ctx, 10) // 10s TTL
	if err != nil {
		return err
	}

	_, err = ins.client.Put(ctx, key, string(val), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	// Keep alive
	ch, err := ins.client.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return err
	}

	// 存储 LeaseID 以便注销
	ins.leaseIDs.Store(instance.ID, leaseResp.ID)

	go func() {
		for range ch {
			// consume keepalive responses
		}
		// 如果 keepalive 通道关闭，可能需要重新注册或处理错误，这里简化处理
	}()

	return nil
}

func (ins *Instance) Deregister(ctx context.Context, instanceID string) error {
	// 由于接口只提供了 instanceID，我们无法直接构建 Key (缺少 ServiceName 等)
	// 但如果我们保存了 leaseID，我们可以直接 Revoke Lease，这会导致 Key 自动删除
	if val, ok := ins.leaseIDs.Load(instanceID); ok {
		leaseID := val.(clientv3.LeaseID)
		_, err := ins.client.Revoke(ctx, leaseID)
		ins.leaseIDs.Delete(instanceID)
		return err
	}
	return nil
}

// Heartbeat Etcd 使用 Lease KeepAlive 机制，通常不需要手动心跳，除非是为了兼容接口
func (ins *Instance) Heartbeat(ctx context.Context, instanceID string) error {
	return nil
}

func (ins *Instance) GetService(ctx context.Context, serviceName string, groupName string) ([]*types.ServiceInstance, error) {
	prefix := ins.buildServicePrefix(serviceName, groupName)
	resp, err := ins.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var instances []*types.ServiceInstance
	for _, kv := range resp.Kvs {
		var inst types.ServiceInstance
		if err := json.Unmarshal(kv.Value, &inst); err != nil {
			continue
		}
		instances = append(instances, &inst)
	}
	return instances, nil
}

func (ins *Instance) WatchService(ctx context.Context, serviceName string, groupName string) (<-chan []*types.ServiceInstance, error) {
	prefix := ins.buildServicePrefix(serviceName, groupName)
	ch := make(chan []*types.ServiceInstance)

	go func() {
		rch := ins.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
		for range rch {
			// 当有变化时，重新获取完整列表（简化实现）
			// 也可以解析 Event 进行增量更新
			insts, err := ins.GetService(context.Background(), serviceName, groupName)
			if err == nil {
				ch <- insts
			}
		}
	}()
	return ch, nil
}

func (ins *Instance) ListServices(ctx context.Context, groupName string) ([]string, error) {
	// 这是一个比较重的操作，需要扫描指定 Group 下的所有 Key
	// 前缀: /<namespace>/services/<group>/
	var sb strings.Builder
	sb.WriteString("/")
	if ins.namespace != "" {
		sb.WriteString(ins.namespace)
		sb.WriteString("/")
	}
	sb.WriteString("services/")
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}
	sb.WriteString(groupName)
	sb.WriteString("/")

	prefix := sb.String()

	resp, err := ins.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}

	servicesMap := make(map[string]struct{})
	for _, kv := range resp.Kvs {
		// Key 格式: prefix/<serviceName>/<instanceID>
		key := string(kv.Key)
		suffix := strings.TrimPrefix(key, prefix)
		parts := strings.Split(suffix, "/")
		if len(parts) >= 1 {
			servicesMap[parts[0]] = struct{}{}
		}
	}

	var list []string
	for k := range servicesMap {
		list = append(list, k)
	}
	return list, nil
}

func (ins *Instance) GetConfig(ctx context.Context, key, group string) (*types.ConfigItem, error) {
	etcdKey := ins.buildConfigKey(key, group)
	resp, err := ins.client.Get(ctx, etcdKey)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("config not found: %s", key)
	}
	return &types.ConfigItem{
		Key:   key,
		Value: string(resp.Kvs[0].Value),
		Group: group,
	}, nil
}

func (ins *Instance) WatchConfig(ctx context.Context, key, group string) (<-chan *types.ConfigItem, error) {
	etcdKey := ins.buildConfigKey(key, group)
	ch := make(chan *types.ConfigItem)

	go func() {
		rch := ins.client.Watch(context.Background(), etcdKey)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				ch <- &types.ConfigItem{
					Key:   key,
					Value: string(ev.Kv.Value),
					Group: group,
				}
			}
		}
	}()
	return ch, nil
}

func (ins *Instance) PublishConfig(ctx context.Context, config *types.ConfigItem) error {
	etcdKey := ins.buildConfigKey(config.Key, config.Group)
	_, err := ins.client.Put(ctx, etcdKey, config.Value)
	return err
}

func (ins *Instance) RemoveConfig(ctx context.Context, key, group string) error {
	etcdKey := ins.buildConfigKey(key, group)
	_, err := ins.client.Delete(ctx, etcdKey)
	return err
}

func (ins *Instance) Close() error {
	return ins.client.Close()
}
