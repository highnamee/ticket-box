package middleware

import (
	"slices"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a configured CORS middleware instance
// If allowedOrigins is empty, all origins are allowed (suitable for development)
// In production, origins are checked against the provided whitelist
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// If no specific origins are configured, allow all (dev mode)
			if len(allowedOrigins) == 0 {
				return true
			}
			// In production, verify against domain whitelist
			return slices.Contains(allowedOrigins, origin)
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept",
			"Accept-Encoding",
			"Authorization",
			"X-CSRF-Token",
			"X-Requested-With",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
