package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("ENABLE_ADMIN", "")
	t.Setenv("ADMIN_USERNAME", "")
	t.Setenv("ADMIN_PASSWORD", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	cfg := Load()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.AppEnv)
	assert.True(t, cfg.EnableAdmin)
	assert.Equal(t, "admin", cfg.Admin.Username)
	assert.Equal(t, "admin", cfg.Admin.Password)
	assert.Equal(t, 5432, cfg.DB.Port)
	assert.Equal(t, "localhost", cfg.DB.Host)
	assert.Equal(t, "ticketbox_db", cfg.DB.DBName)
	assert.Empty(t, cfg.CORS.AllowedOrigins)
	assert.Empty(t, cfg.JWT.Secret) // No default — must be supplied via JWT_SECRET env var
	assert.Equal(t, "15m", cfg.JWT.AccessTokenTTL)
	assert.Equal(t, "168h", cfg.JWT.RefreshTokenTTL)
}

func TestLoad_CustomEnvValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("APP_ENV", "production")
	t.Setenv("ENABLE_ADMIN", "")
	t.Setenv("ADMIN_USERNAME", "superadmin")
	t.Setenv("ADMIN_PASSWORD", "Secret@2026!")
	t.Setenv("DB_HOST", "db.internal")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "production_db")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://ticketbox.com, https://admin.ticketbox.com")
	t.Setenv("JWT_SECRET", "super-secret-key-for-production-use-only-32chars")

	cfg := Load()

	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "production", cfg.AppEnv)
	assert.False(t, cfg.EnableAdmin) // In production, defaults to false unless explicitly ENABLE_ADMIN=true
	assert.Equal(t, "superadmin", cfg.Admin.Username)
	assert.Equal(t, "Secret@2026!", cfg.Admin.Password)
	assert.Equal(t, "db.internal", cfg.DB.Host)
	assert.Equal(t, 5433, cfg.DB.Port)
	assert.Equal(t, "production_db", cfg.DB.DBName)
	assert.Equal(t, []string{"https://ticketbox.com", "https://admin.ticketbox.com"}, cfg.CORS.AllowedOrigins)
}

func TestConfig_EnvironmentHelpers(t *testing.T) {
	devCfg := &Config{AppEnv: string(EnvDevelopment)}
	assert.True(t, devCfg.IsDevelopment())
	assert.False(t, devCfg.IsProduction())
	assert.False(t, devCfg.IsStaging())
	assert.False(t, devCfg.IsTest())

	prodCfg := &Config{AppEnv: string(EnvProduction)}
	assert.True(t, prodCfg.IsProduction())
	assert.False(t, prodCfg.IsDevelopment())
	assert.False(t, prodCfg.IsStaging())
	assert.False(t, prodCfg.IsTest())

	stagingCfg := &Config{AppEnv: string(EnvStaging)}
	assert.True(t, stagingCfg.IsStaging())
	assert.False(t, stagingCfg.IsProduction())

	testCfg := &Config{AppEnv: string(EnvTest)}
	assert.True(t, testCfg.IsTest())
	assert.False(t, testCfg.IsProduction())
}
