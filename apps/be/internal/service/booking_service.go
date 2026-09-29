package service

import (
	"context"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
)

type BookingService struct {
	bookingRepo domain.BookingRepository
}

func NewBookingService(bookingRepo domain.BookingRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo}
}

// BookTicket orchestrates ticket booking with concurrency protection and business rules
func (s *BookingService) BookTicket(ctx context.Context, userID, ticketID uuid.UUID, quantity int) (*domain.Booking, error) {
	if quantity <= 0 {
		return nil, domain.ErrInvalidQuantity
	}
	if userID == uuid.Nil {
		return nil, domain.ErrUnauthorized
	}
	if ticketID == uuid.Nil {
		return nil, domain.ErrTicketNotFound
	}

	return s.bookingRepo.CreateBookingWithLock(ctx, userID, ticketID, quantity)
}

func (s *BookingService) GetBookingByID(ctx context.Context, userID, bookingID uuid.UUID) (*domain.Booking, error) {
	if bookingID == uuid.Nil {
		return nil, domain.ErrBookingNotFound
	}

	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	// Security: ensure users can only access their own bookings
	if booking.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return booking, nil
}

func (s *BookingService) GetUserBookings(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.Booking, int64, error) {
	if userID == uuid.Nil {
		return nil, 0, domain.ErrUnauthorized
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	return s.bookingRepo.FindByUserID(ctx, userID, page, limit)
}
