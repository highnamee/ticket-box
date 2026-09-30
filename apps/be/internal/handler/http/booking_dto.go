package http

import (
	"time"

	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
)

// BookTicketRequest payload for reserving and purchasing tickets
type BookTicketRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1,max=100"`
}

// BookingResponse represents the serialized view of a confirmed booking
type BookingResponse struct {
	ID          uuid.UUID            `json:"id"`
	UserID      uuid.UUID            `json:"user_id"`
	TicketID    uuid.UUID            `json:"ticket_id"`
	Quantity    int                  `json:"quantity"`
	UnitPrice   float64              `json:"unit_price"`
	TotalAmount float64              `json:"total_amount"`
	Status      domain.BookingStatus `json:"status"`
	CreatedAt   time.Time            `json:"created_at"`
	TicketName  string               `json:"ticket_name"`
}

func toBookingResponse(b *domain.Booking) BookingResponse {
	res := BookingResponse{
		ID:          b.ID,
		UserID:      b.UserID,
		TicketID:    b.TicketID,
		Quantity:    b.Quantity,
		UnitPrice:   b.UnitPrice,
		TotalAmount: b.TotalAmount,
		Status:      b.Status,
		CreatedAt:   b.CreatedAt,
	}
	if b.Ticket != nil {
		res.TicketName = b.Ticket.Name
	}
	return res
}

func toBookingResponseList(bookings []domain.Booking) []BookingResponse {
	res := make([]BookingResponse, len(bookings))
	for i := range bookings {
		res[i] = toBookingResponse(&bookings[i])
	}
	return res
}

// UserTicketQuotaResponse represents user's booking allowance and current usage for a ticket
type UserTicketQuotaResponse struct {
	TicketID          uuid.UUID `json:"ticket_id"`
	MaxBookingPerUser *int      `json:"max_booking_per_user"`
	CurrentBooked     int       `json:"current_booked"`
	RemainingQuota    *int      `json:"remaining_quota"`
}

func toUserTicketQuotaResponse(q *domain.UserTicketQuota) UserTicketQuotaResponse {
	return UserTicketQuotaResponse{
		TicketID:          q.TicketID,
		MaxBookingPerUser: q.MaxBookingPerUser,
		CurrentBooked:     q.CurrentBooked,
		RemainingQuota:    q.RemainingQuota,
	}
}
