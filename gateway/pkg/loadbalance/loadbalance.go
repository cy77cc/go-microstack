package loadbalance

import (
	"errors"

	"github.com/cy77cc/go-microstack/common/register/types"
)

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
	// Select 选择一个服务实例
	Select(instances []*types.ServiceInstance) (*types.ServiceInstance, error)
}

// ErrNoInstances 无可用实例错误
var ErrNoInstances = errors.New("no instances available")

// NewRoundRobinLoadBalancer 创建轮询负载均衡器
func NewRoundRobinLoadBalancer() LoadBalancer {
	return NewRoundRobin()
}

// NewRandomLoadBalancer 创建随机负载均衡器
func NewRandomLoadBalancer() LoadBalancer {
	return NewRandom()
}
