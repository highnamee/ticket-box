package testutil

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"ticket-box-be/internal/config"
	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/pkg/database"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func loadEnvFromAnywhere() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for i := 0; i < 5; i++ {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

// SetupTestDB initializes connection to PostgreSQL test database and migrates schema.
// It wraps each test in a transaction that automatically rolls back upon test completion.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	loadEnvFromAnywhere()

	dbHost := getEnv("TEST_DB_HOST", getEnv("DB_HOST", "localhost"))
	dbPortStr := getEnv("TEST_DB_PORT", getEnv("DB_PORT", "5432"))
	dbPort, _ := strconv.Atoi(dbPortStr)
	dbUser := getEnv("TEST_DB_USER", getEnv("DB_USER", "postgres"))
	dbPassword := getEnv("TEST_DB_PASSWORD", getEnv("DB_PASSWORD", ""))
	dbName := getEnv("TEST_DB_NAME", "ticketbox_test")
	dbSSLMode := getEnv("TEST_DB_SSLMODE", "disable")

	cfg := &config.Config{
		AppEnv: "test",
		Port:   "8080",
		DB: config.DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		},
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Skipf("Skipping PostgreSQL integration test (database unreachable: %v)", err)
		return nil
	}

	// Ensure schema exists
	if err := db.AutoMigrate(&domain.Ticket{}); err != nil {
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
