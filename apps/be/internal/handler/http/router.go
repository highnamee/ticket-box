package http

import (
	"log/slog"

	_ "ticket-box-be/docs"
	"ticket-box-be/internal/admin"
	"ticket-box-be/internal/config"
	"ticket-box-be/internal/middleware"
	"ticket-box-be/internal/pkg/validator"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type RouterConfig struct {
	AppConfig          *config.Config
	AppEnv             string
	HealthHandler      *HealthHandler
	TicketHandler      *TicketHandler
	CORSAllowedOrigins []string
	EnableAdmin        bool
	DB                 *gorm.DB
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	// Register custom validation rules for Gin
	validator.RegisterCustomValidators()

	appEnv := config.Environment(cfg.AppEnv)
	switch appEnv {
	case config.EnvProduction:
		gin.SetMode(gin.ReleaseMode)
	case config.EnvTest:
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
		if err := admin.Mount(r, cfg.AppConfig, cfg.DB); err != nil {
			slog.Warn("Failed to mount GoAdmin engine", "error", err)
		}
	}

	// Swagger Documentation UI (disabled in production)
	if appEnv != config.EnvProduction {
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
