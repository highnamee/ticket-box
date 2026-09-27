package postgres

import (
	"context"
	"errors"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// CreateBookingWithLock executes ticket booking using an atomic conditional UPDATE and ACID transaction.
// The atomic UPDATE decrements available_stock only when available_stock >= quantity and status == ACTIVE in a single statement.
// This holds locks for mere microseconds (drastically cutting lock contention under 10k requests) and fails immediately when out of stock.
func (r *BookingRepository) CreateBookingWithLock(ctx context.Context, userID, ticketID uuid.UUID, quantity int) (*domain.Booking, error) {
	if quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}

	var createdBooking *domain.Booking

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ticket domain.Ticket

		// 1. Atomic Single-Statement Decrement
		// PostgreSQL acquires row lock only during the execution of this single statement.
		// If stock is insufficient or status != ACTIVE, RowsAffected == 0 immediately without blocking other connections.
		updateQuery := `
			UPDATE tickets
			SET available_stock = available_stock - ?,
			    status = CASE WHEN available_stock - ? = 0 THEN ? ELSE status END,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = ? AND status = ? AND available_stock >= ?
			RETURNING id, name, price, available_stock, status, created_at, updated_at
		`
		result := tx.Raw(updateQuery, quantity, quantity, domain.TicketStatusSoldOut, ticketID, domain.TicketStatusActive, quantity).Scan(&ticket)
		if result.Error != nil {
			return result.Error
		}

		// 2. If update didn't affect any row, inspect without locks to diagnose the exact error for client
		if result.RowsAffected == 0 {
			var existing domain.Ticket
			if err := tx.First(&existing, "id = ?", ticketID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return domain.ErrTicketNotFound
				}
				return err
			}

			if existing.Status != domain.TicketStatusActive {
				if existing.Status == domain.TicketStatusSoldOut {
					return domain.ErrTicketSoldOut
				}
				// INACTIVE is a hidden system status -> ErrTicketNotFound
				return domain.ErrTicketNotFound
			}

			if existing.AvailableStock < quantity {
				return domain.ErrInsufficientStock
			}

			return domain.ErrInsufficientStock
		}

		// 3. Create confirmed booking record
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
