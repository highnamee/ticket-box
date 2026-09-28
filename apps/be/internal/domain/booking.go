package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCanceled  BookingStatus = "CANCELED"
	BookingStatusPending   BookingStatus = "PENDING"
)

type Booking struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	TicketID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"ticket_id"`
	Quantity    int            `gorm:"not null" json:"quantity"`
	UnitPrice   float64        `gorm:"type:decimal(12,2);not null" json:"unit_price"`
	TotalAmount float64        `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status      BookingStatus  `gorm:"type:varchar(50);not null;default:'CONFIRMED'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Ticket *Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
}

func (Booking) TableName() string {
	return "bookings"
}

// BeforeCreate hook to generate UUID v7 (Time-Ordered) if not provided and ensure default status
func (b *Booking) BeforeCreate(_ *gorm.DB) error {
	if b.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		b.ID = id
	}
	if b.Status == "" {
		b.Status = BookingStatusConfirmed
	}
	return nil
}

// BookingRepository defines contract for database operations with concurrency protection
type BookingRepository interface {
	CreateBookingWithLock(ctx context.Context, userID, ticketID uuid.UUID, quantity int) (*Booking, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Booking, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]Booking, int64, error)
}

// BookingService defines contract for booking business use cases
type BookingService interface {
	BookTicket(ctx context.Context, userID, ticketID uuid.UUID, quantity int) (*Booking, error)
	GetBookingByID(ctx context.Context, userID, bookingID uuid.UUID) (*Booking, error)
	GetUserBookings(ctx context.Context, userID uuid.UUID, page, limit int) ([]Booking, int64, error)
}
