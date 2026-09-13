package postgres_test

import (
	"context"
	"testing"
	"time"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/pkg/testutil"
	"ticket-box-be/internal/repository/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTicketRepository_Create_And_FindByID(t *testing.T) {
	tx := testutil.SetupTestDB(t)
	repo := postgres.NewTicketRepository(tx)
	ctx := context.Background()

	ticket := &domain.Ticket{
		Name:           "Early Bird",
		Description:    "Early bird discount ticket",
		Price:          49.99,
		TotalQuantity:  50,
		AvailableStock: 50,
		Status:         domain.TicketStatusActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := repo.Create(ctx, ticket)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, ticket.ID)

	// Query from PostgreSQL
	found, err := repo.FindByID(ctx, ticket.ID)
	require.NoError(t, err)
	assert.Equal(t, ticket.ID, found.ID)
	assert.Equal(t, "Early Bird", found.Name)
	assert.Equal(t, 49.99, found.Price)
	assert.Equal(t, 50, found.AvailableStock)
}

func TestTicketRepository_FindAll(t *testing.T) {
	tx := testutil.SetupTestDB(t)
	repo := postgres.NewTicketRepository(tx)
	ctx := context.Background()

	t1 := &domain.Ticket{
		Name:           "VIP Ticket",
		Price:          199.00,
		TotalQuantity:  20,
		AvailableStock: 20,
		Status:         domain.TicketStatusActive,
	}
	t2 := &domain.Ticket{
		Name:           "General Admission",
		Price:          79.00,
		TotalQuantity:  100,
		AvailableStock: 100,
		Status:         domain.TicketStatusActive,
	}

	require.NoError(t, repo.Create(ctx, t1))
	require.NoError(t, repo.Create(ctx, t2))

	tickets, err := repo.FindAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(tickets), 2)
}

func TestTicketRepository_Update(t *testing.T) {
	tx := testutil.SetupTestDB(t)
	repo := postgres.NewTicketRepository(tx)
	ctx := context.Background()

	ticket := &domain.Ticket{
		Name:           "Old Ticket Name",
		Price:          100.00,
		TotalQuantity:  10,
		AvailableStock: 10,
		Status:         domain.TicketStatusActive,
	}
	require.NoError(t, repo.Create(ctx, ticket))

	// Update fields
	ticket.Name = "Updated Ticket Name"
	ticket.Price = 120.00
	ticket.Status = domain.TicketStatusSoldOut
	err := repo.Update(ctx, ticket)
	require.NoError(t, err)

	updated, err := repo.FindByID(ctx, ticket.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Ticket Name", updated.Name)
	assert.Equal(t, 120.00, updated.Price)
	assert.Equal(t, domain.TicketStatusSoldOut, updated.Status)
}

func TestTicketRepository_Delete(t *testing.T) {
	tx := testutil.SetupTestDB(t)
	repo := postgres.NewTicketRepository(tx)
	ctx := context.Background()

	ticket := &domain.Ticket{
		Name:           "Ticket To Delete",
		Price:          50.00,
		TotalQuantity:  5,
		AvailableStock: 5,
		Status:         domain.TicketStatusActive,
	}
	require.NoError(t, repo.Create(ctx, ticket))

	// Soft delete
	err := repo.Delete(ctx, ticket.ID)
	require.NoError(t, err)

	// Should not be found via normal find
	found, err := repo.FindByID(ctx, ticket.ID)
	assert.ErrorIs(t, err, domain.ErrTicketNotFound)
	assert.Nil(t, found)
}
