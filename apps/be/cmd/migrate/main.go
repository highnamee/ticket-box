package main

import (
	"fmt"
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

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	// Subcommand 'create' does not require connecting to database
	if command == "create" {
		if len(os.Args) < 3 || os.Args[2] == "" {
			slog.Error("Migration name is required. Usage: go run ./cmd/migrate/main.go create <migration_name>")
			os.Exit(1)
		}
		name := os.Args[2]
		filePath, err := database.CreateMigration(name)
		if err != nil {
			slog.Error("Failed to create migration file", "error", err)
			os.Exit(1)
		}
		slog.Info("Created new migration file", "path", filePath)
		return
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		slog.Error("Failed to connect to database for migration", "error", err)
		os.Exit(1)
	}

	switch command {
	case "up":
		// 1. Run Goose Versioned SQL Migrations
		if err := database.MigrateUp(db); err != nil {
			slog.Error("Goose migration (up) failed", "error", err)
			os.Exit(1)
		}

		// 2. Ensure GoAdmin Internal Schema is initialized
		if err := admin.InitSchema(db, cfg); err != nil {
			slog.Error("GoAdmin schema initialization failed", "error", err)
			os.Exit(1)
		}

		slog.Info("All database migrations finished successfully")

	case "down":
		if err := database.MigrateDown(db); err != nil {
			slog.Error("Goose migration (down) failed", "error", err)
			os.Exit(1)
		}

	case "status":
		if err := database.MigrateStatus(db); err != nil {
			slog.Error("Goose migration (status) failed", "error", err)
			os.Exit(1)
		}

	default:
		slog.Error(fmt.Sprintf("Unknown migration command '%s'. Supported commands: up, down, status, create", command))
		os.Exit(1)
	}
}
