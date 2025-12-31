package middleware

import (
	"github.com/cy77cc/go-microstack/gateway/internal/circuitbreaker"
	"github.com/cy77cc/go-microstack/gateway/internal/config"
	"github.com/gin-gonic/gin"
)

// breakerManager 熔断器管理器
var breakerManager *CircuitBreakerManager

// CircuitBreakerManager 熔断器管理器结构
type CircuitBreakerManager struct {
	breakers map[string]*circuitbreaker.CircuitBreaker
}

// Get 获取指定路由的熔断器
func (cm *CircuitBreakerManager) Get(routeKey string) *circuitbreaker.CircuitBreaker {
	if cb, exists := cm.breakers[routeKey]; exists {
		return cb
	}

	// 如果不存在，创建新的熔断器（使用10个时间窗口桶）
	cb := circuitbreaker.NewCircuitBreaker(10)
	cm.breakers[routeKey] = cb
	return cb
}

// InitBreakerManager 初始化熔断器管理器
func InitBreakerManager() {
	if breakerManager == nil {
		breakerManager = &CircuitBreakerManager{
			breakers: make(map[string]*circuitbreaker.CircuitBreaker),
		}
	}
}

// CircuitBreakerMiddleware 熔断器中间件
func CircuitBreakerMiddleware(route *config.Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取熔断器
		// 使用路由前缀作为唯一标识
		cb := breakerManager.Get(route.PathPrefix)

		// 2. 检查是否允许通过
		if !cb.Allow() {
			c.AbortWithStatus(503) // Service Unavailable
			return
		}

		// 3. 执行请求
		c.Next()

		// 4. 记录结果
		success := c.Writer.Status() < 500
		cb.OnResult(success, route.CircuitBreakerConfig)
	}
}
