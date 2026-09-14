package testutil

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"ticket-box-be/internal/config"
	"ticket-box-be/internal/pkg/database"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func loadTestEnv() {
	// 1. Direct try .env.test from current directory or apps/be/.env.test
	if err := godotenv.Load(".env.test", "apps/be/.env.test"); err == nil {
		return
	}

	// 2. Search upwards strictly for .env.test if executed inside subpackages
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for i := 0; i < 5; i++ {
		testEnvPath := filepath.Join(dir, ".env.test")
		if _, err := os.Stat(testEnvPath); err == nil {
			_ = godotenv.Load(testEnvPath)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

// GetTestConfig builds a Config struct configured strictly for test environment and test database
func GetTestConfig(t *testing.T) *config.Config {
	t.Helper()
	loadTestEnv()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPortStr := getEnv("DB_PORT", "5432")
	dbPort, _ := strconv.Atoi(dbPortStr)
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "ticketbox_test")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	return &config.Config{
		AppEnv:      string(config.EnvTest),
		Port:        getEnv("PORT", "8080"),
		EnableAdmin: true,
		Admin: config.AdminConfig{
			Username: getEnv("ADMIN_USERNAME", "admin"),
			Password: getEnv("ADMIN_PASSWORD", "admin"),
		},
		DB: config.DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		},
	}
}

// SetupTestDB initializes connection to PostgreSQL test database and migrates schema.
// It wraps each test in a transaction that automatically rolls back upon test completion.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	cfg := GetTestConfig(t)

	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Skipf("Skipping PostgreSQL integration test (database unreachable: %v)", err)
		return nil
	}

	// Ensure schema exists via Goose versioned migrations
	if err := database.MigrateUp(db); err != nil {
		t.Fatalf("Failed to migrate test schema: %v", err)
	}

	// Begin isolated transaction for the test
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})

	return tx
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultVal
}
