package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ticket-box-be/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouter_Endpoints(t *testing.T) {
	healthHandler := NewHealthHandler()
	mockRepo := new(MockTicketRepository)
	mockRepo.On("FindPublic", mock.Anything, mock.Anything, mock.Anything).Return([]domain.Ticket{}, int64(0), nil)
	ticketHandler := NewTicketHandler(mockRepo)

	router := NewRouter(RouterConfig{
		AppEnv:        "development",
		HealthHandler: healthHandler,
		TicketHandler: ticketHandler,
	})

	tests := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{
			name:         "Root health endpoint",
			url:          "/health",
			expectedCode: http.StatusOK,
		},
		{
			name:         "API v1 ping endpoint",
			url:          "/api/v1/ping",
			expectedCode: http.StatusOK,
		},
		{
			name:         "API v1 tickets endpoint",
			url:          "/api/v1/tickets",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Swagger UI in development (redirects or serves)",
			url:          "/swagger/index.html",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Not found route",
			url:          "/unknown-route",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			req.RequestURI = tt.url
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestRouter_SwaggerDisabledInProduction(t *testing.T) {
	healthHandler := NewHealthHandler()
	router := NewRouter(RouterConfig{
		AppEnv:        "production",
		HealthHandler: healthHandler,
	})

	req, _ := http.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	req.RequestURI = "/swagger/index.html"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
