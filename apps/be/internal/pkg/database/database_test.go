package database_test

import (
	"testing"

	"ticket-box-be/internal/config"
	"ticket-box-be/internal/pkg/database"
	"ticket-box-be/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase_BuildDSN(t *testing.T) {
	cfgWithPassword := &config.Config{
		DB: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "secretpassword",
			DBName:   "ticketbox",
			SSLMode:  "disable",
		},
	}

	dsn := database.BuildDSN(cfgWithPassword)
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "password=secretpassword")

	cfgWithoutPassword := &config.Config{
		DB: config.DatabaseConfig{
			Host:    "localhost",
			Port:    5432,
			User:    "postgres",
			DBName:  "ticketbox",
			SSLMode: "disable",
		},
	}

	dsnNoPass := database.BuildDSN(cfgWithoutPassword)
	assert.Contains(t, dsnNoPass, "host=localhost")
	assert.NotContains(t, dsnNoPass, "password=")
}

func TestDatabase_NewDatabase_Success(t *testing.T) {
	cfg := testutil.GetTestConfig(t)
	db, err := database.NewDatabase(cfg)
	require.NoError(t, err)
	require.NotNil(t, db)
}

func TestDatabase_Close(t *testing.T) {
	t.Run("nil database returns nil", func(t *testing.T) {
		err := database.Close(nil)
		assert.NoError(t, err)
	})

	t.Run("valid database closes connection", func(t *testing.T) {
		cfg := testutil.GetTestConfig(t)
		db, err := database.NewDatabase(cfg)
		require.NoError(t, err)
		require.NotNil(t, db)

		err = database.Close(db)
		assert.NoError(t, err)

		// Subsequent operations should fail because connection pool is closed
		sqlDB, err := db.DB()
		require.NoError(t, err)
		assert.Error(t, sqlDB.Ping())
	})
}
