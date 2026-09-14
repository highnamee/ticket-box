package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (m *MockTicketRepository) FindAll(ctx context.Context) ([]domain.Ticket, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Ticket), args.Error(1)
}

func (m *MockTicketRepository) FindPublic(ctx context.Context, page, limit int) ([]domain.Ticket, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Ticket), args.Get(1).(int64), args.Error(2)
}

func (m *MockTicketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTicketService_GetPublicTickets(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	ticketSvc := service.NewTicketService(mockRepo)

	ctx := context.Background()
	mockTickets := []domain.Ticket{
		{
			ID:             uuid.New(),
			Name:           "VIP Ticket",
			Price:          100.0,
			TotalQuantity:  50,
			AvailableStock: 50,
			Status:         domain.TicketStatusActive,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("FindPublic", ctx, 1, 10).Return(mockTickets, int64(1), nil).Once()

		tickets, total, err := ticketSvc.GetPublicTickets(ctx, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, tickets, 1)
		assert.Equal(t, "VIP Ticket", tickets[0].Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("FindPublic", ctx, 1, 10).Return(nil, int64(0), errors.New("db error")).Once()

		tickets, total, err := ticketSvc.GetPublicTickets(ctx, 1, 10)
		require.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, tickets)
		mockRepo.AssertExpectations(t)
	})
}

func TestTicketService_GetTicketByID(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	ticketSvc := service.NewTicketService(mockRepo)

	ctx := context.Background()
	testID := uuid.New()
	mockTicket := &domain.Ticket{
		ID:    testID,
		Name:  "Standard Ticket",
		Price: 50.0,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, testID).Return(mockTicket, nil).Once()

		ticket, err := ticketSvc.GetTicketByID(ctx, testID)
		require.NoError(t, err)
		assert.Equal(t, testID, ticket.ID)
		assert.Equal(t, "Standard Ticket", ticket.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, testID).Return(nil, domain.ErrTicketNotFound).Once()

		ticket, err := ticketSvc.GetTicketByID(ctx, testID)
		require.ErrorIs(t, err, domain.ErrTicketNotFound)
		assert.Nil(t, ticket)
		mockRepo.AssertExpectations(t)
	})
}
