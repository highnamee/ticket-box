package http

import (
	"ticket-box-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	HealthHandler      *HealthHandler
	CORSAllowedOrigins []string
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.New()

	// Global Middlewares
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Health Check Route
	r.GET("/health", cfg.HealthHandler.Check)

	// API V1 Group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", cfg.HealthHandler.Check)
	}

	return r
}
