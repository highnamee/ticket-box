package postgres_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/pkg/database"
	"ticket-box-be/internal/pkg/testutil"
	"ticket-box-be/internal/repository/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookingRepository_CreateBookingWithLock(t *testing.T) {
	tx := testutil.SetupTestDB(t)
	ticketRepo := postgres.NewTicketRepository(tx)
	userRepo := postgres.NewUserRepository(tx)
	bookingRepo := postgres.NewBookingRepository(tx)
	ctx := context.Background()

	// 1. Setup User
	user := &domain.User{
		Email:        "buyer@example.com",
		PasswordHash: "hashedsecret",
		FullName:     "Ticket Buyer",
		Role:         domain.RoleUser,
		Status:       domain.StatusActive,
	}
	require.NoError(t, userRepo.Create(ctx, user))

	// 2. Setup Ticket
	ticket := &domain.Ticket{
		Name:           "Rock Concert",
		Price:          100.0,
		TotalQuantity:  5,
		AvailableStock: 5,
		Status:         domain.TicketStatusActive,
	}
	require.NoError(t, ticketRepo.Create(ctx, ticket))

	t.Run("successful booking decrements stock correctly", func(t *testing.T) {
		booking, err := bookingRepo.CreateBookingWithLock(ctx, user.ID, ticket.ID, 2)
		require.NoError(t, err)
		assert.NotNil(t, booking)
		assert.Equal(t, user.ID, booking.UserID)
		assert.Equal(t, ticket.ID, booking.TicketID)
		assert.Equal(t, 2, booking.Quantity)
		assert.Equal(t, 100.0, booking.UnitPrice)
		assert.Equal(t, 200.0, booking.TotalAmount)
		assert.Equal(t, domain.BookingStatusConfirmed, booking.Status)

		// Verify ticket stock in DB
		updatedTicket, err := ticketRepo.FindByID(ctx, ticket.ID)
		require.NoError(t, err)
		assert.Equal(t, 3, updatedTicket.AvailableStock)
		assert.Equal(t, domain.TicketStatusActive, updatedTicket.Status)
	})

	t.Run("fails when requested quantity exceeds available stock", func(t *testing.T) {
		booking, err := bookingRepo.CreateBookingWithLock(ctx, user.ID, ticket.ID, 10)
		assert.ErrorIs(t, err, domain.ErrInsufficientStock)
		assert.Nil(t, booking)

		// Verify stock unchanged
		updatedTicket, err := ticketRepo.FindByID(ctx, ticket.ID)
		require.NoError(t, err)
		assert.Equal(t, 3, updatedTicket.AvailableStock)
	})

	t.Run("booking remaining stock marks ticket as SOLD_OUT", func(t *testing.T) {
		booking, err := bookingRepo.CreateBookingWithLock(ctx, user.ID, ticket.ID, 3)
		require.NoError(t, err)
		assert.NotNil(t, booking)

		// Verify ticket stock is 0 and status is SOLD_OUT
		updatedTicket, err := ticketRepo.FindByID(ctx, ticket.ID)
		require.NoError(t, err)
		assert.Equal(t, 0, updatedTicket.AvailableStock)
		assert.Equal(t, domain.TicketStatusSoldOut, updatedTicket.Status)
	})

	t.Run("booking sold out ticket returns ErrTicketSoldOut", func(t *testing.T) {
		booking, err := bookingRepo.CreateBookingWithLock(ctx, user.ID, ticket.ID, 1)
		assert.ErrorIs(t, err, domain.ErrTicketSoldOut)
		assert.Nil(t, booking)
	})

	t.Run("booking non-existent ticket returns ErrTicketNotFound", func(t *testing.T) {
		booking, err := bookingRepo.CreateBookingWithLock(ctx, user.ID, uuid.New(), 1)
		assert.ErrorIs(t, err, domain.ErrTicketNotFound)
		assert.Nil(t, booking)
	})

	t.Run("booking inactive ticket returns ErrTicketNotFound", func(t *testing.T) {
		inactiveTicket := &domain.Ticket{
			Name:           "Inactive Draft Ticket",
			Price:          100.0,
			TotalQuantity:  5,
			AvailableStock: 5,
			Status:         domain.TicketStatusInactive,
		}
		require.NoError(t, ticketRepo.Create(ctx, inactiveTicket))

		booking, err := bookingRepo.CreateBookingWithLock(ctx, user.ID, inactiveTicket.ID, 1)
		assert.ErrorIs(t, err, domain.ErrTicketNotFound)
		assert.Nil(t, booking)
	})

	t.Run("booking invalid quantity returns ErrInvalidQuantity", func(t *testing.T) {
		booking, err := bookingRepo.CreateBookingWithLock(ctx, user.ID, ticket.ID, 0)
		assert.ErrorIs(t, err, domain.ErrInvalidQuantity)
		assert.Nil(t, booking)

		bookingNeg, err := bookingRepo.CreateBookingWithLock(ctx, user.ID, ticket.ID, -5)
		assert.ErrorIs(t, err, domain.ErrInvalidQuantity)
		assert.Nil(t, bookingNeg)
	})
}

