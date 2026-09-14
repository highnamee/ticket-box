package middleware

import (
	"log/slog"
	"time"

	"ticket-box-be/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RequestLoggerMiddleware logs HTTP requests with structured attributes using slog
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()
		reqID := logger.GetRequestID(c.Request.Context())

		fullPath := path
		if query != "" {
			fullPath = path + "?" + query
		}

		attrs := []slog.Attr{
			slog.String("request_id", reqID),
			slog.String("method", method),
			slog.String("path", fullPath),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", userAgent),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("error", c.Errors.String()))
		}

		l := logger.FromContext(c.Request.Context())
		if statusCode >= 500 {
			l.LogAttrs(c.Request.Context(), slog.LevelError, "HTTP Request", attrs...)
		} else if statusCode >= 400 {
			l.LogAttrs(c.Request.Context(), slog.LevelWarn, "HTTP Request", attrs...)
		} else {
			l.LogAttrs(c.Request.Context(), slog.LevelInfo, "HTTP Request", attrs...)
		}
	}
}
