package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"ticket-box-be/migrations"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func init() {
	goose.SetBaseFS(migrations.EmbedFS)
	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("Failed to set goose postgres dialect", "error", err)
	}
}

// MigrateUp executes all pending migrations
func MigrateUp(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return MigrateUpSQL(sqlDB)
}

// MigrateUpSQL executes all pending migrations using *sql.DB
func MigrateUpSQL(sqlDB *sql.DB) error {
	slog.Info("Running database migrations (Up)...")
	if err := goose.Up(sqlDB, "."); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}
	slog.Info("Database migrations (Up) completed successfully")
	return nil
}

// MigrateDown rolls back the most recent migration
func MigrateDown(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return MigrateDownSQL(sqlDB)
}

// MigrateDownSQL rolls back the most recent migration using *sql.DB
func MigrateDownSQL(sqlDB *sql.DB) error {
	slog.Info("Rolling back database migration (Down)...")
	if err := goose.Down(sqlDB, "."); err != nil {
		return fmt.Errorf("goose down failed: %w", err)
	}
	slog.Info("Database rollback (Down) completed successfully")
	return nil
}

// MigrateStatus displays the status of all migrations
func MigrateStatus(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return goose.Status(sqlDB, ".")
}

// CreateMigration generates a new sequential SQL migration file in the migrations folder
func CreateMigration(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("migration name cannot be empty")
	}

	// Clean name: replace spaces and dashes with underscores
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	// Determine migrations directory
	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = "apps/be/migrations"
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			return "", fmt.Errorf("migrations directory not found")
		}
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return "", fmt.Errorf("failed to read migrations directory: %w", err)
	}

	highestSeq := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) > 0 {
			if seq, err := strconv.Atoi(parts[0]); err == nil && seq > highestSeq {
				highestSeq = seq
			}
		}
	}

	nextSeq := highestSeq + 1
	filename := fmt.Sprintf("%05d_%s.sql", nextSeq, name)
	targetPath := filepath.Join(migrationsDir, filename)

	content := `-- +goose Up
-- +goose StatementBegin
-- Write your migration UP SQL here
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Write your migration DOWN SQL here
-- +goose StatementEnd
`

	if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write migration file: %w", err)
	}

	return targetPath, nil
}
