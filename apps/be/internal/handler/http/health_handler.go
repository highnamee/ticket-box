package http

import (
	"net/http"
	"time"

	"ticket-box-be/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(c, http.StatusOK, "Ticket Box API is healthy", gin.H{
		"status":    "UP",
		"timestamp": time.Now().UTC(),
	})
}
