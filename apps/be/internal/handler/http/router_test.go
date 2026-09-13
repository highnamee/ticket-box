package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter_HealthEndpoints(t *testing.T) {
	healthHandler := NewHealthHandler()
	router := NewRouter(RouterConfig{
		HealthHandler: healthHandler,
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
			name:         "Not found route",
			url:          "/unknown-route",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
