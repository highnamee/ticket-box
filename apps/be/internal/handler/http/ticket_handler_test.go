package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockTicketRepository for unit testing
type MockTicketRepository struct {
	mock.Mock
}

func (m *MockTicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	ticket, _ := args.Get(0).(*domain.Ticket)
	return ticket, args.Error(1)
}

func (m *MockTicketRepository) FindAll(ctx context.Context) ([]domain.Ticket, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	tickets, _ := args.Get(0).([]domain.Ticket)
	return tickets, args.Error(1)
}

func (m *MockTicketRepository) FindPublic(ctx context.Context, page, limit int) ([]domain.Ticket, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	tickets, _ := args.Get(0).([]domain.Ticket)
	return tickets, args.Get(1).(int64), args.Error(2)
}

func (m *MockTicketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type PaginatedTicketResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    response.PaginatedData `json:"data"`
	Error   interface{}            `json:"error"`
}

func TestTicketHandler_GetPublicTickets_Success(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	handler := NewTicketHandler(mockRepo)

	mockTickets := []domain.Ticket{
		{
			ID:             uuid.New(),
			Name:           "VIP Pass",
			Price:          150.0,
			TotalQuantity:  50,
			AvailableStock: 50,
			Status:         domain.TicketStatusActive,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			Name:           "Early Bird",
			Price:          50.0,
			TotalQuantity:  30,
			AvailableStock: 0,
			Status:         domain.TicketStatusSoldOut,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	// Default pagination page=1, limit=10
	mockRepo.On("FindPublic", mock.Anything, 1, 10).Return(mockTickets, int64(2), nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	c.Request = req

	handler.GetPublicTickets(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var res PaginatedTicketResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.True(t, res.Success)
	assert.Equal(t, "Tickets retrieved successfully", res.Message)
	require.NotNil(t, res.Data.Pagination)
	assert.Equal(t, 1, res.Data.Pagination.Page)
	assert.Equal(t, 10, res.Data.Pagination.Limit)
	assert.Equal(t, int64(2), res.Data.Pagination.TotalItems)
	assert.Equal(t, 1, res.Data.Pagination.TotalPages)
	mockRepo.AssertExpectations(t)
}

func TestTicketHandler_GetPublicTickets_WithCustomPagination(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	handler := NewTicketHandler(mockRepo)

	mockTickets := []domain.Ticket{
		{
			ID:             uuid.New(),
			Name:           "VIP Pass",
			Price:          150.0,
			TotalQuantity:  50,
			AvailableStock: 50,
			Status:         domain.TicketStatusActive,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	mockRepo.On("FindPublic", mock.Anything, 2, 5).Return(mockTickets, int64(11), nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets?page=2&limit=5", nil)
	c.Request = req

	handler.GetPublicTickets(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var res PaginatedTicketResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.True(t, res.Success)
	require.NotNil(t, res.Data.Pagination)
	assert.Equal(t, 2, res.Data.Pagination.Page)
	assert.Equal(t, 5, res.Data.Pagination.Limit)
	assert.Equal(t, int64(11), res.Data.Pagination.TotalItems)
	assert.Equal(t, 3, res.Data.Pagination.TotalPages)
	mockRepo.AssertExpectations(t)
}

func TestTicketHandler_GetPublicTickets_InvalidQueryParams(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	handler := NewTicketHandler(mockRepo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets?page=0", nil) // page < 1
	c.Request = req

	handler.GetPublicTickets(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.False(t, res.Success)
	assert.Equal(t, "Invalid pagination query", res.Message)
	assert.Contains(t, res.Error, "Field validation for 'Page' failed on the 'min' tag")
}

func TestTicketHandler_GetPublicTickets_InvalidLimit(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	handler := NewTicketHandler(mockRepo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets?limit=150", nil) // limit > 100
	c.Request = req

	handler.GetPublicTickets(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.False(t, res.Success)
	assert.Equal(t, "Invalid pagination query", res.Message)
	assert.Contains(t, res.Error, "Field validation for 'Limit' failed on the 'max' tag")
}

func TestTicketHandler_GetPublicTickets_NonNumericQuery(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	handler := NewTicketHandler(mockRepo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets?page=abc", nil)
	c.Request = req

	handler.GetPublicTickets(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.False(t, res.Success)
	assert.Equal(t, "Invalid pagination query", res.Message)
	assert.Contains(t, res.Error, "strconv.ParseInt: parsing \"abc\": invalid syntax")
}

func TestTicketHandler_GetPublicTickets_Error(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	handler := NewTicketHandler(mockRepo)

	mockRepo.On("FindPublic", mock.Anything, 1, 10).Return(nil, int64(0), errors.New("database failure"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	c.Request = req

	handler.GetPublicTickets(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.False(t, res.Success)
	assert.Equal(t, "Failed to retrieve tickets", res.Message)
	assert.Equal(t, "database failure", res.Error)
	mockRepo.AssertExpectations(t)
}
