package models

import (
	"time"
)

// AuditAction represents audit action types
type AuditAction string

const (
	// Authentication actions
	ActionLogin        AuditAction = "LOGIN"
	ActionLogout       AuditAction = "LOGOUT"
	ActionLoginFailed  AuditAction = "LOGIN_FAILED"
	ActionTokenRefresh AuditAction = "TOKEN_REFRESH"

	// User management actions
	ActionUserCreate     AuditAction = "USER_CREATE"
	ActionUserUpdate     AuditAction = "USER_UPDATE"
	ActionUserDelete     AuditAction = "USER_DELETE"
	ActionPasswordChange AuditAction = "PASSWORD_CHANGE"

	// Holiday management actions
	ActionHolidayCreate AuditAction = "HOLIDAY_CREATE"
	ActionHolidayUpdate AuditAction = "HOLIDAY_UPDATE"
	ActionHolidayDelete AuditAction = "HOLIDAY_DELETE"
	ActionHolidayView   AuditAction = "HOLIDAY_VIEW"

	// System actions
	ActionSystemAccess AuditAction = "SYSTEM_ACCESS"
	ActionConfigChange AuditAction = "CONFIG_CHANGE"
)

// AuditResource represents audit resource types
type AuditResource string

const (
	ResourceAuth    AuditResource = "auth"
	ResourceUser    AuditResource = "user"
	ResourceHoliday AuditResource = "holiday"
	ResourceSystem  AuditResource = "system"
)

// AuditLog represents audit log entry
type AuditLog struct {
	ID         int           `json:"id" db:"id"`
	UserID     *int          `json:"user_id,omitempty" db:"user_id"`
	Username   string        `json:"username" db:"username"`
	Action     AuditAction   `json:"action" db:"action"`
	Resource   AuditResource `json:"resource" db:"resource"`
	ResourceID *int          `json:"resource_id,omitempty" db:"resource_id"`
	Details    string        `json:"details" db:"details"`
	IPAddress  string        `json:"ip_address" db:"ip_address"`
	UserAgent  string        `json:"user_agent" db:"user_agent"`
	Success    bool          `json:"success" db:"success"`
	CreatedAt  time.Time     `json:"created_at" db:"created_at"`
}

// AuditLogFilter represents filters for audit log queries
type AuditLogFilter struct {
	UserID    *int           `json:"user_id,omitempty"`
	Action    *AuditAction   `json:"action,omitempty"`
	Resource  *AuditResource `json:"resource,omitempty"`
	Success   *bool          `json:"success,omitempty"`
	StartDate *time.Time     `json:"start_date,omitempty"`
	EndDate   *time.Time     `json:"end_date,omitempty"`
	Limit     int            `json:"limit,omitempty"`
	Offset    int            `json:"offset,omitempty"`
}

// AuditLogResponse represents paginated audit log response
type AuditLogResponse struct {
	Data       []AuditLog `json:"data"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
}
