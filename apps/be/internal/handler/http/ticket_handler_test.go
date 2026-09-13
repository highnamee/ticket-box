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

func (m *MockTicketRepository) FindPublic(ctx context.Context) ([]domain.Ticket, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	tickets, _ := args.Get(0).([]domain.Ticket)
	return tickets, args.Error(1)
}

func (m *MockTicketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
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

	mockRepo.On("FindPublic", mock.Anything).Return(mockTickets, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	c.Request = req

	handler.GetPublicTickets(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var res response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.True(t, res.Success)
	assert.Equal(t, "Tickets retrieved successfully", res.Message)
	mockRepo.AssertExpectations(t)
}

func TestTicketHandler_GetPublicTickets_Error(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	handler := NewTicketHandler(mockRepo)

	mockRepo.On("FindPublic", mock.Anything).Return(nil, errors.New("database failure"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	c.Request = req

	handler.GetPublicTickets(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.False(t, res.Success)
	assert.Equal(t, "Failed to retrieve tickets", res.Message)
	mockRepo.AssertExpectations(t)
}
