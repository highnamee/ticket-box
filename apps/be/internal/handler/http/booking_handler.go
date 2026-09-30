package http

import (
	"errors"
	"net/http"

	"ticket-box-be/internal/domain"
	"ticket-box-be/internal/middleware"
	"ticket-box-be/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler struct {
	bookingService domain.BookingService
}

func NewBookingHandler(bookingService domain.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

// BookTicket godoc
// @Summary      Book a ticket
// @Description  Reserve and purchase tickets with concurrency protection against race conditions and overselling
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string             true  "Ticket UUID"
// @Param        request body      BookTicketRequest  true  "Booking details"
// @Success      201     {object}  response.APIResponse{data=BookingResponse}  "Ticket booked successfully"
// @Failure      400     {object}  response.APIResponse                        "Bad request - invalid quantity, payload, or exceeded max booking limit per user"
// @Failure      401     {object}  response.APIResponse                        "Unauthorized"
// @Failure      404     {object}  response.APIResponse                        "Ticket not found"
// @Failure      409     {object}  response.APIResponse                        "Conflict - ticket sold out or insufficient stock"
// @Failure      500     {object}  response.APIResponse                        "Internal server error"
// @Router       /tickets/{id}/book [post]
func (h *BookingHandler) BookTicket(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid ticket ID format", err.Error())
		return
	}

	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", domain.ErrUnauthorized)
		return
	}

	var req BookTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid booking request payload", err.Error())
		return
	}

	booking, err := h.bookingService.BookTicket(c.Request.Context(), userID, ticketID, req.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTicketNotFound):
			response.NotFound(c, "Ticket not found", err.Error())
		case errors.Is(err, domain.ErrTicketSoldOut):
			response.Error(c, http.StatusConflict, "Ticket is sold out", err.Error())
		case errors.Is(err, domain.ErrInsufficientStock):
			response.Error(c, http.StatusConflict, "Insufficient ticket stock available", err.Error())
		case errors.Is(err, domain.ErrInvalidQuantity):
			response.BadRequest(c, "Invalid ticket quantity requested", err.Error())
		case errors.Is(err, domain.ErrMaxBookingLimitExceeded):
			response.BadRequest(c, "Exceeded maximum booking limit per user", err.Error())
		case errors.Is(err, domain.ErrUnauthorized):
			response.Error(c, http.StatusUnauthorized, "Unauthorized", err.Error())
		default:
			response.InternalServerError(c, "Failed to book ticket", err)
		}
		return
	}

	response.Success(c, http.StatusCreated, "Ticket booked successfully", toBookingResponse(booking))
}

// GetMyBookings godoc
// @Summary      Get current user's bookings
// @Description  Retrieve all bookings belonging to the currently authenticated user
// @Tags         Bookings
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Page number (default 1, min 1)"              default(1)
// @Param        limit  query     int  false  "Items per page (default 10, min 1, max 100)" default(10)
// @Success      200    {object}  response.APIResponse{data=response.PaginatedData}    "Bookings retrieved successfully"
// @Failure      400    {object}  response.APIResponse                                 "Bad request"
// @Failure      401    {object}  response.APIResponse                                 "Unauthorized"
// @Failure      500    {object}  response.APIResponse                                 "Internal server error"
// @Router       /users/me/bookings [get]
func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", domain.ErrUnauthorized)
		return
	}

	var query PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "Invalid pagination query", err.Error())
		return
	}

	bookings, total, err := h.bookingService.GetUserBookings(c.Request.Context(), userID, query.Page, query.Limit)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve bookings", err)
		return
	}

	items := toBookingResponseList(bookings)
	pagination := response.NewPaginationMeta(query.Page, query.Limit, total)
	response.SuccessWithPagination(c, http.StatusOK, "Bookings retrieved successfully", items, pagination)
}

// GetMyTicketQuota godoc
// @Summary      Get current user's quota for a ticket
// @Description  Retrieve the booking limit, current booked count, and remaining quota for the authenticated user
// @Tags         Bookings
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Ticket UUID"
// @Success      200  {object}  response.APIResponse{data=UserTicketQuotaResponse}  "User ticket quota retrieved successfully"
// @Failure      400  {object}  response.APIResponse                                "Invalid ticket ID"
// @Failure      401  {object}  response.APIResponse                                "Unauthorized"
// @Failure      404  {object}  response.APIResponse                                "Ticket not found"
// @Failure      500  {object}  response.APIResponse                                "Internal server error"
// @Router       /tickets/{id}/my-quota [get]
func (h *BookingHandler) GetMyTicketQuota(c *gin.Context) {
	ticketIDStr := c.Param("id")
	ticketID, err := uuid.Parse(ticketIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid ticket ID format", err.Error())
		return
	}

	userID, ok := middleware.GetAuthUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", domain.ErrUnauthorized)
		return
	}

	quota, err := h.bookingService.GetUserTicketQuota(c.Request.Context(), userID, ticketID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTicketNotFound):
			response.NotFound(c, "Ticket not found", err.Error())
		case errors.Is(err, domain.ErrUnauthorized):
			response.Error(c, http.StatusUnauthorized, "Unauthorized", err.Error())
		default:
			response.InternalServerError(c, "Failed to retrieve user ticket quota", err)
		}
		return
	}

	response.Success(c, http.StatusOK, "User ticket quota retrieved successfully", toUserTicketQuotaResponse(quota))
}
