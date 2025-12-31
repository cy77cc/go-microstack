package loadbalance

import (
	"math/rand"
	"time"

	"github.com/cy77cc/go-microstack/common/register/types"
)

// Random 随机负载均衡策略
type Random struct{}

// NewRandom 创建随机实例
func NewRandom() *Random {
	return &Random{}
}

func init() {
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())
}

// Select 选择实例
func (r *Random) Select(instances []*types.ServiceInstance) (*types.ServiceInstance, error) {
	if len(instances) == 0 {
		return nil, ErrNoInstances
	}

	index := rand.Intn(len(instances))
	return instances[index], nil
}
