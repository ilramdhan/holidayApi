// Package main provides the entry point for the Holiday API server
// @title Holiday API Indonesia
// @version 2.0
// @description API untuk mendapatkan informasi hari libur nasional dan cuti bersama Indonesia berdasarkan SKB 3 Menteri dengan sistem authentication JWT yang lengkap
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description Legacy Admin API Key for backward compatibility

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer Token. Format: Bearer {token}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "holidayapi/docs" // Import for swagger docs
	"holidayapi/internal/config"
	"holidayapi/internal/database"
	"holidayapi/internal/handlers"
	"holidayapi/internal/repository"
	"holidayapi/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewConnection(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(cfg.Database.MigrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	holidayRepo := repository.NewHolidayRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	auditRepo := repository.NewAuditRepository(db.DB)

	// Initialize services
	jwtService := services.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	auditService := services.NewAuditService(auditRepo)
	authService := services.NewAuthService(userRepo, auditRepo, jwtService)
	holidayService := services.NewHolidayService(holidayRepo)

	// Setup router
	router := handlers.SetupRouter(cfg, holidayService, authService, jwtService, auditService)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Swagger documentation available at: http://%s:%s/swagger/index.html", cfg.Server.Host, cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
