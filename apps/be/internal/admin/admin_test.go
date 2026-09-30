package admin

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"ticket-box-be/internal/config"
	"ticket-box-be/internal/pkg/database"
	"ticket-box-be/internal/pkg/testutil"
)

func setupAdminTest(t *testing.T) (*gin.Engine, *gorm.DB, string) {
	t.Helper()
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

	err = Mount(r, cfg, db)
	assert.NoError(t, err)

	formData := "username=admin&password=admin"
	req, _ := http.NewRequest(http.MethodPost, "/admin/signin", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	cookie := w.Header().Get("Set-Cookie")
	assert.NotEmpty(t, cookie)
	return r, db, cookie
}

func extractCSRFToken(body string) string {
	tokenKey := `name="__go_admin_t_" value='`
	if idx := strings.Index(body, tokenKey); idx != -1 {
		token := body[idx+len(tokenKey):]
		if endIdx := strings.Index(token, `'`); endIdx != -1 {
			return token[:endIdx]
		}
	}
	return ""
}

func TestGoAdmin(t *testing.T) {
	r, db, cookie := setupAdminTest(t)

	t.Run("Auth and Navigation", func(t *testing.T) {
		reqLogin, _ := http.NewRequest(http.MethodGet, "/admin/login", nil)
		wLogin := httptest.NewRecorder()
		r.ServeHTTP(wLogin, reqLogin)
		assert.Equal(t, http.StatusOK, wLogin.Code)

		reqAdmin, _ := http.NewRequest(http.MethodGet, "/admin", nil)
		wAdmin := httptest.NewRecorder()
		r.ServeHTTP(wAdmin, reqAdmin)
		assert.Equal(t, http.StatusFound, wAdmin.Code)
		assert.Equal(t, "/admin/info/tickets", wAdmin.Header().Get("Location"))
	})

	t.Run("UUID Primary Key Views", func(t *testing.T) {
		ticketID := "11111111-2222-3333-4444-555555555555"
		err := db.Exec("INSERT INTO tickets (id, name, price, total_quantity, available_stock, status) VALUES (?, 'VIP Concert', 150.00, 100, 100, 'ACTIVE') ON CONFLICT (id) DO NOTHING", ticketID).Error
		assert.NoError(t, err)
		defer db.Exec("DELETE FROM tickets WHERE id = ?", ticketID)

		// List view renders valid UUID links
		reqList, _ := http.NewRequest(http.MethodGet, "/admin/info/tickets", nil)
		reqList.Header.Set("Cookie", cookie)
		wList := httptest.NewRecorder()
		r.ServeHTTP(wList, reqList)
		assert.Equal(t, http.StatusOK, wList.Code)
		assert.Contains(t, wList.Body.String(), ticketID)
		assert.Contains(t, wList.Body.String(), "__goadmin_detail_pk="+ticketID)
		assert.Contains(t, wList.Body.String(), "__goadmin_edit_pk="+ticketID)

		// Detail view parses UUID without errors
		reqDetail, _ := http.NewRequest(http.MethodGet, "/admin/info/tickets/detail?__goadmin_detail_pk="+ticketID, nil)
		reqDetail.Header.Set("Cookie", cookie)
		wDetail := httptest.NewRecorder()
		r.ServeHTTP(wDetail, reqDetail)
		assert.Equal(t, http.StatusOK, wDetail.Code)
		assert.NotContains(t, wDetail.Body.String(), "invalid input syntax for type uuid")
	})

	t.Run("Create Ticket with Auto-Timestamps and Nullable MaxBooking", func(t *testing.T) {
		reqNew, _ := http.NewRequest(http.MethodGet, "/admin/info/tickets/new", nil)
		reqNew.Header.Set("Cookie", cookie)
		wNew := httptest.NewRecorder()
		r.ServeHTTP(wNew, reqNew)
		assert.Equal(t, http.StatusOK, wNew.Code)
		token := extractCSRFToken(wNew.Body.String())
		assert.NotEmpty(t, token)

		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		_ = mw.WriteField("name", "Festival Ticket")
		_ = mw.WriteField("price", "120.00")
		_ = mw.WriteField("total_quantity", "50")
		_ = mw.WriteField("available_stock", "50")
		_ = mw.WriteField("status", "ACTIVE")
		_ = mw.WriteField("max_booking_per_user", "")
		_ = mw.WriteField("__go_admin_t_", token)
		_ = mw.WriteField("__go_admin_previous_", "/admin/info/tickets")
		_ = mw.Close()

		reqCreate, _ := http.NewRequest(http.MethodPost, "/admin/new/tickets", &buf)
		reqCreate.Header.Set("Cookie", cookie)
		reqCreate.Header.Set("Content-Type", mw.FormDataContentType())
		wCreate := httptest.NewRecorder()
		r.ServeHTTP(wCreate, reqCreate)
		assert.Equal(t, http.StatusOK, wCreate.Code)

		var createdTicket struct {
			ID                string
			MaxBookingPerUser *int
			CreatedAt         *time.Time
			UpdatedAt         *time.Time
		}
		err := db.Table("tickets").Where("name = ?", "Festival Ticket").Take(&createdTicket).Error
		assert.NoError(t, err)
		defer db.Exec("DELETE FROM tickets WHERE name = 'Festival Ticket'")

		assert.NotEmpty(t, createdTicket.ID)
		assert.Nil(t, createdTicket.MaxBookingPerUser, "max_booking_per_user should be NULL when left empty")
		assert.NotNil(t, createdTicket.CreatedAt)
		assert.NotNil(t, createdTicket.UpdatedAt)
		assert.NotZero(t, createdTicket.CreatedAt.Nanosecond(), "created_at should match PostgreSQL microsecond precision")

		// List view renders 'Unlimited' when max_booking_per_user is NULL
		reqList, _ := http.NewRequest(http.MethodGet, "/admin/info/tickets", nil)
		reqList.Header.Set("Cookie", cookie)
		wList := httptest.NewRecorder()
		r.ServeHTTP(wList, reqList)
		assert.Equal(t, http.StatusOK, wList.Code)
		assert.Contains(t, wList.Body.String(), "Unlimited")
	})

	t.Run("Update Ticket with Value and Revert to NULL", func(t *testing.T) {
		ticketID := "22222222-3333-4444-5555-666666666666"
		err := db.Exec("INSERT INTO tickets (id, name, price, total_quantity, available_stock, status) VALUES (?, 'Update Ticket', 100.00, 50, 50, 'ACTIVE') ON CONFLICT (id) DO NOTHING", ticketID).Error
		assert.NoError(t, err)
		defer db.Exec("DELETE FROM tickets WHERE id = ?", ticketID)

		// 1. Update with max_booking_per_user = 3
		reqEdit1, _ := http.NewRequest(http.MethodGet, "/admin/info/tickets/edit?__goadmin_edit_pk="+ticketID, nil)
		reqEdit1.Header.Set("Cookie", cookie)
		wEdit1 := httptest.NewRecorder()
		r.ServeHTTP(wEdit1, reqEdit1)
		token1 := extractCSRFToken(wEdit1.Body.String())
		assert.NotEmpty(t, token1)

		var buf1 bytes.Buffer
		mw1 := multipart.NewWriter(&buf1)
		_ = mw1.WriteField("id", ticketID)
		_ = mw1.WriteField("name", "Update Ticket")
		_ = mw1.WriteField("price", "100.00")
		_ = mw1.WriteField("total_quantity", "50")
		_ = mw1.WriteField("available_stock", "50")
		_ = mw1.WriteField("status", "ACTIVE")
		_ = mw1.WriteField("max_booking_per_user", "3")
		_ = mw1.WriteField("__go_admin_t_", token1)
		_ = mw1.WriteField("__goadmin_edit_pk", ticketID)
		_ = mw1.WriteField("__go_admin_previous_", "/admin/info/tickets")
		_ = mw1.Close()

		reqPost1, _ := http.NewRequest(http.MethodPost, "/admin/edit/tickets", &buf1)
		reqPost1.Header.Set("Cookie", cookie)
		reqPost1.Header.Set("Content-Type", mw1.FormDataContentType())
		wPost1 := httptest.NewRecorder()
		r.ServeHTTP(wPost1, reqPost1)
		assert.Equal(t, http.StatusOK, wPost1.Code)

		var ticketWithVal struct {
			MaxBookingPerUser *int
		}
		err = db.Table("tickets").Where("id = ?", ticketID).Take(&ticketWithVal).Error
		assert.NoError(t, err)
		assert.NotNil(t, ticketWithVal.MaxBookingPerUser)
		assert.Equal(t, 3, *ticketWithVal.MaxBookingPerUser)

		// 2. Clear max_booking_per_user back to empty / NULL
		reqEdit2, _ := http.NewRequest(http.MethodGet, "/admin/info/tickets/edit?__goadmin_edit_pk="+ticketID, nil)
		reqEdit2.Header.Set("Cookie", cookie)
		wEdit2 := httptest.NewRecorder()
		r.ServeHTTP(wEdit2, reqEdit2)
		token2 := extractCSRFToken(wEdit2.Body.String())
		assert.NotEmpty(t, token2)

		var buf2 bytes.Buffer
		mw2 := multipart.NewWriter(&buf2)
		_ = mw2.WriteField("id", ticketID)
		_ = mw2.WriteField("name", "Update Ticket")
		_ = mw2.WriteField("price", "100.00")
		_ = mw2.WriteField("total_quantity", "50")
		_ = mw2.WriteField("available_stock", "50")
		_ = mw2.WriteField("status", "ACTIVE")
		_ = mw2.WriteField("max_booking_per_user", "")
		_ = mw2.WriteField("__go_admin_t_", token2)
		_ = mw2.WriteField("__goadmin_edit_pk", ticketID)
		_ = mw2.WriteField("__go_admin_previous_", "/admin/info/tickets")
		_ = mw2.Close()

		reqPost2, _ := http.NewRequest(http.MethodPost, "/admin/edit/tickets", &buf2)
		reqPost2.Header.Set("Cookie", cookie)
		reqPost2.Header.Set("Content-Type", mw2.FormDataContentType())
		wPost2 := httptest.NewRecorder()
		r.ServeHTTP(wPost2, reqPost2)
		assert.Equal(t, http.StatusOK, wPost2.Code)

		var ticketCleared struct {
			MaxBookingPerUser *int
		}
		err = db.Table("tickets").Where("id = ?", ticketID).Take(&ticketCleared).Error
		assert.NoError(t, err)
		assert.Nil(t, ticketCleared.MaxBookingPerUser, "max_booking_per_user should revert to NULL when cleared")
	})
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
