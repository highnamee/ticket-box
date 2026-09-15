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

func TestPasswordResetRepository_Create_Find_MarkUsed(t *testing.T) {
	tx := testutil.SetupTestDB(t)
	userRepo := postgres.NewUserRepository(tx)
	resetRepo := postgres.NewPasswordResetRepository(tx)
	ctx := context.Background()

	// 1. Create a user first for FK reference
	user := &domain.User{
		Email:        "reset.user@example.com",
		PasswordHash: "hashed_pwd",
		FullName:     "Reset User",
	}
	require.NoError(t, userRepo.Create(ctx, user))

	// 2. Create PasswordReset token
	tokenHash := "aabbccddeeff00112233445566778899"
	reset := &domain.PasswordReset{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err := resetRepo.Create(ctx, reset)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, reset.ID)

	// 3. Find by TokenHash
	found, err := resetRepo.FindByTokenHash(ctx, tokenHash)
	require.NoError(t, err)
	assert.Equal(t, reset.ID, found.ID)
	assert.Equal(t, user.ID, found.UserID)
	assert.Nil(t, found.UsedAt)

	// 4. Mark as used
	err = resetRepo.MarkAsUsed(ctx, reset.ID)
	require.NoError(t, err)

	// 5. Trying to mark as used again should return ErrResetTokenUsed
	err = resetRepo.MarkAsUsed(ctx, reset.ID)
	assert.ErrorIs(t, err, domain.ErrResetTokenUsed)
}
