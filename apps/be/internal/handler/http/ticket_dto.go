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
	ID                uuid.UUID           `json:"id"`
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	Price             float64             `json:"price"`
	AvailableStock    int                 `json:"available_stock"`
	MaxBookingPerUser *int                `json:"max_booking_per_user,omitempty"`
	Status            domain.TicketStatus `json:"status"`
	IsSoldOut         bool                `json:"is_sold_out"`
	CategoryID        *uuid.UUID          `json:"category_id,omitempty"`
	Category          string              `json:"category"`
	Venue             string              `json:"venue"`
	Date              string              `json:"date"`
	Time              string              `json:"time"`
	ImageURL          string              `json:"image_url"`
	Tags              []string            `json:"tags"`
	Featured          bool                `json:"featured"`
}

func toPublicTicketResponse(t *domain.Ticket) PublicTicketResponse {
	tags := []string(t.Tags)
	if tags == nil {
		tags = []string{}
	}
	categoryName := ""
	if t.Category != nil {
		categoryName = t.Category.Name
	}
	return PublicTicketResponse{
		ID:                t.ID,
		Name:              t.Name,
		Description:       t.Description,
		Price:             t.Price,
		AvailableStock:    t.AvailableStock,
		MaxBookingPerUser: t.MaxBookingPerUser,
		Status:            t.Status,
		IsSoldOut:         t.IsSoldOut(),
		CategoryID:        t.CategoryID,
		Category:          categoryName,
		Venue:             t.Venue,
		Date:              t.Date,
		Time:              t.Time,
		ImageURL:          t.ImageURL,
		Tags:              tags,
		Featured:          t.Featured,
	}
}

func toPublicTicketResponseList(tickets []domain.Ticket) []PublicTicketResponse {
	res := make([]PublicTicketResponse, len(tickets))
	for i := range tickets {
		res[i] = toPublicTicketResponse(&tickets[i])
	}
	return res
}
