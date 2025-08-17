package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"holidayapi/internal/models"
	"holidayapi/internal/services"
)

// JWTAuthMiddleware validates JWT tokens
func JWTAuthMiddleware(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "Token is required",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(requiredRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "User role not found in context",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(models.UserRole)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Message: "Unauthorized",
				Error:   "Invalid user role format",
			})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasPermission := false
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Success: false,
				Message: "Forbidden",
				Error:   "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSuperAdmin middleware checks if user is super admin
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole(models.SuperAdminRole)
}

// RequireAdminOrSuperAdmin middleware checks if user is admin or super admin
func RequireAdminOrSuperAdmin() gin.HandlerFunc {
	return RequireRole(models.AdminRole, models.SuperAdminRole)
}

// GetCurrentUser helper function to get current user from context
func GetCurrentUser(c *gin.Context) (*models.JWTClaims, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, fmt.Errorf("user ID not found in context")
	}

	username, exists := c.Get("username")
	if !exists {
		return nil, fmt.Errorf("username not found in context")
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		return nil, fmt.Errorf("user role not found in context")
	}

	return &models.JWTClaims{
		UserID:   userID.(int),
		Username: username.(string),
		Role:     userRole.(models.UserRole),
	}, nil
}
