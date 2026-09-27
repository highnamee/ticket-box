package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/middleware"
	"ticket-box-be/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockBookingService struct {
	mock.Mock
}

func (m *MockBookingService) BookTicket(ctx context.Context, userID, ticketID uuid.UUID, quantity int) (*domain.Booking, error) {
	args := m.Called(ctx, userID, ticketID, quantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingService) GetBookingByID(ctx context.Context, userID, bookingID uuid.UUID) (*domain.Booking, error) {
	args := m.Called(ctx, userID, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingService) GetUserBookings(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.Booking, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Booking), args.Get(1).(int64), args.Error(2)
}

type SingleBookingResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    BookingResponse `json:"data"`
	Error   interface{}     `json:"error"`
}

func setupBookingTestRouter(handler *BookingHandler, userID *uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// If userID is provided, simulate AuthMiddleware setting context
	if userID != nil {
		r.Use(func(c *gin.Context) {
			c.Set(middleware.ContextKeyUserID, *userID)
			c.Next()
		})
	}

	r.POST("/api/v1/tickets/:id/book", handler.BookTicket)
	r.GET("/api/v1/users/me/bookings", handler.GetMyBookings)
	return r
}

func TestBookingHandler_BookTicket(t *testing.T) {
	userID := uuid.New()
	ticketID := uuid.New()

	t.Run("201 Created on successful booking", func(t *testing.T) {
		mockService := new(MockBookingService)
		handler := NewBookingHandler(mockService)
		router := setupBookingTestRouter(handler, &userID)

		expectedBooking := &domain.Booking{
			ID:          uuid.New(),
			UserID:      userID,
			TicketID:    ticketID,
			Quantity:    2,
			UnitPrice:   50.0,
			TotalAmount: 100.0,
			Status:      domain.BookingStatusConfirmed,
			CreatedAt:   time.Now(),
		}

		mockService.On("BookTicket", mock.Anything, userID, ticketID, 2).
			Return(expectedBooking, nil).Once()

		reqBody := `{"quantity": 2}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets/"+ticketID.String()+"/book", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp SingleBookingResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, 2, resp.Data.Quantity)
		assert.Equal(t, 100.0, resp.Data.TotalAmount)
		mockService.AssertExpectations(t)
	})

	t.Run("400 Bad Request on invalid ticket ID", func(t *testing.T) {
		mockService := new(MockBookingService)
		handler := NewBookingHandler(mockService)
		router := setupBookingTestRouter(handler, &userID)

		reqBody := `{"quantity": 1}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets/not-a-uuid/book", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("400 Bad Request on invalid quantity", func(t *testing.T) {
		mockService := new(MockBookingService)
		handler := NewBookingHandler(mockService)
		router := setupBookingTestRouter(handler, &userID)

		reqBody := `{"quantity": 0}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets/"+ticketID.String()+"/book", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("401 Unauthorized when user context is missing", func(t *testing.T) {
		mockService := new(MockBookingService)
		handler := NewBookingHandler(mockService)
		router := setupBookingTestRouter(handler, nil) // no auth

		reqBody := `{"quantity": 1}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets/"+ticketID.String()+"/book", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("404 Not Found when ticket does not exist", func(t *testing.T) {
		mockService := new(MockBookingService)
		handler := NewBookingHandler(mockService)
		router := setupBookingTestRouter(handler, &userID)

		mockService.On("BookTicket", mock.Anything, userID, ticketID, 1).
			Return(nil, domain.ErrTicketNotFound).Once()

		reqBody := `{"quantity": 1}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets/"+ticketID.String()+"/book", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("409 Conflict when ticket is sold out", func(t *testing.T) {
		mockService := new(MockBookingService)
		handler := NewBookingHandler(mockService)
		router := setupBookingTestRouter(handler, &userID)

		mockService.On("BookTicket", mock.Anything, userID, ticketID, 1).
			Return(nil, domain.ErrTicketSoldOut).Once()

		reqBody := `{"quantity": 1}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets/"+ticketID.String()+"/book", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("409 Conflict when stock is insufficient", func(t *testing.T) {
		mockService := new(MockBookingService)
		handler := NewBookingHandler(mockService)
		router := setupBookingTestRouter(handler, &userID)

		mockService.On("BookTicket", mock.Anything, userID, ticketID, 5).
			Return(nil, domain.ErrInsufficientStock).Once()

		reqBody := `{"quantity": 5}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/tickets/"+ticketID.String()+"/book", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestBookingHandler_GetMyBookings(t *testing.T) {
	userID := uuid.New()

	t.Run("200 OK returns user bookings with pagination", func(t *testing.T) {
		mockService := new(MockBookingService)
		handler := NewBookingHandler(mockService)
		router := setupBookingTestRouter(handler, &userID)

		bookings := []domain.Booking{
			{
				ID:          uuid.New(),
				UserID:      userID,
				TicketID:    uuid.New(),
				Quantity:    1,
				UnitPrice:   100.0,
				TotalAmount: 100.0,
				Status:      domain.BookingStatusConfirmed,
				CreatedAt:   time.Now(),
			},
		}

		mockService.On("GetUserBookings", mock.Anything, userID, 1, 10).
			Return(bookings, int64(1), nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/me/bookings?page=1&limit=10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		mockService.AssertExpectations(t)
	})

	t.Run("401 Unauthorized when not logged in", func(t *testing.T) {
		mockService := new(MockBookingService)
		handler := NewBookingHandler(mockService)
		router := setupBookingTestRouter(handler, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/me/bookings", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
