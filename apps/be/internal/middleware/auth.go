package middleware

import (
	"errors"
	"net/http"
	"strings"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/pkg/response"
	"ticket-box-be/internal/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationTypeBearer = "bearer"
	ContextKeyUserID        = "auth_user_id"
	ContextKeyEmail         = "auth_user_email"
	ContextKeyRole          = "auth_user_role"
)

// AuthMiddleware creates a Gin middleware for verifying JWT Access Tokens
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is not provided")
			response.Error(c, http.StatusUnauthorized, "Missing authorization header", err)
			c.Abort()
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			response.Error(c, http.StatusUnauthorized, "Invalid authorization header format", err)
			c.Abort()
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != AuthorizationTypeBearer {
			err := errors.New("unsupported authorization type")
			response.Error(c, http.StatusUnauthorized, "Unsupported authorization type, expected Bearer", err)
			c.Abort()
			return
		}

		accessToken := fields[1]
		claims, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired access token", err)
			c.Abort()
			return
		}

		if claims.TokenType != token.TypeAccessToken {
			response.Error(c, http.StatusUnauthorized, "Invalid token type, access token required", domain.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyRole, domain.UserRole(claims.Role))
		c.Next()
	}
}

// RequireRole enforces role-based access control for protected routes
func RequireRole(roles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, ok := GetAuthUserRole(c)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Unauthorized", domain.ErrUnauthorized)
			c.Abort()
			return
		}

		for _, allowedRole := range roles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "You do not have permission to access this resource", domain.ErrForbidden)
		c.Abort()
	}
}

// GetAuthUserID extracts the authenticated user ID from Gin context
func GetAuthUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(ContextKeyUserID)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

// GetAuthUserEmail extracts the authenticated user email from Gin context
func GetAuthUserEmail(c *gin.Context) (string, bool) {
	val, exists := c.Get(ContextKeyEmail)
	if !exists {
		return "", false
	}
	email, ok := val.(string)
	return email, ok
}

// GetAuthUserRole extracts the authenticated user role from Gin context
func GetAuthUserRole(c *gin.Context) (domain.UserRole, bool) {
	val, exists := c.Get(ContextKeyRole)
	if !exists {
		return "", false
	}
	role, ok := val.(domain.UserRole)
	return role, ok
}
