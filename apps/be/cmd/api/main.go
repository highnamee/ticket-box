package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ticket-box-be/internal/config"
	httpHandler "ticket-box-be/internal/handler/http"
	"ticket-box-be/internal/pkg/database"
	"ticket-box-be/internal/pkg/logger"
	"ticket-box-be/internal/pkg/token"
	"ticket-box-be/internal/repository/postgres"
	"ticket-box-be/internal/service"
)

// @title                      Ticket Box Backend API
// @version                    1.0
// @description                RESTful API documentation for the Ticket Box event booking platform.
// @contact.name               Ticket Box Team
// @host                       localhost:8080
// @BasePath                   /api/v1
// @schemes                    http https
// @produce                    json
// @consume                    json
//
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type 'Bearer ' followed by your JWT access token.

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Initialize Structured Logger
	logger.InitLogger(cfg.AppEnv)

	// 3. Initialize Database Connection
	db, err := database.NewDatabase(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// 4. Initialize JWT Token Maker
	tokenMaker, err := token.NewJWTMaker(cfg.JWT.Secret)
	if err != nil {
		slog.Error("Failed to initialize JWT maker", "error", err)
		os.Exit(1)
	}

	// 5. Initialize Repositories
	ticketRepo := postgres.NewTicketRepository(db)
	userRepo := postgres.NewUserRepository(db)
	resetRepo := postgres.NewPasswordResetRepository(db)

	// 6. Initialize Services
	ticketService := service.NewTicketService(ticketRepo)
	authService := service.NewAuthService(userRepo, resetRepo, tokenMaker, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)

	// 7. Initialize Handlers
	healthHandler := httpHandler.NewHealthHandler()
	ticketHandler := httpHandler.NewTicketHandler(ticketService)
	authHandler := httpHandler.NewAuthHandler(authService)

	// 8. Setup Router
	router := httpHandler.NewRouter(httpHandler.RouterConfig{
		AppConfig:          cfg,
		AppEnv:             cfg.AppEnv,
		HealthHandler:      healthHandler,
		TicketHandler:      ticketHandler,
		AuthHandler:        authHandler,
		TokenMaker:         tokenMaker,
		CORSAllowedOrigins: cfg.CORS.AllowedOrigins,
		EnableAdmin:        cfg.EnableAdmin,
		DB:                 db,
	})

	// 9. Configure HTTP Server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 9. Start Server in a separate goroutine
	go func() {
		slog.Info("Ticket Box Backend starting",
			"port", cfg.Port,
			"env", cfg.AppEnv,
			"url", fmt.Sprintf("http://localhost:%s", cfg.Port),
		)
		if cfg.EnableAdmin {
			slog.Info("GoAdmin Panel available", "url", fmt.Sprintf("http://localhost:%s/admin", cfg.Port))
		}
		if !cfg.IsProduction() {
			slog.Info("Swagger UI available", "url", fmt.Sprintf("http://localhost:%s/swagger/index.html", cfg.Port))
		}
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// 10. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	// 11. Close Database Connection Pool
	slog.Info("Closing database connection pool...")
	if err := database.Close(db); err != nil {
		slog.Error("Failed to close database connection pool", "error", err)
	} else {
		slog.Info("Database connection pool closed successfully")
	}

	slog.Info("Server exited cleanly.")
}
