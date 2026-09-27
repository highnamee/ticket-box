package postgres

import (
	"context"
	"errors"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// CreateBookingWithLock executes ticket booking inside an ACID transaction with pessimistic row-level locking (SELECT ... FOR UPDATE).
// This guarantees race-condition prevention, stock accuracy under massive concurrency, and avoids overselling.
func (r *BookingRepository) CreateBookingWithLock(ctx context.Context, userID, ticketID uuid.UUID, quantity int) (*domain.Booking, error) {
	if quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}

	var createdBooking *domain.Booking

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ticket domain.Ticket

		// 1. Pessimistic Locking: SELECT ... FOR UPDATE
		// Serializes concurrent booking requests for the exact same ticket row at database level.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&ticket, "id = ?", ticketID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrTicketNotFound
			}
			return err
		}

		// 2. Validate ticket status (INACTIVE is a hidden system status, treated as not found)
		if ticket.Status != domain.TicketStatusActive {
			if ticket.Status == domain.TicketStatusSoldOut {
				return domain.ErrTicketSoldOut
			}
			return domain.ErrTicketNotFound
		}

		// 3. Check available stock
		if ticket.AvailableStock < quantity {
			return domain.ErrInsufficientStock
		}

		// 4. Atomically decrement stock and transition status if exhausted
		ticket.AvailableStock -= quantity
		if ticket.AvailableStock == 0 {
			ticket.Status = domain.TicketStatusSoldOut
		}

		if err := tx.Save(&ticket).Error; err != nil {
			return err
		}

		// 5. Create booking record
		totalAmount := ticket.Price * float64(quantity)
		booking := &domain.Booking{
			UserID:      userID,
			TicketID:    ticketID,
			Quantity:    quantity,
			UnitPrice:   ticket.Price,
			TotalAmount: totalAmount,
			Status:      domain.BookingStatusConfirmed,
		}

		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		booking.Ticket = &ticket
		createdBooking = booking
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdBooking, nil
}

func (r *BookingRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	var booking domain.Booking
	err := r.db.WithContext(ctx).
		Preload("Ticket").
		Preload("User").
		First(&booking, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) FindByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.Booking, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	var bookings []domain.Booking
	var total int64

	baseQuery := r.db.WithContext(ctx).
		Model(&domain.Booking{}).
		Where("user_id = ?", userID)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := baseQuery.
		Preload("Ticket").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error

	return bookings, total, err
}
