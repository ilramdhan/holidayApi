package services

import (
	"fmt"
	"time"

	"holidayapi/internal/models"
	"holidayapi/internal/repository"
)

// AuditService handles audit logging operations
type AuditService interface {
	LogAction(userID *int, username string, action models.AuditAction, resource models.AuditResource, resourceID *int, details, ipAddress, userAgent string, success bool) error
	GetAuditLogs(filter models.AuditLogFilter) (*models.AuditLogResponse, error)
	GetUserAuditLogs(userID int, limit, offset int) ([]models.AuditLog, error)
	CleanupOldLogs(daysToKeep int) error
}

// auditService implements AuditService
type auditService struct {
	auditRepo repository.AuditRepository
}

// NewAuditService creates a new audit service
func NewAuditService(auditRepo repository.AuditRepository) AuditService {
	return &auditService{
		auditRepo: auditRepo,
	}
}

// LogAction logs an audit action
func (s *auditService) LogAction(userID *int, username string, action models.AuditAction, resource models.AuditResource, resourceID *int, details, ipAddress, userAgent string, success bool) error {
	auditLog := &models.AuditLog{
		UserID:     userID,
		Username:   username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Success:    success,
	}

	if err := s.auditRepo.Create(auditLog); err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetAuditLogs retrieves audit logs with filters and pagination
func (s *auditService) GetAuditLogs(filter models.AuditLogFilter) (*models.AuditLogResponse, error) {
	// Set default pagination
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	logs, total, err := s.auditRepo.GetAll(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	page := (filter.Offset / filter.Limit) + 1
	totalPages := (total + filter.Limit - 1) / filter.Limit

	return &models.AuditLogResponse{
		Data:       logs,
		Total:      total,
		Page:       page,
		PerPage:    filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetUserAuditLogs retrieves audit logs for a specific user
func (s *auditService) GetUserAuditLogs(userID int, limit, offset int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	logs, err := s.auditRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user audit logs: %w", err)
	}

	return logs, nil
}

// CleanupOldLogs removes audit logs older than specified days
func (s *auditService) CleanupOldLogs(daysToKeep int) error {
	if daysToKeep <= 0 {
		daysToKeep = 90 // Default to 90 days
	}

	// Calculate cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)

	if err := s.auditRepo.DeleteOldLogs(cutoffDate); err != nil {
		return fmt.Errorf("failed to cleanup old audit logs: %w", err)
	}

	return nil
}
