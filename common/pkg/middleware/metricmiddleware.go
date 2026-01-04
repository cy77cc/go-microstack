package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/timex"
)

const serverNamespace = "microstack"

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path", "method", "code"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "method", "code"},
	})
)

type MetricMiddleware struct {
}

func NewMetricMiddleware() *MetricMiddleware {
	return &MetricMiddleware{}
}

func (m *MetricMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := timex.Now()

		sw := &statusWriter{ResponseWriter: w}

		next(sw, r)

		duration := timex.Since(startTime)
		code := strconv.Itoa(sw.status)

		// Record metrics
		metricServerReqDur.Observe(int64(duration/time.Millisecond), r.URL.Path, r.Method, code)
		metricServerReqCodeTotal.Inc(r.URL.Path, r.Method, code)
	}
}

// Reusing statusWriter from auditmiddleware if in same package,
// but since we are in same package 'middleware', we don't need to redefine it if it's not exported.
// However, in Go, unexported types are package private.
// If auditmiddleware.go is in the same package 'middleware', we can reuse `statusWriter`.
// I will assume they are in the same package and `statusWriter` is defined in `auditmiddleware.go`.
// If I need to be safe, I can rename it or check.
// Since I just wrote `auditmiddleware.go` with `type statusWriter struct`, it is available here.
