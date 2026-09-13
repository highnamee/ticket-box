package domain_test

import (
	"testing"
	"time"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTicket_TableName(t *testing.T) {
	ticket := domain.Ticket{}
	assert.Equal(t, "tickets", ticket.TableName())
}

func TestTicket_DefaultsAndStatus(t *testing.T) {
	now := time.Now()
	testID := uuid.New()
	ticket := domain.Ticket{
		ID:             testID,
		Name:           "VIP Pass",
		Description:    "Access to VIP lounge",
		Price:          150.00,
		TotalQuantity:  100,
		AvailableStock: 100,
		Status:         domain.TicketStatusActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	assert.Equal(t, testID, ticket.ID)
	assert.Equal(t, "VIP Pass", ticket.Name)
	assert.Equal(t, 150.00, ticket.Price)
	assert.Equal(t, 100, ticket.TotalQuantity)
	assert.Equal(t, 100, ticket.AvailableStock)
	assert.Equal(t, domain.TicketStatusActive, ticket.Status)
}

func TestTicket_BeforeCreate(t *testing.T) {
	t.Run("generates UUIDv7 when ID is nil", func(t *testing.T) {
		ticket := domain.Ticket{Name: "General Admission"}

		err := ticket.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, ticket.ID)
		assert.Equal(t, uuid.Version(7), ticket.ID.Version())
	})

	t.Run("preserves existing ID when provided", func(t *testing.T) {
		customID, err := uuid.NewV7()
		assert.NoError(t, err)

		ticket := domain.Ticket{
			ID:   customID,
			Name: "VIP Pass",
		}

		err = ticket.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.Equal(t, customID, ticket.ID)
	})

	t.Run("defaults status to INACTIVE when empty", func(t *testing.T) {
		ticket := domain.Ticket{Name: "Early Bird"}

		err := ticket.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.Equal(t, domain.TicketStatusInactive, ticket.Status)
	})

	t.Run("preserves existing status when provided", func(t *testing.T) {
		ticket := domain.Ticket{
			Name:   "VIP Pass",
			Status: domain.TicketStatusActive,
		}

		err := ticket.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.Equal(t, domain.TicketStatusActive, ticket.Status)
	})
}
