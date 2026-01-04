package loadbalance

import (
	"sync/atomic"

	"github.com/cy77cc/go-microstack/common/pkg/register/types"
)

// RoundRobin 轮询负载均衡策略
type RoundRobin struct {
	counter uint64
}

// NewRoundRobin 创建轮询实例
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

// Select 选择实例
func (rb *RoundRobin) Select(instances []*types.ServiceInstance) (*types.ServiceInstance, error) {
	if len(instances) == 0 {
		return nil, ErrNoInstances
	}

	// 原子递增
	count := atomic.AddUint64(&rb.counter, 1)
	index := (count - 1) % uint64(len(instances))
	return instances[index], nil
}
