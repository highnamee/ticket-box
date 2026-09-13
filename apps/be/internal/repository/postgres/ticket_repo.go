package postgres

import (
	"context"
	"errors"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *TicketRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).First(&ticket, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTicketNotFound
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) FindAll(ctx context.Context) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	err := r.db.WithContext(ctx).
		Order("price ASC").
		Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	result := r.db.WithContext(ctx).Save(ticket)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrTicketNotFound
	}
	return nil
}

func (r *TicketRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.Ticket{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrTicketNotFound
	}
	return nil
}
