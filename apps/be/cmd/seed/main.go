package main

import (
	"context"
	"log/slog"
	"os"

	"ticket-box-be/internal/config"
	"ticket-box-be/internal/pkg/database"
	"ticket-box-be/internal/pkg/logger"
)

func main() {
	cfg := config.Load()
	logger.InitLogger(cfg.AppEnv)

	slog.Info("Connecting to database for seeding...")
	db, err := database.NewDatabase(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// 1. Ensure schema is migrated up
	slog.Info("Ensuring database migrations are up to date...")
	if err := database.MigrateUp(db); err != nil {
		slog.Error("Failed to run database migrations", "error", err)
		os.Exit(1)
	}

	// 2. Run Ticket Seeder
	slog.Info("Seeding tickets matching frontend catalog...")
	ctx := context.Background()
	if err := database.SeedTickets(ctx, db); err != nil {
		slog.Error("Failed to seed tickets", "error", err)
		os.Exit(1)
	}

	slog.Info("✅ Seeding completed successfully!")
}
