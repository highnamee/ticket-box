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

type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) CreateBookingWithLock(ctx context.Context, userID, ticketID uuid.UUID, quantity int) (*domain.Booking, error) {
	args := m.Called(ctx, userID, ticketID, quantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	booking, _ := args.Get(0).(*domain.Booking)
	return booking, args.Error(1)
}

func (m *MockBookingRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	booking, _ := args.Get(0).(*domain.Booking)
	return booking, args.Error(1)
}

func (m *MockBookingRepository) FindByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.Booking, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	total, _ := args.Get(1).(int64)
	if args.Get(0) == nil {
		return nil, total, args.Error(2)
	}
	bookings, _ := args.Get(0).([]domain.Booking)
	return bookings, total, args.Error(2)
}

func TestBookingService_BookTicket(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	svc := service.NewBookingService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	ticketID := uuid.New()

	t.Run("successful booking", func(t *testing.T) {
		expectedBooking := &domain.Booking{
			ID:          uuid.New(),
			UserID:      userID,
			TicketID:    ticketID,
			Quantity:    2,
			UnitPrice:   100.0,
			TotalAmount: 200.0,
			Status:      domain.BookingStatusConfirmed,
			CreatedAt:   time.Now(),
		}

		mockRepo.On("CreateBookingWithLock", ctx, userID, ticketID, 2).
			Return(expectedBooking, nil).Once()

		booking, err := svc.BookTicket(ctx, userID, ticketID, 2)
		require.NoError(t, err)
		assert.Equal(t, expectedBooking.ID, booking.ID)
		assert.Equal(t, 2, booking.Quantity)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid quantity", func(t *testing.T) {
		booking, err := svc.BookTicket(ctx, userID, ticketID, 0)
		assert.ErrorIs(t, err, domain.ErrInvalidQuantity)
		assert.Nil(t, booking)

		bookingNeg, err := svc.BookTicket(ctx, userID, ticketID, -1)
		assert.ErrorIs(t, err, domain.ErrInvalidQuantity)
		assert.Nil(t, bookingNeg)
	})

	t.Run("nil user ID returns unauthorized", func(t *testing.T) {
		booking, err := svc.BookTicket(ctx, uuid.Nil, ticketID, 1)
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
		assert.Nil(t, booking)
	})

	t.Run("nil ticket ID returns not found", func(t *testing.T) {
		booking, err := svc.BookTicket(ctx, userID, uuid.Nil, 1)
		assert.ErrorIs(t, err, domain.ErrTicketNotFound)
		assert.Nil(t, booking)
	})

	t.Run("propagates insufficient stock error", func(t *testing.T) {
		mockRepo.On("CreateBookingWithLock", ctx, userID, ticketID, 50).
			Return(nil, domain.ErrInsufficientStock).Once()

		booking, err := svc.BookTicket(ctx, userID, ticketID, 50)
		assert.ErrorIs(t, err, domain.ErrInsufficientStock)
		assert.Nil(t, booking)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingService_GetBookingByID(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	svc := service.NewBookingService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	bookingID := uuid.New()

	t.Run("success when user owns booking", func(t *testing.T) {
		expectedBooking := &domain.Booking{
			ID:       bookingID,
			UserID:   userID,
			Quantity: 1,
		}

		mockRepo.On("FindByID", ctx, bookingID).Return(expectedBooking, nil).Once()

		booking, err := svc.GetBookingByID(ctx, userID, bookingID)
		require.NoError(t, err)
		assert.Equal(t, bookingID, booking.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("forbidden when user does not own booking", func(t *testing.T) {
		otherUser := uuid.New()
		bookingOfOtherUser := &domain.Booking{
			ID:       bookingID,
			UserID:   otherUser,
			Quantity: 1,
		}

		mockRepo.On("FindByID", ctx, bookingID).Return(bookingOfOtherUser, nil).Once()

		booking, err := svc.GetBookingByID(ctx, userID, bookingID)
		assert.ErrorIs(t, err, domain.ErrForbidden)
		assert.Nil(t, booking)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, bookingID).Return(nil, domain.ErrBookingNotFound).Once()

		booking, err := svc.GetBookingByID(ctx, userID, bookingID)
		assert.ErrorIs(t, err, domain.ErrBookingNotFound)
		assert.Nil(t, booking)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingService_GetUserBookings(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	svc := service.NewBookingService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		list := []domain.Booking{
			{ID: uuid.New(), UserID: userID, Quantity: 1},
		}

		mockRepo.On("FindByUserID", ctx, userID, 1, 10).Return(list, int64(1), nil).Once()

		bookings, total, err := svc.GetUserBookings(ctx, userID, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, bookings, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("nil user returns unauthorized", func(t *testing.T) {
		bookings, total, err := svc.GetUserBookings(ctx, uuid.Nil, 1, 10)
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, bookings)
	})

	t.Run("error propagation", func(t *testing.T) {
		mockRepo.On("FindByUserID", ctx, userID, 1, 10).Return(nil, int64(0), errors.New("db error")).Once()

		bookings, total, err := svc.GetUserBookings(ctx, userID, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, bookings)
		mockRepo.AssertExpectations(t)
	})
}
