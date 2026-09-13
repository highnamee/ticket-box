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

func TestTicket_BeforeCreate_GeneratesUUIDv7(t *testing.T) {
	ticket := Ticket{
		Name: "General Admission",
	}

	err := ticket.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, ticket.ID)
	assert.Equal(t, uuid.Version(7), ticket.ID.Version())

	// Preserves existing ID if already set
	customID, _ := uuid.NewV7()
	ticketWithID := Ticket{
		ID:   customID,
		Name: "VIP",
	}
	err = ticketWithID.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.Equal(t, customID, ticketWithID.ID)
}
