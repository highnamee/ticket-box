package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, http.StatusOK, "Operation successful", map[string]string{"key": "value"})

	assert.Equal(t, http.StatusOK, w.Code)

	var res APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.True(t, res.Success)
	assert.Equal(t, "Operation successful", res.Message)
	assert.NotNil(t, res.Data)
	assert.Nil(t, res.Error)
}

func TestSuccessWithPaginationResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	meta := NewPaginationMeta(2, 10, 25)
	SuccessWithPagination(c, http.StatusOK, "Tickets retrieved", []string{"ticket1", "ticket2"}, meta)

	assert.Equal(t, http.StatusOK, w.Code)

	var res struct {
		Success bool          `json:"success"`
		Message string        `json:"message"`
		Data    PaginatedData `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.True(t, res.Success)
	assert.Equal(t, "Tickets retrieved", res.Message)
	require.NotNil(t, res.Data.Pagination)
	assert.Equal(t, 2, res.Data.Pagination.Page)
	assert.Equal(t, 10, res.Data.Pagination.Limit)
	assert.Equal(t, int64(25), res.Data.Pagination.TotalItems)
	assert.Equal(t, 3, res.Data.Pagination.TotalPages)
}

func TestPaginationMeta_EdgeCases(t *testing.T) {
	// Zero items
	m0 := NewPaginationMeta(1, 10, 0)
	assert.Equal(t, 0, m0.TotalPages)

	// Invalid page and limit defaults
	mInvalid := NewPaginationMeta(0, 0, 15)
	assert.Equal(t, 1, mInvalid.Page)
	assert.Equal(t, 10, mInvalid.Limit)
	assert.Equal(t, 2, mInvalid.TotalPages)
}

func TestErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	BadRequest(c, "Invalid input", "Field 'email' is required")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.False(t, res.Success)
	assert.Equal(t, "Invalid input", res.Message)
	assert.Equal(t, "Field 'email' is required", res.Error)
}

func TestNotFoundResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	NotFound(c, "Resource not found", errors.New("record not found"))

	assert.Equal(t, http.StatusNotFound, w.Code)

	var res APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.False(t, res.Success)
	assert.Equal(t, "Resource not found", res.Message)
	assert.Equal(t, "record not found", res.Error)
}

func TestInternalServerErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	InternalServerError(c, "", errors.New("unexpected error"))

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.False(t, res.Success)
	assert.Equal(t, "Internal server error", res.Message)
	assert.Nil(t, res.Error)
}
