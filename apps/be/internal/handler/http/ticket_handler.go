package http

import (
	"net/http"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketService domain.TicketService
}

func NewTicketHandler(ticketService domain.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

// GetPublicTickets godoc
// @Summary      Get public tickets
// @Description  Retrieve all public tickets (ACTIVE and SOLD_OUT, excluding private INACTIVE tickets) with pagination
// @Tags         Tickets
// @Produce      json
// @Param        page   query     int  false  "Page number (default 1, min 1)"                default(1)
// @Param        limit  query     int  false  "Items per page (default 10, min 1, max 100)"   default(10)
// @Success      200    {object}  response.APIResponse{data=response.PaginatedData}              "Tickets retrieved successfully"
// @Failure      400    {object}  response.APIResponse                                           "Bad request"
// @Failure      500    {object}  response.APIResponse                                           "Internal server error"
// @Router       /tickets [get]
func (h *TicketHandler) GetPublicTickets(c *gin.Context) {
	var query PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "Invalid pagination query", err.Error())
		return
	}

	tickets, total, err := h.ticketService.GetPublicTickets(c.Request.Context(), query.Page, query.Limit)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve tickets", err.Error())
		return
	}

	items := toPublicTicketResponseList(tickets)
	pagination := response.NewPaginationMeta(query.Page, query.Limit, total)
	response.SuccessWithPagination(c, http.StatusOK, "Tickets retrieved successfully", items, pagination)
}
