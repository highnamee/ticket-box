package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword_And_CheckPassword(t *testing.T) {
	password := "SecretPassword123!"

	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)

	// Valid password check
	assert.True(t, CheckPassword(password, hashed))

	// Invalid password check
	assert.False(t, CheckPassword("WrongPassword123!", hashed))
	assert.False(t, CheckPassword("", hashed))
}

func TestGenerateRandomToken(t *testing.T) {
	rawToken, tokenHash, err := GenerateRandomToken(32)
	require.NoError(t, err)
	assert.Len(t, rawToken, 64)  // 32 bytes hex encoded = 64 characters
	assert.Len(t, tokenHash, 64) // SHA-256 hex encoded = 64 characters

	// Verify HashToken matches
	computedHash := HashToken(rawToken)
	assert.Equal(t, tokenHash, computedHash)

	// Ensure consecutive tokens are different
	rawToken2, _, err := GenerateRandomToken(32)
	require.NoError(t, err)
	assert.NotEqual(t, rawToken, rawToken2)
}
