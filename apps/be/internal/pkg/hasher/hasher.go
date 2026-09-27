package hasher

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a bcrypt hash of the plain text password
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateRandomToken generates cryptographically secure random bytes,
// returns the hex-encoded raw token and its SHA-256 hash.
func GenerateRandomToken(byteLength int) (rawToken string, tokenHash string, err error) {
	if byteLength <= 0 {
		byteLength = 32
	}
	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	rawToken = hex.EncodeToString(bytes)
	tokenHash = HashToken(rawToken)
	return rawToken, tokenHash, nil
}

// HashToken computes the SHA-256 hash of a raw token string (hex format)
func HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}
