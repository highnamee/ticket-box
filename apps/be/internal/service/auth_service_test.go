package service_test

import (
	"context"
	"testing"
	"time"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/pkg/hasher"
	"ticket-box-be/internal/pkg/token"
	"ticket-box-be/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	args := m.Called(ctx, id, passwordHash)
	return args.Error(0)
}

// MockPasswordResetRepository
type MockPasswordResetRepository struct {
	mock.Mock
}

func (m *MockPasswordResetRepository) Create(ctx context.Context, reset *domain.PasswordReset) error {
	args := m.Called(ctx, reset)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.PasswordReset, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PasswordReset), args.Error(1)
}

func (m *MockPasswordResetRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupAuthService(t *testing.T) (*service.AuthService, *MockUserRepository, *MockPasswordResetRepository, token.Maker) {
	userRepo := new(MockUserRepository)
	resetRepo := new(MockPasswordResetRepository)
	maker, err := token.NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	svc := service.NewAuthService(userRepo, resetRepo, maker, 15*time.Minute, 7*24*time.Hour)
	return svc, userRepo, resetRepo, maker
}

func TestAuthService_Register_Success(t *testing.T) {
	svc, userRepo, _, _ := setupAuthService(t)
	ctx := context.Background()

	userRepo.On("FindByEmail", ctx, "newuser@example.com").Return(nil, domain.ErrUserNotFound)
	userRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	user, tokens, err := svc.Register(ctx, "newuser@example.com", "SecurePass123!", "New User")
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser@example.com", user.Email)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	svc, userRepo, _, _ := setupAuthService(t)
	ctx := context.Background()

	existingUser := &domain.User{Email: "existing@example.com"}
	userRepo.On("FindByEmail", ctx, "existing@example.com").Return(existingUser, nil)

	user, tokens, err := svc.Register(ctx, "existing@example.com", "SecurePass123!", "New User")
	assert.ErrorIs(t, err, domain.ErrEmailAlreadyExists)
	assert.Nil(t, user)
	assert.Nil(t, tokens)
}

func TestAuthService_Login_Success(t *testing.T) {
	svc, userRepo, _, _ := setupAuthService(t)
	ctx := context.Background()

	hashed, _ := hasher.HashPassword("Password123!")
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "user@example.com",
		PasswordHash: hashed,
		Role:         domain.RoleUser,
		Status:       domain.StatusActive,
	}

	userRepo.On("FindByEmail", ctx, "user@example.com").Return(user, nil)

	loggedInUser, tokens, err := svc.Login(ctx, "user@example.com", "Password123!")
	require.NoError(t, err)
	assert.Equal(t, user.ID, loggedInUser.ID)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	svc, userRepo, _, _ := setupAuthService(t)
	ctx := context.Background()

	hashed, _ := hasher.HashPassword("CorrectPassword123!")
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "user@example.com",
		PasswordHash: hashed,
		Status:       domain.StatusActive,
	}

	userRepo.On("FindByEmail", ctx, "user@example.com").Return(user, nil)

	loggedInUser, tokens, err := svc.Login(ctx, "user@example.com", "WrongPassword")
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	assert.Nil(t, loggedInUser)
	assert.Nil(t, tokens)
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	svc, userRepo, _, maker := setupAuthService(t)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:     userID,
		Email:  "user@example.com",
		Role:   domain.RoleUser,
		Status: domain.StatusActive,
	}

	refreshToken, _, err := maker.CreateToken(userID, "user@example.com", "USER", token.TypeRefreshToken, time.Hour)
	require.NoError(t, err)

	userRepo.On("FindByID", ctx, userID).Return(user, nil)

	tokens, err := svc.RefreshToken(ctx, refreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
}

func TestAuthService_ForgotPassword_And_ResetPassword(t *testing.T) {
	svc, userRepo, resetRepo, _ := setupAuthService(t)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:     userID,
		Email:  "forgot@example.com",
		Status: domain.StatusActive,
	}

	userRepo.On("FindByEmail", ctx, "forgot@example.com").Return(user, nil)
	resetRepo.On("Create", ctx, mock.AnythingOfType("*domain.PasswordReset")).Return(nil)

	rawToken, err := svc.ForgotPassword(ctx, "forgot@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, rawToken)

	// Reset password flow
	resetID := uuid.New()
	tokenHash := hasher.HashToken(rawToken)
	resetRecord := &domain.PasswordReset{
		ID:        resetID,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		UsedAt:    nil,
	}

	resetRepo.On("FindByTokenHash", ctx, tokenHash).Return(resetRecord, nil)
	userRepo.On("UpdatePassword", ctx, userID, mock.AnythingOfType("string")).Return(nil)
	resetRepo.On("MarkAsUsed", ctx, resetID).Return(nil)

	err = svc.ResetPassword(ctx, rawToken, "NewSecurePassword123!")
	require.NoError(t, err)

	userRepo.AssertExpectations(t)
	resetRepo.AssertExpectations(t)
}
