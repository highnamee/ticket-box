package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecretKey = "12345678901234567890123456789012" // 32 chars

func TestJWTMaker_KeySizeValidation(t *testing.T) {
	_, err := NewJWTMaker("too_short_key")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidSecret)

	maker, err := NewJWTMaker(testSecretKey)
	assert.NoError(t, err)
	assert.NotNil(t, maker)
}

func TestJWTMaker_CreateAndVerifyToken_Success(t *testing.T) {
	maker, err := NewJWTMaker(testSecretKey)
	require.NoError(t, err)

	userID := uuid.New()
	email := "test@example.com"
	role := "user"
	duration := time.Minute

	tokenString, claims, err := maker.CreateToken(userID, email, role, TypeAccessToken, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotNil(t, claims)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, TypeAccessToken, claims.TokenType)

	// Verify Token
	verifiedClaims, err := maker.VerifyToken(tokenString)
	require.NoError(t, err)
	require.NotNil(t, verifiedClaims)

	assert.Equal(t, userID, verifiedClaims.UserID)
	assert.Equal(t, email, verifiedClaims.Email)
	assert.Equal(t, role, verifiedClaims.Role)
	assert.Equal(t, TypeAccessToken, verifiedClaims.TokenType)
	assert.WithinDuration(t, time.Now().Add(duration), verifiedClaims.ExpiresAt.Time, 5*time.Second)
}

func TestJWTMaker_ExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(testSecretKey)
	require.NoError(t, err)

	userID := uuid.New()
	email := "expired@example.com"
	role := "user"

	tokenString, _, err := maker.CreateToken(userID, email, role, TypeAccessToken, -time.Minute)
	require.NoError(t, err)

	verifiedClaims, err := maker.VerifyToken(tokenString)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrExpiredToken)
	assert.Nil(t, verifiedClaims)
}

func TestJWTMaker_InvalidToken(t *testing.T) {
	maker, err := NewJWTMaker(testSecretKey)
	require.NoError(t, err)

	// Tampered / invalid token string
	verifiedClaims, err := maker.VerifyToken("invalid.jwt.token.format")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidToken)
	assert.Nil(t, verifiedClaims)

	// Token created with different secret key
	maker2, err := NewJWTMaker("different_secret_key_123456789012")
	require.NoError(t, err)

	tokenString, _, err := maker2.CreateToken(uuid.New(), "user@test.com", "user", TypeAccessToken, time.Minute)
	require.NoError(t, err)

	verifiedClaims, err = maker.VerifyToken(tokenString)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidToken)
	assert.Nil(t, verifiedClaims)
}
