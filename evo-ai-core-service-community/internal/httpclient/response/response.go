package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SuccessResponseData represents the standard success response structure
type SuccessResponseData struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    MetaInfo    `json:"meta"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponseData represents the standard error response structure
type ErrorResponseData struct {
	Success bool      `json:"success"`
	Error   ErrorInfo `json:"error"`
	Meta    MetaInfo  `json:"meta"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// MetaInfo contains metadata about the response
type MetaInfo struct {
	Timestamp  string          `json:"timestamp"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	Page            int   `json:"page"`
	PageSize        int   `json:"page_size"`
	TotalPages      int   `json:"total_pages"`
	Total           int   `json:"total"`
	HasNextPage     *bool `json:"has_next_page"`
	HasPreviousPage *bool `json:"has_previous_page"`
}

// SuccessResponse sends a standardized success response
func SuccessResponse(c *gin.Context, data interface{}, message string, status ...int) {
	statusCode := http.StatusOK
	if len(status) > 0 {
		statusCode = status[0]
	}

	response := SuccessResponseData{
		Success: true,
		Data:    data,
		Meta: MetaInfo{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	if message != "" {
		response.Message = message
	}

	c.JSON(statusCode, response)
}

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, code string, message string, details interface{}, status int) {
	response := ErrorResponseData{
		Success: false,
		Error: ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: MetaInfo{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	c.JSON(status, response)
}

// PaginatedResponse sends a standardized paginated response
func PaginatedResponse(c *gin.Context, data interface{}, page, pageSize, totalItems int, message string, status ...int) {
	statusCode := http.StatusOK
	if len(status) > 0 {
		statusCode = status[0]
	}

	totalPages := totalItems / pageSize
	if totalItems%pageSize != 0 {
		totalPages++
	}

	if totalPages == 0 {
		totalPages = 1
	}

	response := SuccessResponseData{
		Success: true,
		Data:    data,
		Meta: MetaInfo{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Pagination: &PaginationMeta{
				Page:            page,
				PageSize:        pageSize,
				TotalPages:      totalPages,
				Total:           totalItems,
				HasNextPage:     &[]bool{page < totalPages}[0],
				HasPreviousPage: &[]bool{page > 1}[0],
			},
		},
	}

	if message != "" {
		response.Message = message
	}

	c.JSON(statusCode, response)
}
