package http

import (
	"time"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
)

// RegisterRequest defines the payload for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8,max=72" example:"SecurePassword123!"`
	FullName string `json:"full_name" binding:"required,notblank,max=100" example:"Nguyen Van A"`
}

// LoginRequest defines the payload for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePassword123!"`
}

// RefreshTokenRequest defines the payload for refreshing access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ForgotPasswordRequest defines the payload for initiating password reset
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// ResetPasswordRequest defines the payload for completing password reset
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required" example:"3a7f8e9b..."`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72" example:"NewSecurePassword123!"`
}

// UserResponse represents the public user profile without sensitive fields
type UserResponse struct {
	ID        uuid.UUID         `json:"id" example:"01956550-9d3e-7a4c-81b4-2e91cf98f01b"`
	Email     string            `json:"email" example:"user@example.com"`
	FullName  string            `json:"full_name" example:"Nguyen Van A"`
	Role      domain.UserRole   `json:"role" example:"USER"`
	Status    domain.UserStatus `json:"status" example:"ACTIVE"`
	CreatedAt time.Time         `json:"created_at"`
}

// AuthTokenResponse represents issued JWT access and refresh tokens
type AuthTokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int64  `json:"expires_in" example:"900"`
	TokenType    string `json:"token_type" example:"Bearer"`
}

// AuthResponse represents the complete authentication response
type AuthResponse struct {
	User   UserResponse      `json:"user"`
	Tokens AuthTokenResponse `json:"tokens"`
}

func toUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}

func toAuthTokenResponse(t *domain.AuthTokens) AuthTokenResponse {
	return AuthTokenResponse{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresIn:    t.ExpiresIn,
		TokenType:    t.TokenType,
	}
}
