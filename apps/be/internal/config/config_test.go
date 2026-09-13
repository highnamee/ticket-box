package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear relevant environment variables
	os.Unsetenv("PORT")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("CORS_ALLOWED_ORIGINS")

	cfg := Load()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, 5432, cfg.DB.Port)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, "ticketbox_db", cfg.DB.DBName)
	assert.Empty(t, cfg.CORS.AllowedOrigins)
}

func TestLoad_CustomEnvValues(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("APP_ENV", "production")
	os.Setenv("DB_HOST", "db.internal")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_NAME", "production_db")
	os.Setenv("CORS_ALLOWED_ORIGINS", "https://ticketbox.com, https://admin.ticketbox.com")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("APP_ENV")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
	}()

	cfg := Load()

	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, "db.internal", cfg.DB.Host)
	assert.Equal(t, 5433, cfg.DB.Port)
	assert.Equal(t, "production_db", cfg.DB.DBName)
	assert.Equal(t, []string{"https://ticketbox.com", "https://admin.ticketbox.com"}, cfg.CORS.AllowedOrigins)
}
