package service

import (
	"context"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
)

type TicketService struct {
	ticketRepo domain.TicketRepository
}

func NewTicketService(ticketRepo domain.TicketRepository) *TicketService {
	return &TicketService{ticketRepo: ticketRepo}
}

func (s *TicketService) GetPublicTickets(ctx context.Context, page, limit int) ([]domain.Ticket, int64, error) {
	return s.ticketRepo.FindPublic(ctx, page, limit)
}

func (s *TicketService) GetTicketByID(ctx context.Context, id uuid.UUID) (*domain.Ticket, error) {
	return s.ticketRepo.FindByID(ctx, id)
}
