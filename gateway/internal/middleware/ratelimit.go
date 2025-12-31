package middleware

import (
	"fmt"
	"sync"

	"github.com/cy77cc/go-microstack/gateway/internal/config"
	"github.com/cy77cc/go-microstack/gateway/internal/ratelimit"
	"github.com/gin-gonic/gin"
)

var bucketManager *BucketManager

// BucketManager 令牌桶管理器
type BucketManager struct {
	buckets map[string]*ratelimit.TokenBucket
	mutex   sync.RWMutex
}

// Get 获取指定路由的令牌桶
// route: 路由配置
// 返回: 令牌桶实例
func (bm *BucketManager) Get(route *config.Route) *ratelimit.TokenBucket {
	if route == nil || route.RateLimitConfig == nil {
		return nil
	}

	key := fmt.Sprintf("%s_%s", route.Service, route.PathPrefix)

	bm.mutex.RLock()
	if bucket, exists := bm.buckets[key]; exists {
		bm.mutex.RUnlock()
		return bucket
	}
	bm.mutex.RUnlock()

	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if bucket, exists := bm.buckets[key]; exists {
		return bucket
	}

	bucket := ratelimit.NewTokenBucket(int64(route.RateLimitConfig.Burst), int64(route.RateLimitConfig.QPS))
	bm.buckets[key] = bucket
	return bucket
}

// InitBucketManager 初始化桶管理器
func InitBucketManager() {
	if bucketManager == nil {
		bucketManager = &BucketManager{
			buckets: make(map[string]*ratelimit.TokenBucket),
		}
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(route *config.Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		if route == nil || route.RateLimitConfig == nil {
			c.Next()
			return
		}

		bucket := bucketManager.Get(route)
		if bucket == nil || !bucket.Allow() {
			c.AbortWithStatus(429) // Too Many Requests
			return
		}

		c.Next()
	}
}
