package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"holidayapi/internal/models"
)

// AdminAuthMiddleware validates admin API key
func AdminAuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		providedKey := c.GetHeader("X-API-Key")
		
		if providedKey == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "API key is required",
			})
			c.Abort()
			return
		}

		if providedKey != apiKey {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
