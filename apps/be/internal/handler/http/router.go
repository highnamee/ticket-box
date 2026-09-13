package http

import (
	_ "ticket-box-be/docs"
	"ticket-box-be/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	AppEnv             string
	HealthHandler      *HealthHandler
	CORSAllowedOrigins []string
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.New()

	// Global Middlewares
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Swagger Documentation UI (disabled in production)
	if cfg.AppEnv != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Health Check Route
	r.GET("/health", cfg.HealthHandler.Check)

	// API V1 Group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", cfg.HealthHandler.Check)
	}

	return r
}
