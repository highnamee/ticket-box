package admin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ticket-box-be/internal/config"
	"ticket-box-be/internal/pkg/testutil"
)

func TestGoAdmin_Routes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db := testutil.SetupTestDB(t)
	if db == nil {
		return
	}

	cfg := config.Load()
	err := Mount(r, cfg)
	assert.NoError(t, err)

	for _, route := range r.Routes() {
		t.Logf("Route: %-6s %s", route.Method, route.Path)
	}

	// Test GET /admin/login
	req, _ := http.NewRequest(http.MethodGet, "/admin/login", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test GET /admin (redirects to /admin/info/tickets)
	reqAdmin, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	wAdmin := httptest.NewRecorder()
	r.ServeHTTP(wAdmin, reqAdmin)
	assert.Equal(t, http.StatusFound, wAdmin.Code)
	assert.Equal(t, "/admin/info/tickets", wAdmin.Header().Get("Location"))

	// Test POST /admin/signin with credentials
	formData := "username=admin&password=admin"
	req3, _ := http.NewRequest(http.MethodPost, "/admin/signin", strings.NewReader(formData))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	// Extract cookie and test authenticated dashboard route
	cookie := w3.Header().Get("Set-Cookie")
	if cookie != "" {
		req4, _ := http.NewRequest(http.MethodGet, "/admin/info/tickets", nil)
		req4.Header.Set("Cookie", cookie)
		w4 := httptest.NewRecorder()
		r.ServeHTTP(w4, req4)
		assert.Equal(t, http.StatusOK, w4.Code)
	}
}
