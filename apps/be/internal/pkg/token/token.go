package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

var (
	ErrInvalidToken  = errors.New("token is invalid")
	ErrExpiredToken  = errors.New("token has expired")
	ErrInvalidSecret = errors.New("invalid key size: secret key must be at least 32 characters")
)

type TokenType string

const (
	TypeAccessToken  TokenType = "access"
	TypeRefreshToken TokenType = "refresh"
)

// Claims contains custom user claims embedded inside JWT
type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// Maker is an interface for managing authentication tokens
type Maker interface {
	CreateToken(userID uuid.UUID, email, role string, tokenType TokenType, duration time.Duration) (string, *Claims, error)
	VerifyToken(tokenString string) (*Claims, error)
}

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker with minimum key size validation
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("%w (got %d chars)", ErrInvalidSecret, len(secretKey))
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a signed JWT token for a specific user and duration
func (maker *JWTMaker) CreateToken(userID uuid.UUID, email, role string, tokenType TokenType, duration time.Duration) (string, *Claims, error) {
	tokenID, err := uuid.NewV7()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token id: %w", err)
	}

	now := time.Now()
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, claims, nil
}

// VerifyToken checks if the token is valid and unexpired
func (maker *JWTMaker) VerifyToken(tokenString string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*Claims)
	if !ok || !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
