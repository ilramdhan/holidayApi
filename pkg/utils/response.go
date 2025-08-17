package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"holidayapi/internal/models"
)

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message, error string) {
	c.JSON(statusCode, models.ErrorResponse{
		Success: false,
		Message: message,
		Error:   error,
	})
}

// CreatedResponse sends a created response
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}
