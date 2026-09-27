package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_BeforeCreate(t *testing.T) {
	// 1. User with no ID, Role, Status
	u1 := &User{
		Email:    "test@example.com",
		FullName: "Test User",
	}

	err := u1.BeforeCreate(nil)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, u1.ID)
	assert.Equal(t, RoleUser, u1.Role)
	assert.Equal(t, StatusActive, u1.Status)
	assert.Equal(t, "users", u1.TableName())

	// 2. User with custom ID and Role
	customID := uuid.New()
	u2 := &User{
		ID:       customID,
		Email:    "admin@example.com",
		FullName: "Admin User",
		Role:     RoleAdmin,
		Status:   StatusInactive,
	}

	err = u2.BeforeCreate(nil)
	require.NoError(t, err)
	assert.Equal(t, customID, u2.ID)
	assert.Equal(t, RoleAdmin, u2.Role)
	assert.Equal(t, StatusInactive, u2.Status)
}

func TestPasswordReset_BeforeCreate(t *testing.T) {
	pr := &PasswordReset{
		UserID:    uuid.New(),
		TokenHash: "dummyhash",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err := pr.BeforeCreate(nil)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, pr.ID)
	assert.Equal(t, "password_resets", pr.TableName())
}
