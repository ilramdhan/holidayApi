package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"holidayapi/internal/models"
)

// AuditRepository interface defines audit log data access methods
type AuditRepository interface {
	Create(log *models.AuditLog) error
	GetAll(filter models.AuditLogFilter) ([]models.AuditLog, int, error)
	GetByUserID(userID int, limit, offset int) ([]models.AuditLog, error)
	DeleteOldLogs(olderThan time.Time) error
}

// auditRepository implements AuditRepository
type auditRepository struct {
	db *sql.DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository(db *sql.DB) AuditRepository {
	return &auditRepository{db: db}
}

// Create creates a new audit log entry
func (r *auditRepository) Create(log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (user_id, username, action, resource, resource_id, details, ip_address, user_agent, success, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	log.CreatedAt = time.Now()

	result, err := r.db.Exec(query, log.UserID, log.Username, log.Action, log.Resource,
		log.ResourceID, log.Details, log.IPAddress, log.UserAgent, log.Success, log.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	log.ID = int(id)
	return nil
}

// GetAll retrieves audit logs with filters
func (r *auditRepository) GetAll(filter models.AuditLogFilter) ([]models.AuditLog, int, error) {
	// Build WHERE clause
	whereConditions := []string{}
	args := []interface{}{}

	if filter.UserID != nil {
		whereConditions = append(whereConditions, "user_id = ?")
		args = append(args, *filter.UserID)
	}

	if filter.Action != nil {
		whereConditions = append(whereConditions, "action = ?")
		args = append(args, string(*filter.Action))
	}

	if filter.Resource != nil {
		whereConditions = append(whereConditions, "resource = ?")
		args = append(args, string(*filter.Resource))
	}

	if filter.Success != nil {
		whereConditions = append(whereConditions, "success = ?")
		args = append(args, *filter.Success)
	}

	if filter.StartDate != nil {
		whereConditions = append(whereConditions, "created_at >= ?")
		args = append(args, filter.StartDate.Format("2006-01-02 15:04:05"))
	}

	if filter.EndDate != nil {
		whereConditions = append(whereConditions, "created_at <= ?")
		args = append(args, filter.EndDate.Format("2006-01-02 15:04:05"))
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT id, user_id, username, action, resource, resource_id, details, ip_address, user_agent, success, created_at
		FROM audit_logs
		%s
		ORDER BY created_at DESC
	`, whereClause)

	// Add pagination
	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)

		if filter.Offset > 0 {
			query += " OFFSET ?"
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var userID sql.NullInt64
		var resourceID sql.NullInt64
		var ipAddress sql.NullString
		var userAgent sql.NullString
		var details sql.NullString

		err := rows.Scan(
			&log.ID, &userID, &log.Username, &log.Action, &log.Resource,
			&resourceID, &details, &ipAddress, &userAgent, &log.Success, &log.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit log: %w", err)
		}

		if userID.Valid {
			uid := int(userID.Int64)
			log.UserID = &uid
		}

		if resourceID.Valid {
			rid := int(resourceID.Int64)
			log.ResourceID = &rid
		}

		if details.Valid {
			log.Details = details.String
		}

		if ipAddress.Valid {
			log.IPAddress = ipAddress.String
		}

		if userAgent.Valid {
			log.UserAgent = userAgent.String
		}

		logs = append(logs, log)
	}

	return logs, total, nil
}

// GetByUserID retrieves audit logs for a specific user
func (r *auditRepository) GetByUserID(userID int, limit, offset int) ([]models.AuditLog, error) {
	query := `
		SELECT id, user_id, username, action, resource, resource_id, details, ip_address, user_agent, success, created_at
		FROM audit_logs
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs by user: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var userIDNull sql.NullInt64
		var resourceID sql.NullInt64
		var ipAddress sql.NullString
		var userAgent sql.NullString
		var details sql.NullString

		err := rows.Scan(
			&log.ID, &userIDNull, &log.Username, &log.Action, &log.Resource,
			&resourceID, &details, &ipAddress, &userAgent, &log.Success, &log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		if userIDNull.Valid {
			uid := int(userIDNull.Int64)
			log.UserID = &uid
		}

		if resourceID.Valid {
			rid := int(resourceID.Int64)
			log.ResourceID = &rid
		}

		if details.Valid {
			log.Details = details.String
		}

		if ipAddress.Valid {
			log.IPAddress = ipAddress.String
		}

		if userAgent.Valid {
			log.UserAgent = userAgent.String
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// DeleteOldLogs deletes audit logs older than specified time
func (r *auditRepository) DeleteOldLogs(olderThan time.Time) error {
	query := `DELETE FROM audit_logs WHERE created_at < ?`

	result, err := r.db.Exec(query, olderThan.Format("2006-01-02 15:04:05"))
	if err != nil {
		return fmt.Errorf("failed to delete old audit logs: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	fmt.Printf("Deleted %d old audit log entries\n", rowsAffected)
	return nil
}
