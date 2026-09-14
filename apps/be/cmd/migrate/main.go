package main

import (
	"log/slog"
	"os"

	"ticket-box-be/internal/admin"
	"ticket-box-be/internal/config"
	"ticket-box-be/internal/pkg/database"
	"ticket-box-be/internal/pkg/logger"
)

func main() {
	cfg := config.Load()
	logger.InitLogger(cfg.AppEnv)

	db, err := database.NewDatabase(cfg)
	if err != nil {
		slog.Error("Failed to connect to database for migration", "error", err)
		os.Exit(1)
	}

	// 1. Migrate Domain Models (tickets, etc.)
	if err := database.AutoMigrate(db); err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}

	// 2. Migrate GoAdmin Internal Schema
	if err := admin.InitSchema(db, cfg); err != nil {
		slog.Error("GoAdmin schema initialization failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Migration finished successfully")
}
