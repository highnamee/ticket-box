package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketStatus string

const (
	TicketStatusActive   TicketStatus = "ACTIVE"
	TicketStatusSoldOut  TicketStatus = "SOLD_OUT"
	TicketStatusInactive TicketStatus = "INACTIVE"
)

type Ticket struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description,omitempty"`
	Price          float64        `gorm:"type:decimal(12,2);not null;default:0" json:"price"`
	TotalQuantity  int            `gorm:"not null;default:0" json:"total_quantity"`
	AvailableStock int            `gorm:"not null;default:0" json:"available_stock"`
	Status         TicketStatus   `gorm:"type:varchar(50);not null;default:'INACTIVE'" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Ticket) TableName() string {
	return "tickets"
}

// BeforeCreate hook to generate UUID v7 (Time-Ordered) if not provided and ensure default status
func (t *Ticket) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		t.ID = id
	}
	if t.Status == "" {
		t.Status = TicketStatusInactive
	}
	return nil
}

// TicketRepository defines contract for database operations
type TicketRepository interface {
	Create(ctx context.Context, ticket *Ticket) error
	FindByID(ctx context.Context, id uuid.UUID) (*Ticket, error)
	FindAll(ctx context.Context) ([]Ticket, error)
	FindPublic(ctx context.Context, page, limit int) ([]Ticket, int64, error)
	Update(ctx context.Context, ticket *Ticket) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// TicketService defines contract for ticket business logic use cases
type TicketService interface {
	GetPublicTickets(ctx context.Context, page, limit int) ([]Ticket, int64, error)
	GetTicketByID(ctx context.Context, id uuid.UUID) (*Ticket, error)
}

