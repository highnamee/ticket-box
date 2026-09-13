package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCORSMiddleware_DevelopmentMode(t *testing.T) {
	r := gin.New()
	r.Use(CORSMiddleware(nil)) // Empty means dev mode
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Authorization")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORSMiddleware_ProductionWhitelist_Allowed(t *testing.T) {
	whitelist := []string{"https://ticketbox.com", "https://admin.ticketbox.com"}
	r := gin.New()
	r.Use(CORSMiddleware(whitelist))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://ticketbox.com")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "https://ticketbox.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_ProductionWhitelist_Disallowed(t *testing.T) {
	whitelist := []string{"https://ticketbox.com"}
	r := gin.New()
	r.Use(CORSMiddleware(whitelist))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://malicious-site.com")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// gin-contrib/cors rejects disallowed origins with 403 Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}
