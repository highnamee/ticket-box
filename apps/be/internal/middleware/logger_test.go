package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"ticket-box-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))

	r := gin.New()
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RequestLoggerMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/error", func(c *gin.Context) {
		c.String(http.StatusInternalServerError, "server error")
	})

	t.Run("logs 200 OK request with attributes", func(t *testing.T) {
		buf.Reset()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		logOutput := buf.String()
		assert.Contains(t, logOutput, `"method":"GET"`)
		assert.Contains(t, logOutput, `"path":"/ping"`)
		assert.Contains(t, logOutput, `"status":200`)
		assert.Contains(t, logOutput, `"request_id":`)
	})

	t.Run("logs 500 Error request at ERROR level", func(t *testing.T) {
		buf.Reset()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		logOutput := buf.String()
		assert.Contains(t, logOutput, `"level":"ERROR"`)
		assert.Contains(t, logOutput, `"path":"/error"`)
		assert.Contains(t, logOutput, `"status":500`)
	})
}
