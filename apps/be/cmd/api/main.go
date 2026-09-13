package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ticket-box-be/internal/config"
	httpHandler "ticket-box-be/internal/handler/http"
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

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Initialize Handlers
	healthHandler := httpHandler.NewHealthHandler()

	// 3. Setup Router
	router := httpHandler.NewRouter(httpHandler.RouterConfig{
		AppConfig:          cfg,
		AppEnv:             cfg.AppEnv,
		HealthHandler:      healthHandler,
		CORSAllowedOrigins: cfg.CORS.AllowedOrigins,
		EnableAdmin:        cfg.EnableAdmin,
	})

	// 4. Configure HTTP Server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Start Server in a separate goroutine
	go func() {
		log.Printf("🚀 Ticket Box Backend running on http://localhost:%s (Env: %s)", cfg.Port, cfg.AppEnv)
		if cfg.EnableAdmin {
			log.Printf("🛠️  GoAdmin Panel available at http://localhost:%s/admin", cfg.Port)
		}
		if cfg.AppEnv != "production" {
			log.Printf("📖 Swagger UI available at http://localhost:%s/swagger/index.html", cfg.Port)
		}
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly.")
}
