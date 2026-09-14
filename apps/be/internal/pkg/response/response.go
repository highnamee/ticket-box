package response

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the unified envelope for all JSON responses
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// PaginationMeta provides pagination metadata inside the Data payload
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginationMeta(page, limit int, totalItems int64) *PaginationMeta {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalItems == 0 {
		totalPages = 0
	}
	return &PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

// PaginatedData wraps paginated items and metadata inside APIResponse.Data
type PaginatedData struct {
	Items      interface{}     `json:"items"`
	Pagination *PaginationMeta `json:"pagination"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithPagination(c *gin.Context, statusCode int, message string, items interface{}, pagination *PaginationMeta) {
	Success(c, statusCode, message, PaginatedData{
		Items:      items,
		Pagination: pagination,
	})
}

func Error(c *gin.Context, statusCode int, message string, err interface{}) {
	var errVal interface{} = err
	if e, ok := err.(error); ok {
		errVal = e.Error()
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   errVal,
	})
}

func BadRequest(c *gin.Context, message string, err interface{}) {
	if message == "" {
		message = "Bad request"
	}
	Error(c, http.StatusBadRequest, message, err)
}

func NotFound(c *gin.Context, message string, err interface{}) {
	if message == "" {
		message = "Resource not found"
	}
	Error(c, http.StatusNotFound, message, err)
}

func InternalServerError(c *gin.Context, message string, err interface{}) {
	if message == "" {
		message = "Internal server error"
	}
	Error(c, http.StatusInternalServerError, message, err)
}
