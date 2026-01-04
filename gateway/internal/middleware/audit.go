package middleware

import (
	"time"

	"github.com/cy77cc/go-microstack/common/audit"
	"github.com/gin-gonic/gin"
)

func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		duration := time.Since(startTime)

		// Extract user info if available
		// Assuming "userId" is set in context by some auth middleware, 
		// or we can extract from header "X-User-Id" if gateway has already verified it or passed it.
		// For gateway, it might be the entry point, so maybe it doesn't know userId yet unless verified.
		userID := c.GetString("userId")
		if userID == "" {
			userID = c.GetHeader("X-User-Id")
		}

		audit.Log(c.Request.Context(), audit.AuditLog{
			Timestamp: startTime.UnixMilli(),
			UserID:    userID,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			Duration:  duration.String(),
			ClientIP:  c.ClientIP(),
		})
	}
}
