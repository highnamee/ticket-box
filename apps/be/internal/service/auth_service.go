package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/pkg/hasher"
	"ticket-box-be/internal/pkg/token"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo   domain.UserRepository
	resetRepo  domain.PasswordResetRepository
	tokenMaker token.Maker
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(
	userRepo domain.UserRepository,
	resetRepo domain.PasswordResetRepository,
	tokenMaker token.Maker,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *AuthService {
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &AuthService{
		userRepo:   userRepo,
		resetRepo:  resetRepo,
		tokenMaker: tokenMaker,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// generateAuthTokens creates a fresh pair of access and refresh JWT tokens
func (s *AuthService) generateAuthTokens(user *domain.User) (*domain.AuthTokens, error) {
	accessToken, _, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Email,
		string(user.Role),
		token.TypeAccessToken,
		s.accessTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Email,
		string(user.Role),
		token.TypeRefreshToken,
		s.refreshTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &domain.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (s *AuthService) Register(ctx context.Context, email, password, fullName string) (*domain.User, *domain.AuthTokens, error) {
	cleanEmail := strings.TrimSpace(email)
	if cleanEmail == "" || password == "" || strings.TrimSpace(fullName) == "" {
		return nil, nil, domain.ErrInvalidCredentials
	}

	// Check existing user
	existing, err := s.userRepo.FindByEmail(ctx, cleanEmail)
	if err == nil && existing != nil {
		return nil, nil, domain.ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, nil, err
	}

	// Hash password
	passwordHash, err := hasher.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	user := &domain.User{
		Email:        cleanEmail,
		PasswordHash: passwordHash,
		FullName:     strings.TrimSpace(fullName),
		Role:         domain.RoleUser,
		Status:       domain.StatusActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	tokens, err := s.generateAuthTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, *domain.AuthTokens, error) {
	cleanEmail := strings.TrimSpace(email)
	user, err := s.userRepo.FindByEmail(ctx, cleanEmail)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, nil, domain.ErrInvalidCredentials
		}
		return nil, nil, err
	}

	// Check account status before running expensive bcrypt comparison
	if user.Status != domain.StatusActive {
		return nil, nil, domain.ErrUserInactive
	}

	if !hasher.CheckPassword(password, user.PasswordHash) {
		return nil, nil, domain.ErrInvalidCredentials
	}

	tokens, err := s.generateAuthTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthTokens, error) {
	claims, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	if claims.TokenType != token.TypeRefreshToken {
		return nil, domain.ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	if user.Status != domain.StatusActive {
		return nil, domain.ErrUserInactive
	}

	return s.generateAuthTokens(user)
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	cleanEmail := strings.TrimSpace(email)
	user, err := s.userRepo.FindByEmail(ctx, cleanEmail)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// Anti-enumeration: return empty string without error
			return "", nil
		}
		return "", err
	}

	rawToken, tokenHash, err := hasher.GenerateRandomToken(32)
	if err != nil {
		return "", err
	}

	reset := &domain.PasswordReset{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := s.resetRepo.Create(ctx, reset); err != nil {
		return "", err
	}

	return rawToken, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
	if strings.TrimSpace(rawToken) == "" || strings.TrimSpace(newPassword) == "" {
		return domain.ErrResetTokenNotFound
	}

	tokenHash := hasher.HashToken(rawToken)
	reset, err := s.resetRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return err
	}

	if reset.UsedAt != nil {
		return domain.ErrResetTokenUsed
	}

	if time.Now().After(reset.ExpiresAt) {
		return domain.ErrResetTokenExpired
	}

	newPasswordHash, err := hasher.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, reset.UserID, newPasswordHash); err != nil {
		return err
	}

	if err := s.resetRepo.MarkAsUsed(ctx, reset.ID); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
