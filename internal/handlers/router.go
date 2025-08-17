package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"holidayapi/internal/config"
	"holidayapi/internal/middleware"
	"holidayapi/internal/services"
)

// SetupRouter sets up the HTTP router with all routes and middleware
func SetupRouter(cfg *config.Config, holidayService services.HolidayService, authService services.AuthService, jwtService services.JWTService, auditService services.AuditService) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Global middleware
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.InputSanitizationMiddleware())
	router.Use(middleware.SQLInjectionProtectionMiddleware())
	router.Use(middleware.XSSProtectionMiddleware())
	router.Use(gin.Recovery())

	// Rate limiting
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.RequestsPerMinute, cfg.RateLimit.BurstSize)
	router.Use(rateLimiter.RateLimitMiddleware())

	// Initialize handlers
	holidayHandler := NewHolidayHandler(holidayService)
	adminHandler := NewAdminHandler(holidayService)
	authHandler := NewAuthHandler(authService)
	auditHandler := NewAuditHandler(auditService)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": gin.H{"unix": gin.H{"seconds": gin.H{"value": "1692025200"}}},
			"service":   "Holiday API Indonesia",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication endpoints (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)

			// Protected auth endpoints
			authProtected := auth.Group("")
			authProtected.Use(middleware.JWTAuthMiddleware(jwtService))
			{
				authProtected.GET("/profile", authHandler.GetProfile)
				authProtected.POST("/change-password", authHandler.ChangePassword)
				authProtected.GET("/audit-logs", auditHandler.GetMyAuditLogs)

				// Super admin only endpoints
				authProtected.POST("/register", middleware.RequireSuperAdmin(), authHandler.Register)
				authProtected.GET("/users", middleware.RequireAdminOrSuperAdmin(), authHandler.GetAllUsers)
				authProtected.DELETE("/users/:id", middleware.RequireSuperAdmin(), authHandler.DeleteUser)
			}
		}

		// Public holiday endpoints
		holidays := v1.Group("/holidays")
		{
			holidays.GET("", holidayHandler.GetHolidays)
			holidays.GET("/year/:year", holidayHandler.GetHolidaysByYear)
			holidays.GET("/month/:year/:month", holidayHandler.GetHolidaysByMonth)
			holidays.GET("/today", holidayHandler.GetHolidayToday)
			holidays.GET("/upcoming", holidayHandler.GetUpcomingHolidays)
			holidays.GET("/this-year", holidayHandler.GetHolidaysThisYear)
			holidays.GET("/this-month", holidayHandler.GetHolidaysThisMonth)
		}

		// Admin endpoints (JWT protected)
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuthMiddleware(jwtService))
		admin.Use(middleware.RequireAdminOrSuperAdmin())
		{
			// Holiday management
			admin.POST("/holidays", adminHandler.CreateHoliday)
			admin.GET("/holidays/:id", adminHandler.GetHoliday)
			admin.PUT("/holidays/:id", adminHandler.UpdateHoliday)
			admin.DELETE("/holidays/:id", adminHandler.DeleteHoliday)

			// Audit logs
			admin.GET("/audit-logs", auditHandler.GetAuditLogs)
			admin.GET("/audit-logs/user/:id", auditHandler.GetUserAuditLogs)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("doc.json")))

	return router
}
