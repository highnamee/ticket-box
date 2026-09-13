package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTicket_TableName(t *testing.T) {
	ticket := Ticket{}
	assert.Equal(t, "tickets", ticket.TableName())
}

func TestTicket_DefaultsAndStatus(t *testing.T) {
	now := time.Now()
	testID := uuid.New()
	ticket := Ticket{
		ID:             testID,
		Name:           "VIP Pass",
		Description:    "Access to VIP lounge",
		Price:          150.00,
		TotalQuantity:  100,
		AvailableStock: 100,
		Status:         TicketStatusActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	assert.Equal(t, testID, ticket.ID)
	assert.Equal(t, "VIP Pass", ticket.Name)
	assert.Equal(t, 150.00, ticket.Price)
	assert.Equal(t, 100, ticket.TotalQuantity)
	assert.Equal(t, 100, ticket.AvailableStock)
	assert.Equal(t, TicketStatusActive, ticket.Status)
}
