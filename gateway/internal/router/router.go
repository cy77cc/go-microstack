package router

import (
	"net/http"

	"github.com/cy77cc/go-microstack/gateway/internal/config"
	"github.com/cy77cc/go-microstack/gateway/internal/middleware"
	"github.com/cy77cc/go-microstack/gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

// Router 路由管理器
type Router struct {
}

// NewRouter 创建路由管理器
func NewRouter() *Router {
	return &Router{}
}

// RegisterRoutes 注册路由
func (*Router) RegisterRoutes(r *gin.Engine, routes []config.Route, proxyHandler *proxy.Handler) {
	// 全局中间件
	r.Use(middleware.AuditMiddleware())
	r.Use(middleware.MetricMiddleware())

	// 初始化中间件管理器（单例模式，只需调用一次）
	middleware.InitBreakerManager()
	middleware.InitBucketManager()

	for _, route := range routes {
		// 复制 route 变量，避免闭包捕获循环变量问题
		currentRoute := route

		var handlers []gin.HandlerFunc

		// 1. 添加熔断器中间件
		if currentRoute.CircuitBreakerConfig != nil {
			handlers = append(handlers, middleware.CircuitBreakerMiddleware(&currentRoute))
		}

		// 2. 添加限流中间件
		if currentRoute.RateLimitConfig != nil {
			handlers = append(handlers, middleware.RateLimitMiddleware(&currentRoute))
		}

		// 3. 添加代理处理器
		handlers = append(handlers, proxyHandler.HandleRoute(currentRoute.Service, currentRoute.StripPrefix))

		// 注册路由
		r.Any(currentRoute.PathPrefix+"/*path", handlers...)
	}

	// 默认/兜底路由 (不带特定限流/熔断配置，或者使用默认配置)
	r.Any("/api/:service/*path", proxyHandler.HandleGeneric)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
