package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/metric"
)

const gatewayNamespace = "gateway"

var (
	metricGatewayReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: gatewayNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "gateway requests duration(ms).",
		Labels:    []string{"path", "method", "code"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricGatewayReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: gatewayNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "gateway requests error count.",
		Labels:    []string{"path", "method", "code"},
	})
)

func MetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		code := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		metricGatewayReqDur.Observe(int64(duration/time.Millisecond), path, c.Request.Method, code)
		metricGatewayReqCodeTotal.Inc(path, c.Request.Method, code)
	}
}
