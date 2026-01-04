package audit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// AuditLog defines the structure of an audit log entry
type AuditLog struct {
	TraceID   string      `json:"trace_id"`
	SpanID    string      `json:"span_id"`
	Timestamp int64       `json:"timestamp"`
	UserID    string      `json:"user_id"`
	Method    string      `json:"method"`
	Path      string      `json:"path"`
	Status    int         `json:"status"`
	Duration  string      `json:"duration"`
	ClientIP  string      `json:"client_ip"`
	Body      interface{} `json:"body,omitempty"` // Be careful with sensitive data
	Error     string      `json:"error,omitempty"`
}

// Log records an audit log entry
func Log(ctx context.Context, entry AuditLog) {
	// In a real system, this might push to Kafka, Elasticsearch, or a database.
	// For now, we use structured logging via logx.

	// Ensure timestamp is set
	if entry.Timestamp == 0 {
		entry.Timestamp = time.Now().UnixMilli()
	}

	content, err := json.Marshal(entry)
	if err != nil {
		logx.Errorf("Failed to marshal audit log: %v", err)
		return
	}

	// Use a specific key/marker to easily filter these logs later
	logx.WithContext(ctx).Infof("[AUDIT] %s", string(content))
}
