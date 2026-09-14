package database

import (
	"fmt"
	"log/slog"

	"ticket-box-be/internal/domain"

	"gorm.io/gorm"
)

// AutoMigrate runs database migrations for all domain models
func AutoMigrate(db *gorm.DB) error {
	slog.Info("Running database migrations...")

	models := []interface{}{
		&domain.Ticket{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("migration failed for model %T: %w", model, err)
		}
	}

	slog.Info("Database migrations completed successfully")
	return nil
}
