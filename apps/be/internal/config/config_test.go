package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	cfg := Load()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, 5432, cfg.DB.Port)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, "ticketbox_db", cfg.DB.DBName)
	assert.Empty(t, cfg.CORS.AllowedOrigins)
}

func TestLoad_CustomEnvValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("APP_ENV", "production")
	t.Setenv("DB_HOST", "db.internal")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "production_db")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://ticketbox.com, https://admin.ticketbox.com")

	cfg := Load()

	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, "db.internal", cfg.DB.Host)
	assert.Equal(t, 5433, cfg.DB.Port)
	assert.Equal(t, "production_db", cfg.DB.DBName)
	assert.Equal(t, []string{"https://ticketbox.com", "https://admin.ticketbox.com"}, cfg.CORS.AllowedOrigins)
}
