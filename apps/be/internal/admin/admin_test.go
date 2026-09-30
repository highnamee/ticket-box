package admin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ticket-box-be/internal/config"
	"ticket-box-be/internal/pkg/database"
	"ticket-box-be/internal/pkg/testutil"
)

func TestGoAdmin_Routes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	cfg := testutil.GetTestConfig(t)
	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}
	if err := database.MigrateUp(db); err != nil {
		t.Fatalf("Failed to migrate test db: %v", err)
	}
	// We do NOT use a transaction (db.Begin()) here because GoAdmin uses its own connection pool
	// based on the DSN and needs to see the schema changes committed globally.
	err = Mount(r, cfg, db)
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
		ticketID := "11111111-2222-3333-4444-555555555555"
		err = db.Exec("INSERT INTO tickets (id, name, price, total_quantity, available_stock, status) VALUES (?, 'VIP Concert', 150.00, 100, 100, 'ACTIVE') ON CONFLICT (id) DO NOTHING", ticketID).Error
		assert.NoError(t, err)
		t.Cleanup(func() {
			_ = db.Exec("DELETE FROM tickets WHERE id = ?", ticketID).Error
		})

		req4, _ := http.NewRequest(http.MethodGet, "/admin/info/tickets", nil)
		req4.Header.Set("Cookie", cookie)
		w4 := httptest.NewRecorder()
		r.ServeHTTP(w4, req4)
		assert.Equal(t, http.StatusOK, w4.Code)
		assert.Contains(t, w4.Body.String(), ticketID)
		assert.Contains(t, w4.Body.String(), "__goadmin_detail_pk="+ticketID)
		assert.Contains(t, w4.Body.String(), "__goadmin_edit_pk="+ticketID)

		// Test GET /admin/info/tickets/detail with valid UUID
		reqDetail, _ := http.NewRequest(http.MethodGet, "/admin/info/tickets/detail?__goadmin_detail_pk="+ticketID, nil)
		reqDetail.Header.Set("Cookie", cookie)
		wDetail := httptest.NewRecorder()
		r.ServeHTTP(wDetail, reqDetail)
		assert.Equal(t, http.StatusOK, wDetail.Code)
		assert.NotContains(t, wDetail.Body.String(), "invalid input syntax for type uuid")

		// Test GET /admin/info/tickets/edit with valid UUID
		reqEdit, _ := http.NewRequest(http.MethodGet, "/admin/info/tickets/edit?__goadmin_edit_pk="+ticketID, nil)
		reqEdit.Header.Set("Cookie", cookie)
		wEdit := httptest.NewRecorder()
		r.ServeHTTP(wEdit, reqEdit)
		assert.Equal(t, http.StatusOK, wEdit.Code)
		assert.NotContains(t, wEdit.Body.String(), "invalid input syntax for type uuid")
	}
}

func TestInitSchema_ProductionRejectsDefaultPassword(t *testing.T) {
	prodCfg := &config.Config{
		AppEnv:      string(config.EnvProduction),
		EnableAdmin: true,
		Admin: config.AdminConfig{
			Username: "admin",
			Password: "admin",
		},
	}

	err := InitSchema(nil, prodCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refusing to initialize GoAdmin in production with default password")
}
