package postgres_test

import (
	"context"
	"testing"
	"time"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/pkg/testutil"
	"ticket-box-be/internal/repository/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create_And_Find(t *testing.T) {
	tx := testutil.SetupTestDB(t)
	repo := postgres.NewUserRepository(tx)
	ctx := context.Background()

	user := &domain.User{
		Email:        "john.doe@example.com",
		PasswordHash: "hashed_secret_123",
		FullName:     "John Doe",
		Role:         domain.RoleUser,
		Status:       domain.StatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := repo.Create(ctx, user)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)

	// FindByID
	foundByID, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, foundByID.ID)
	assert.Equal(t, "john.doe@example.com", foundByID.Email)
	assert.Equal(t, "John Doe", foundByID.FullName)

	// FindByEmail (case-insensitive)
	foundByEmail, err := repo.FindByEmail(ctx, "JOHN.DOE@EXAMPLE.COM")
	require.NoError(t, err)
	assert.Equal(t, user.ID, foundByEmail.ID)

	// Find non-existent
	notFound, err := repo.FindByID(ctx, uuid.New())
	assert.ErrorIs(t, err, domain.ErrUserNotFound)
	assert.Nil(t, notFound)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	tx := testutil.SetupTestDB(t)
	repo := postgres.NewUserRepository(tx)
	ctx := context.Background()

	user1 := &domain.User{
		Email:        "dup@example.com",
		PasswordHash: "hash1",
		FullName:     "User 1",
	}
	require.NoError(t, repo.Create(ctx, user1))

	user2 := &domain.User{
		Email:        "dup@example.com",
		PasswordHash: "hash2",
		FullName:     "User 2",
	}
	err := repo.Create(ctx, user2)
	assert.ErrorIs(t, err, domain.ErrEmailAlreadyExists)
}

func TestUserRepository_Update_And_UpdatePassword(t *testing.T) {
	tx := testutil.SetupTestDB(t)
	repo := postgres.NewUserRepository(tx)
	ctx := context.Background()

	user := &domain.User{
		Email:        "jane.doe@example.com",
		PasswordHash: "initial_hash",
		FullName:     "Jane Doe",
		Role:         domain.RoleUser,
		Status:       domain.StatusActive,
	}
	require.NoError(t, repo.Create(ctx, user))

	// Update profile
	user.FullName = "Jane Updated"
	user.Status = domain.StatusInactive
	err := repo.Update(ctx, user)
	require.NoError(t, err)

	updated, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "Jane Updated", updated.FullName)
	assert.Equal(t, domain.StatusInactive, updated.Status)

	// Update password
	err = repo.UpdatePassword(ctx, user.ID, "new_secure_hash")
	require.NoError(t, err)

	updatedAfterPw, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "new_secure_hash", updatedAfterPw.PasswordHash)
}
