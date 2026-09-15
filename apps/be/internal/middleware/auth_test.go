package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/middleware"
	"ticket-box-be/internal/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthMiddlewareTestRouter(maker token.Maker) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	protected := r.Group("/protected")
	protected.Use(middleware.AuthMiddleware(maker))
	{
		protected.GET("/me", func(c *gin.Context) {
			userID, _ := middleware.GetAuthUserID(c)
			email, _ := middleware.GetAuthUserEmail(c)
			role, _ := middleware.GetAuthUserRole(c)
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID.String(),
				"email":   email,
				"role":    role,
			})
		})

		adminOnly := protected.Group("/admin")
		adminOnly.Use(middleware.RequireRole(domain.RoleAdmin))
		{
			adminOnly.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "admin_access_granted"})
			})
		}
	}

	return r
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	maker, err := token.NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	r := setupAuthMiddlewareTestRouter(maker)
	userID := uuid.New()
	tokenString, _, err := maker.CreateToken(userID, "user@test.com", "USER", token.TypeAccessToken, time.Minute)
	require.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "/protected/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_MissingOrInvalidHeader(t *testing.T) {
	maker, err := token.NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)
	r := setupAuthMiddlewareTestRouter(maker)

	// Missing header
	req, _ := http.NewRequest(http.MethodGet, "/protected/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Wrong auth scheme
	req, _ = http.NewRequest(http.MethodGet, "/protected/me", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Malformed Bearer
	req, _ = http.NewRequest(http.MethodGet, "/protected/me", nil)
	req.Header.Set("Authorization", "Bearer")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_RejectRefreshTokenAsAccess(t *testing.T) {
	maker, err := token.NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)
	r := setupAuthMiddlewareTestRouter(maker)

	refreshToken, _, err := maker.CreateToken(uuid.New(), "user@test.com", "USER", token.TypeRefreshToken, time.Minute)
	require.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "/protected/me", nil)
	req.Header.Set("Authorization", "Bearer "+refreshToken)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireRole_RBAC(t *testing.T) {
	maker, err := token.NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)
	r := setupAuthMiddlewareTestRouter(maker)

	// Regular USER accessing /protected/admin/dashboard -> 403 Forbidden
	userToken, _, err := maker.CreateToken(uuid.New(), "user@test.com", "USER", token.TypeAccessToken, time.Minute)
	require.NoError(t, err)

	req, _ := http.NewRequest(http.MethodGet, "/protected/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// ADMIN user accessing /protected/admin/dashboard -> 200 OK
	adminToken, _, err := maker.CreateToken(uuid.New(), "admin@test.com", "ADMIN", token.TypeAccessToken, time.Minute)
	require.NoError(t, err)

	req, _ = http.NewRequest(http.MethodGet, "/protected/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
