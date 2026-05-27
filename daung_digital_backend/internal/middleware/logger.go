package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger middleware logs HTTP requests
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		
		c.Next()
		
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		
		logger.Info("HTTP request",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.Duration("latency", latency),
		)
	}
}
