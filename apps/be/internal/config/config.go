package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
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
	AppEnv      string         `env:"APP_ENV" envDefault:"development"`
	Port        string         `env:"PORT" envDefault:"8080"`
	EnableAdmin bool           `env:"ENABLE_ADMIN" envDefault:"true"`
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
	Username string `env:"ADMIN_USERNAME" envDefault:"admin"`
	Password string `env:"ADMIN_PASSWORD" envDefault:"admin"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:""`
	DBName   string `env:"DB_NAME" envDefault:"ticketbox_db"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

type CORSConfig struct {
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:","`
}

func Load() *Config {
	_ = godotenv.Load(".env", "apps/be/.env")

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		slog.Error("Failed to parse application configuration", "error", err)
		os.Exit(1)
	}

	// Trim whitespace from CORS allowed origins
	var trimmedOrigins []string
	for _, origin := range cfg.CORS.AllowedOrigins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			trimmedOrigins = append(trimmedOrigins, trimmed)
		}
	}
	cfg.CORS.AllowedOrigins = trimmedOrigins

	// In production, default EnableAdmin to false unless explicitly configured
	if cfg.AppEnv == string(EnvProduction) && os.Getenv("ENABLE_ADMIN") == "" {
		cfg.EnableAdmin = false
	}

	return &cfg
}

