package http

import (
	"net/http"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketRepo domain.TicketRepository
}

func NewTicketHandler(ticketRepo domain.TicketRepository) *TicketHandler {
	return &TicketHandler{ticketRepo: ticketRepo}
}

// GetPublicTickets godoc
// @Summary      Get public tickets
// @Description  Retrieve all public tickets (ACTIVE and SOLD_OUT, excluding private INACTIVE tickets)
// @Tags         Tickets
// @Produce      json
// @Success      200  {object}  response.APIResponse{data=[]domain.Ticket}  "Tickets retrieved successfully"
// @Failure      500  {object}  response.APIResponse                        "Internal server error"
// @Router       /tickets [get]
func (h *TicketHandler) GetPublicTickets(c *gin.Context) {
	tickets, err := h.ticketRepo.FindPublic(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve tickets")
		return
	}

	response.Success(c, http.StatusOK, "Tickets retrieved successfully", tickets)
}
