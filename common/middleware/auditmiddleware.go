package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cy77cc/go-microstack/common/audit"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuditMiddleware struct {
}

func NewAuditMiddleware() *AuditMiddleware {
	return &AuditMiddleware{}
}

func (m *AuditMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Wrap ResponseWriter to capture status code
		sw := &statusWriter{ResponseWriter: w}

		// Read body for audit (be careful with large bodies or sensitive data)
		// For now, we skip reading body to avoid performance impact and security risks
		// unless explicitly needed.

		next(sw, r)

		duration := time.Since(startTime)

		// Extract User ID from context (assuming it was set by JWT middleware)
		userID := ""
		if uid, ok := r.Context().Value("userId").(string); ok {
			userID = uid
		}
		// If userId is json.Number
		if uid, ok := r.Context().Value("userId").(json.Number); ok {
			userID = uid.String()
		}

		// Log audit entry
		audit.Log(r.Context(), audit.AuditLog{
			Timestamp: startTime.UnixMilli(),
			UserID:    userID,
			Method:    r.Method,
			Path:      r.URL.Path,
			Status:    sw.status,
			Duration:  duration.String(),
			ClientIP:  httpx.GetRemoteAddr(r),
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}
