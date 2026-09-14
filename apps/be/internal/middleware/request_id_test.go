package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ticket-box-be/internal/middleware"
	"ticket-box-be/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("generates new request ID when not provided", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.RequestIDMiddleware())

		var capturedContextReqID string
		r.GET("/test", func(c *gin.Context) {
			capturedContextReqID = logger.GetRequestID(c.Request.Context())
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		headerID := w.Header().Get(middleware.RequestIDHeader)
		require.NotEmpty(t, headerID)
		assert.Equal(t, headerID, capturedContextReqID)

		// Verify it's a valid UUID
		_, err := uuid.Parse(headerID)
		assert.NoError(t, err)
	})

	t.Run("preserves existing request ID when provided in header", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.RequestIDMiddleware())

		customReqID := "custom-trace-id-123"
		var capturedContextReqID string
		r.GET("/test", func(c *gin.Context) {
			capturedContextReqID = logger.GetRequestID(c.Request.Context())
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(middleware.RequestIDHeader, customReqID)
		r.ServeHTTP(w, req)

		assert.Equal(t, customReqID, w.Header().Get(middleware.RequestIDHeader))
		assert.Equal(t, customReqID, capturedContextReqID)
	})
}
