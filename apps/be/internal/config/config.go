package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	Port   string
	DB     DatabaseConfig
	CORS   CORSConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() *Config {
	// Load .env file if it exists, ignore error if missing (e.g. in CI or Production)
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	appEnv := getEnv("APP_ENV", "development")

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	corsOriginsStr := getEnv("CORS_ALLOWED_ORIGINS", "")
	var allowedOrigins []string
	if corsOriginsStr != "" {
		for _, o := range strings.Split(corsOriginsStr, ",") {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}

	return &Config{
		AppEnv: appEnv,
		Port:   port,
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "ticketbox_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		CORS: CORSConfig{
			AllowedOrigins: allowedOrigins,
		},
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultVal
}
