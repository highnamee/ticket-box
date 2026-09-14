package http

import (
	"ticket-box-be/internal/domain"

	"github.com/google/uuid"
)

// PaginationQuery defines the query parameters for pagination
type PaginationQuery struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}

// PublicTicketResponse represents the serialized public view of a ticket
type PublicTicketResponse struct {
	ID             uuid.UUID           `json:"id"`
	Name           string              `json:"name"`
	Description    string              `json:"description,omitempty"`
	Price          float64             `json:"price"`
	AvailableStock int                 `json:"available_stock"`
	Status         domain.TicketStatus `json:"status"`
}

func toPublicTicketResponse(t *domain.Ticket) PublicTicketResponse {
	return PublicTicketResponse{
		ID:             t.ID,
		Name:           t.Name,
		Description:    t.Description,
		Price:          t.Price,
		AvailableStock: t.AvailableStock,
		Status:         t.Status,
	}
}

func toPublicTicketResponseList(tickets []domain.Ticket) []PublicTicketResponse {
	res := make([]PublicTicketResponse, len(tickets))
	for i := range tickets {
		res[i] = toPublicTicketResponse(&tickets[i])
	}
	return res
}
