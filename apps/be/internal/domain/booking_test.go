package domain_test

import (
	"testing"
	"time"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBooking_TableName(t *testing.T) {
	booking := domain.Booking{}
	assert.Equal(t, "bookings", booking.TableName())
}

func TestBooking_DefaultsAndStatus(t *testing.T) {
	now := time.Now()
	testID := uuid.New()
	userID := uuid.New()
	ticketID := uuid.New()

	booking := domain.Booking{
		ID:          testID,
		UserID:      userID,
		TicketID:    ticketID,
		Quantity:    2,
		UnitPrice:   50.0,
		TotalAmount: 100.0,
		Status:      domain.BookingStatusConfirmed,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, testID, booking.ID)
	assert.Equal(t, userID, booking.UserID)
	assert.Equal(t, ticketID, booking.TicketID)
	assert.Equal(t, 2, booking.Quantity)
	assert.Equal(t, 50.0, booking.UnitPrice)
	assert.Equal(t, 100.0, booking.TotalAmount)
	assert.Equal(t, domain.BookingStatusConfirmed, booking.Status)
}

func TestBooking_BeforeCreate(t *testing.T) {
	t.Run("generates UUIDv7 when ID is nil", func(t *testing.T) {
		booking := domain.Booking{
			UserID:   uuid.New(),
			TicketID: uuid.New(),
			Quantity: 1,
		}

		err := booking.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, booking.ID)
		assert.Equal(t, uuid.Version(7), booking.ID.Version())
	})

	t.Run("defaults status to CONFIRMED when empty", func(t *testing.T) {
		booking := domain.Booking{
			UserID:   uuid.New(),
			TicketID: uuid.New(),
			Quantity: 1,
		}

		err := booking.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.Equal(t, domain.BookingStatusConfirmed, booking.Status)
	})

	t.Run("preserves existing ID and status when provided", func(t *testing.T) {
		customID, err := uuid.NewV7()
		assert.NoError(t, err)

		booking := domain.Booking{
			ID:     customID,
			Status: domain.BookingStatusCanceled,
		}

		err = booking.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.Equal(t, customID, booking.ID)
		assert.Equal(t, domain.BookingStatusCanceled, booking.Status)
	})
}
