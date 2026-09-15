package http

import (
	"errors"
	"net/http"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/middleware"
	"ticket-box-be/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary      Register a new user account
// @Description  Create a new user account with email, password, and full name. Returns user profile and JWT auth tokens.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "User registration details"
// @Success      201 {object} response.APIResponse{data=AuthResponse} "Account created successfully"
// @Failure      400 {object} response.APIResponse "Invalid request payload or validation error"
// @Failure      409 {object} response.APIResponse "Email already registered"
// @Failure      500 {object} response.APIResponse "Internal server error"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid registration request", err)
		return
	}

	user, tokens, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Error(c, http.StatusConflict, "Email is already registered", err)
			return
		}
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.BadRequest(c, "Invalid input data", err)
			return
		}
		response.InternalServerError(c, "Failed to register user", err)
		return
	}

	res := AuthResponse{
		User:   toUserResponse(user),
		Tokens: toAuthTokenResponse(tokens),
	}
	response.Success(c, http.StatusCreated, "User registered successfully", res)
}

// Login godoc
// @Summary      Authenticate user and generate JWT tokens
// @Description  Verify email and password credentials, returning user profile and access/refresh tokens.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} response.APIResponse{data=AuthResponse} "Login successful"
// @Failure      400 {object} response.APIResponse "Invalid request payload"
// @Failure      401 {object} response.APIResponse "Invalid email or password"
// @Failure      403 {object} response.APIResponse "Account inactive or banned"
// @Failure      500 {object} response.APIResponse "Internal server error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid login request", err)
		return
	}

	user, tokens, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "Invalid email or password", err)
			return
		}
		if errors.Is(err, domain.ErrUserInactive) {
			response.Error(c, http.StatusForbidden, "User account is inactive or banned", err)
			return
		}
		response.InternalServerError(c, "Failed to authenticate user", err)
		return
	}

	res := AuthResponse{
		User:   toUserResponse(user),
		Tokens: toAuthTokenResponse(tokens),
	}
	response.Success(c, http.StatusOK, "Login successful", res)
}

// RefreshToken godoc
// @Summary      Refresh expired access token
// @Description  Issue a new pair of JWT access/refresh tokens using a valid unexpired refresh token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Refresh token payload"
// @Success      200 {object} response.APIResponse{data=AuthTokenResponse} "Token refreshed successfully"
// @Failure      400 {object} response.APIResponse "Invalid request payload"
// @Failure      401 {object} response.APIResponse "Invalid or expired refresh token"
// @Failure      500 {object} response.APIResponse "Internal server error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid refresh token request", err)
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidRefreshToken) {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired refresh token", err)
			return
		}
		if errors.Is(err, domain.ErrUserInactive) {
			response.Error(c, http.StatusForbidden, "User account is inactive or banned", err)
			return
		}
		response.InternalServerError(c, "Failed to refresh token", err)
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", toAuthTokenResponse(tokens))
}

// ForgotPassword godoc
// @Summary      Request password reset
// @Description  Initiates password reset process. If email exists, creates single-use reset token (15m validity). Always returns 200 to prevent user enumeration.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body ForgotPasswordRequest true "Email address for password reset"
// @Success      200 {object} response.APIResponse "Instructions sent if account exists"
// @Failure      400 {object} response.APIResponse "Invalid email format"
// @Failure      500 {object} response.APIResponse "Internal server error"
// @Router       /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid email address", err)
		return
	}

	_, err := h.authService.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		response.InternalServerError(c, "Failed to process forgot password request", err)
		return
	}

	// Always respond with generic success message to prevent user enumeration
	response.Success(c, http.StatusOK, "If this email exists in our system, password reset instructions have been sent.", nil)
}

// ResetPassword godoc
// @Summary      Reset password with token
// @Description  Validates single-use reset token and updates the user's password.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body ResetPasswordRequest true "Reset token and new password"
// @Success      200 {object} response.APIResponse "Password reset successful"
// @Failure      400 {object} response.APIResponse "Invalid or expired reset token, or weak password"
// @Failure      500 {object} response.APIResponse "Internal server error"
// @Router       /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid reset password request", err)
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		if errors.Is(err, domain.ErrResetTokenNotFound) || errors.Is(err, domain.ErrResetTokenExpired) || errors.Is(err, domain.ErrResetTokenUsed) {
			response.BadRequest(c, "Invalid, expired, or already used reset token", err)
			return
		}
		response.InternalServerError(c, "Failed to reset password", err)
		return
	}

	response.Success(c, http.StatusOK, "Password has been reset successfully. Please log in with your new password.", nil)
}

// GetMe godoc
// @Summary      Get current authenticated user profile
// @Description  Fetch profile information for the currently authenticated user based on JWT Bearer token.
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.APIResponse{data=UserResponse} "User profile retrieved successfully"
// @Failure      401 {object} response.APIResponse "Unauthorized or missing token"
// @Failure      404 {object} response.APIResponse "User not found"
// @Failure      500 {object} response.APIResponse "Internal server error"
// @Router       /users/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", domain.ErrUnauthorized)
		return
	}

	user, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(c, "User not found", err)
			return
		}
		response.InternalServerError(c, "Failed to retrieve user profile", err)
		return
	}

	response.Success(c, http.StatusOK, "User profile retrieved successfully", toUserResponse(user))
}
