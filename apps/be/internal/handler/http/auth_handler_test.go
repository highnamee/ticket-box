package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ticket-box-be/internal/domain"
	httpHandler "ticket-box-be/internal/handler/http"
	"ticket-box-be/internal/middleware"
	"ticket-box-be/internal/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, email, password, fullName string) (*domain.User, *domain.AuthTokens, error) {
	args := m.Called(ctx, email, password, fullName)
	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User) //nolint:errcheck
	}
	var tokens *domain.AuthTokens
	if args.Get(1) != nil {
		tokens = args.Get(1).(*domain.AuthTokens) //nolint:errcheck
	}
	return user, tokens, args.Error(2)
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) (*domain.User, *domain.AuthTokens, error) {
	args := m.Called(ctx, email, password)
	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User) //nolint:errcheck
	}
	var tokens *domain.AuthTokens
	if args.Get(1) != nil {
		tokens = args.Get(1).(*domain.AuthTokens) //nolint:errcheck
	}
	return user, tokens, args.Error(2)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthTokens, error) {
	args := m.Called(ctx, refreshToken)
	var tokens *domain.AuthTokens
	if args.Get(0) != nil {
		tokens = args.Get(0).(*domain.AuthTokens) //nolint:errcheck
	}
	return tokens, args.Error(1)
}

func (m *MockAuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
	args := m.Called(ctx, rawToken, newPassword)
	return args.Error(0)
}

func (m *MockAuthService) GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, userID)
	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User) //nolint:errcheck
	}
	return user, args.Error(1)
}

func setupAuthRouter(mockSvc *MockAuthService, tokenMaker token.Maker) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := httpHandler.NewAuthHandler(mockSvc)

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", handler.Register)
			authGroup.POST("/login", handler.Login)
			authGroup.POST("/refresh", handler.RefreshToken)
			authGroup.POST("/forgot-password", handler.ForgotPassword)
			authGroup.POST("/reset-password", handler.ResetPassword)
		}

		userGroup := v1.Group("/users")
		if tokenMaker != nil {
			userGroup.Use(middleware.AuthMiddleware(tokenMaker))
		}
		{
			userGroup.GET("/me", handler.GetMe)
		}
	}
	return r
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockSvc := new(MockAuthService)
	r := setupAuthRouter(mockSvc, nil)

	reqBody := httpHandler.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "Password123!",
		FullName: "New User",
	}
	user := &domain.User{
		ID:        uuid.New(),
		Email:     reqBody.Email,
		FullName:  reqBody.FullName,
		Role:      domain.RoleUser,
		Status:    domain.StatusActive,
		CreatedAt: time.Now(),
	}
	tokens := &domain.AuthTokens{
		AccessToken:  "access_token_xyz",
		RefreshToken: "refresh_token_xyz",
		ExpiresIn:    900,
		TokenType:    "Bearer",
	}

	mockSvc.On("Register", mock.Anything, reqBody.Email, reqBody.Password, reqBody.FullName).
		Return(user, tokens, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAuthHandler_Register_EmailConflict(t *testing.T) {
	mockSvc := new(MockAuthService)
	r := setupAuthRouter(mockSvc, nil)

	reqBody := httpHandler.RegisterRequest{
		Email:    "existing@example.com",
		Password: "Password123!",
		FullName: "Existing User",
	}

	mockSvc.On("Register", mock.Anything, reqBody.Email, reqBody.Password, reqBody.FullName).
		Return(nil, nil, domain.ErrEmailAlreadyExists)

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockSvc := new(MockAuthService)
	r := setupAuthRouter(mockSvc, nil)

	reqBody := httpHandler.LoginRequest{
		Email:    "user@example.com",
		Password: "Password123!",
	}
	user := &domain.User{
		ID:       uuid.New(),
		Email:    reqBody.Email,
		FullName: "Test User",
		Role:     domain.RoleUser,
		Status:   domain.StatusActive,
	}
	tokens := &domain.AuthTokens{
		AccessToken:  "access_token_xyz",
		RefreshToken: "refresh_token_xyz",
		ExpiresIn:    900,
		TokenType:    "Bearer",
	}

	mockSvc.On("Login", mock.Anything, reqBody.Email, reqBody.Password).
		Return(user, tokens, nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_ForgotPassword_Always200(t *testing.T) {
	mockSvc := new(MockAuthService)
	r := setupAuthRouter(mockSvc, nil)

	reqBody := httpHandler.ForgotPasswordRequest{
		Email: "random@example.com",
	}

	mockSvc.On("ForgotPassword", mock.Anything, reqBody.Email).
		Return("", nil)

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_GetMe_Success(t *testing.T) {
	mockSvc := new(MockAuthService)
	maker, err := token.NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)
	r := setupAuthRouter(mockSvc, maker)

	userID := uuid.New()
	user := &domain.User{
		ID:       userID,
		Email:    "me@example.com",
		FullName: "Current User",
		Role:     domain.RoleUser,
		Status:   domain.StatusActive,
	}

	tokenString, _, err := maker.CreateToken(userID, user.Email, "USER", token.TypeAccessToken, time.Minute)
	require.NoError(t, err)

	mockSvc.On("GetMe", mock.Anything, userID).Return(user, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