func TestBookingRepository_FindByID_And_FindByUserID(t *testing.T) {
	tx := testutil.SetupTestDB(t)
	ticketRepo := postgres.NewTicketRepository(tx)
	userRepo := postgres.NewUserRepository(tx)
	bookingRepo := postgres.NewBookingRepository(tx)
	ctx := context.Background()

	user := &domain.User{
		Email:        "buyer_history@example.com",
		PasswordHash: "hashedsecret",
		FullName:     "History Buyer",
		Role:         domain.RoleUser,
		Status:       domain.StatusActive,
	}
	require.NoError(t, userRepo.Create(ctx, user))

	ticket := &domain.Ticket{
		Name:           "Jazz Night",
		Price:          75.0,
		TotalQuantity:  10,
		AvailableStock: 10,
		Status:         domain.TicketStatusActive,
	}
	require.NoError(t, ticketRepo.Create(ctx, ticket))

	booking, err := bookingRepo.CreateBookingWithLock(ctx, user.ID, ticket.ID, 2)
	require.NoError(t, err)

	t.Run("FindByID returns existing booking with associations", func(t *testing.T) {
		found, err := bookingRepo.FindByID(ctx, booking.ID)
		require.NoError(t, err)
		assert.Equal(t, booking.ID, found.ID)
		assert.Equal(t, 2, found.Quantity)
		assert.Equal(t, 150.0, found.TotalAmount)
		assert.NotNil(t, found.Ticket)
		assert.Equal(t, ticket.Name, found.Ticket.Name)
	})

	t.Run("FindByID returns ErrBookingNotFound for invalid ID", func(t *testing.T) {
		found, err := bookingRepo.FindByID(ctx, uuid.New())
		assert.ErrorIs(t, err, domain.ErrBookingNotFound)
		assert.Nil(t, found)
	})

	t.Run("FindByUserID returns user bookings with pagination", func(t *testing.T) {
		bookings, total, err := bookingRepo.FindByUserID(ctx, user.ID, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, bookings, 1)
		assert.Equal(t, booking.ID, bookings[0].ID)
	})
}

// TestBookingRepository_ConcurrentBooking simulates high concurrency to verify race condition handling
func TestBookingRepository_ConcurrentBooking(t *testing.T) {
	cfg := testutil.GetTestConfig(t)
	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Skipf("Skipping concurrent integration test (database unreachable: %v)", err)
		return
	}
	if err := database.MigrateUp(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	ticketRepo := postgres.NewTicketRepository(db)
	userRepo := postgres.NewUserRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)
	ctx := context.Background()

	// 1. Create a User for concurrent booking
	user := &domain.User{
		Email:        "concurrent_buyer_" + uuid.NewString()[:8] + "@example.com",
		PasswordHash: "hashedsecret",
		FullName:     "Concurrent Buyer",
		Role:         domain.RoleUser,
		Status:       domain.StatusActive,
	}
	require.NoError(t, userRepo.Create(ctx, user))
	t.Cleanup(func() {
		_ = db.Exec("DELETE FROM bookings WHERE user_id = ?", user.ID).Error
		_ = db.Exec("DELETE FROM users WHERE id = ?", user.ID).Error
	})

	// 2. Setup Ticket with exactly 10 available tickets
	initialStock := 10
	ticket := &domain.Ticket{
		Name:           "Flash Sale Ticket " + uuid.NewString()[:8],
		Price:          50.0,
		TotalQuantity:  initialStock,
		AvailableStock: initialStock,
		Status:         domain.TicketStatusActive,
	}
	require.NoError(t, ticketRepo.Create(ctx, ticket))
	t.Cleanup(func() {
		_ = db.Exec("DELETE FROM bookings WHERE ticket_id = ?", ticket.ID).Error
		_ = db.Exec("DELETE FROM tickets WHERE id = ?", ticket.ID).Error
	})

	// 3. Launch 100 concurrent goroutines, each attempting to book 1 ticket
	concurrentRequests := 100
	var wg sync.WaitGroup
	var successCount int32
	var failedCount int32

	startSignal := make(chan struct{})

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startSignal // synchronize start

			_, bookErr := bookingRepo.CreateBookingWithLock(ctx, user.ID, ticket.ID, 1)
			if bookErr == nil {
				atomic.AddInt32(&successCount, 1)
			} else {
				atomic.AddInt32(&failedCount, 1)
			}
		}()
	}

	// Release all goroutines simultaneously
	close(startSignal)
	wg.Wait()

	// 4. Assertions:
	// Exactly initialStock (10) bookings must succeed
	assert.Equal(t, int32(initialStock), successCount, "Exactly 10 bookings should succeed")
	// The rest (40) must fail
	assert.Equal(t, int32(concurrentRequests-initialStock), failedCount, "40 bookings should fail due to out-of-stock")

	// Final ticket stock in DB must be exactly 0, never negative
	finalTicket, err := ticketRepo.FindByID(ctx, ticket.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, finalTicket.AvailableStock, "Final stock must be exactly 0")
	assert.Equal(t, domain.TicketStatusSoldOut, finalTicket.Status, "Status must be SOLD_OUT")

	// Total quantity recorded in bookings table must equal initialStock
	var totalBookedQuantity int64
	err = db.Model(&domain.Booking{}).
		Where("ticket_id = ?", ticket.ID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalBookedQuantity).Error
	require.NoError(t, err)
	assert.Equal(t, int64(initialStock), totalBookedQuantity, "Sum of booked tickets must match initial stock")
}
