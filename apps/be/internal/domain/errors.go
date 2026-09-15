package domain

import "errors"

var (
	ErrNotFound          = errors.New("record not found")
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrTicketSoldOut     = errors.New("ticket is sold out")
	ErrInsufficientStock = errors.New("insufficient ticket stock")
	ErrInvalidQuantity   = errors.New("invalid quantity requested")
	ErrInvalidPagination = errors.New("invalid pagination parameters")

	// Authentication & User domain errors
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailAlreadyExists   = errors.New("email already registered")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrUserInactive         = errors.New("user account is inactive or banned")
	ErrResetTokenNotFound   = errors.New("invalid reset token")
	ErrResetTokenExpired    = errors.New("reset token has expired")
	ErrResetTokenUsed       = errors.New("reset token has already been used")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidInput         = errors.New("invalid input")
)

