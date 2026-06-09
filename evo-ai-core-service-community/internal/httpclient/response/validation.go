package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse sends a standardized validation error response
func ValidationErrorResponse(c *gin.Context, err error) {
	var fields []FieldError

	// Check if it's a validator.ValidationErrors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fields = append(fields, FieldError{
				Field:   fieldError.Field(),
				Message: getErrorMessage(fieldError),
			})
		}
	} else {
		// Generic validation error
		fields = append(fields, FieldError{
			Field:   "unknown",
			Message: err.Error(),
		})
	}

	details := gin.H{
		"fields": fields,
	}

	ErrorResponse(c, "VALIDATION_ERROR", "Validation failed", details, http.StatusBadRequest)
}

// getErrorMessage returns a human-readable error message for a validation error
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", fe.Param())
	case "uuid":
		return "Must be a valid UUID"
	case "url":
		return "Must be a valid URL"
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s'", fe.Tag())
	}
}
