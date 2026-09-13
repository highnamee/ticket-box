package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Environment represents the application deployment environment
type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
	EnvTest        Environment = "test"
)

type Config struct {
	AppEnv      string
	Port        string
	EnableAdmin bool
	Admin       AdminConfig
	DB          DatabaseConfig
	CORS        CORSConfig
}

// IsProduction checks if current environment is production
func (c *Config) IsProduction() bool {
	return c.AppEnv == string(EnvProduction)
}

// IsDevelopment checks if current environment is development
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == string(EnvDevelopment)
}

// IsStaging checks if current environment is staging
func (c *Config) IsStaging() bool {
	return c.AppEnv == string(EnvStaging)
}

// IsTest checks if current environment is test
func (c *Config) IsTest() bool {
	return c.AppEnv == string(EnvTest)
}

type AdminConfig struct {
	Username string
	Password string
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
	_ = godotenv.Load(".env", "apps/be/.env")

	port := getEnv("PORT", "8080")
	appEnv := getEnv("APP_ENV", string(EnvDevelopment))

	defaultEnableAdmin := "true"
	if appEnv == string(EnvProduction) {
		defaultEnableAdmin = "false"
	}
	enableAdmin, _ := strconv.ParseBool(getEnv("ENABLE_ADMIN", defaultEnableAdmin))

	adminUser := getEnv("ADMIN_USERNAME", "admin")
	adminPass := getEnv("ADMIN_PASSWORD", "admin")

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
		AppEnv:      appEnv,
		Port:        port,
		EnableAdmin: enableAdmin,
		Admin: AdminConfig{
			Username: adminUser,
			Password: adminPass,
		},
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
