package http

import (
	"log/slog"

	_ "ticket-box-be/docs"
	"ticket-box-be/internal/admin"
	"ticket-box-be/internal/config"
	"ticket-box-be/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	AppConfig          *config.Config
	AppEnv             string
	HealthHandler      *HealthHandler
	TicketHandler      *TicketHandler
	CORSAllowedOrigins []string
	EnableAdmin        bool
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	switch cfg.AppEnv {
	case string(config.EnvProduction):
		gin.SetMode(gin.ReleaseMode)
	case string(config.EnvTest):
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// Global Middlewares
	r.Use(gin.Recovery())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RequestLoggerMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Mount GoAdmin Admin Panel
	if cfg.AppConfig != nil && cfg.EnableAdmin {
		if err := admin.Mount(r, cfg.AppConfig); err != nil {
			slog.Warn("Failed to mount GoAdmin engine", "error", err)
		}
	}

	// Swagger Documentation UI (disabled in production)
	if cfg.AppEnv != string(config.EnvProduction) {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Health Check Route
	if cfg.HealthHandler != nil {
		r.GET("/health", cfg.HealthHandler.Check)
	}

	// API V1 Group
	v1 := r.Group("/api/v1")
	{
		if cfg.HealthHandler != nil {
			v1.GET("/ping", cfg.HealthHandler.Check)
		}
		if cfg.TicketHandler != nil {
			v1.GET("/tickets", cfg.TicketHandler.GetPublicTickets)
		}
	}

	return r
}
