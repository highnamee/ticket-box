package database

import (
	"fmt"
	"log"
	"time"

	"ticket-box-be/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// BuildDSN constructs the PostgreSQL connection string
func BuildDSN(cfg *config.Config) string {
	if cfg.DB.Password != "" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.DBName,
			cfg.DB.SSLMode,
		)
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)
}

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := BuildDSN(cfg)

	var logLevel logger.LogLevel
	switch cfg.AppEnv {
	case "test":
		logLevel = logger.Silent
	case "production":
		logLevel = logger.Warn
	default:
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if cfg.AppEnv != "test" {
		log.Println("✅ Database connected successfully")
	}

	return db, nil
}
