package admin

import (
	_ "embed"
	"fmt"

	"github.com/GoAdminGroup/go-admin/modules/auth"
	"gorm.io/gorm"
	"ticket-box-be/internal/config"
)

//go:embed schema.sql
var schemaSQL string

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "admin"
)

// InitSchema automatically ensures GoAdmin tables and initial admin credentials exist using GORM
func InitSchema(db *gorm.DB, appCfg *config.Config) error {
	if appCfg != nil && appCfg.IsProduction() {
		password := appCfg.Admin.Password
		if password == "" || password == defaultAdminPassword {
			return fmt.Errorf("refusing to initialize GoAdmin in production with default password; please set a secure ADMIN_PASSWORD environment variable")
		}
	}

	if err := db.Exec(schemaSQL).Error; err != nil {
		return fmt.Errorf("failed to initialize GoAdmin schema: %w", err)
	}

	// Ensure an admin user exists
	var count int64
	if err := db.Table("goadmin_users").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check admin users: %w", err)
	}

	if count == 0 {
		username := defaultAdminUsername
		password := defaultAdminPassword
		if appCfg != nil {
			if appCfg.Admin.Username != "" {
				username = appCfg.Admin.Username
			}
			if appCfg.Admin.Password != "" {
				password = appCfg.Admin.Password
			}
		}

		hashedPassword := auth.EncodePassword([]byte(password))
		insertSQL := `INSERT INTO goadmin_users (id, username, password, name, avatar, remember_token)
VALUES (1, ?, ?, 'Administrator', '', '')
ON CONFLICT (id) DO NOTHING;`
		if err := db.Exec(insertSQL, username, hashedPassword).Error; err != nil {
			return fmt.Errorf("failed to create default admin user: %w", err)
		}

		_ = db.Exec(`INSERT INTO goadmin_role_users (role_id, user_id) VALUES (1, 1) ON CONFLICT DO NOTHING;`).Error
	}

	return nil
}
