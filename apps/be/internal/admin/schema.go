package admin

import (
	_ "embed"
	"fmt"

	"gorm.io/gorm"
)

//go:embed schema.sql
var schemaSQL string

// InitSchema automatically ensures GoAdmin tables and initial seed data exist using GORM
func InitSchema(db *gorm.DB) error {
	if err := db.Exec(schemaSQL).Error; err != nil {
		return fmt.Errorf("failed to initialize GoAdmin schema: %w", err)
	}
	return nil
}
