package models

import (
	"time"
)

// HolidayType represents the type of holiday
type HolidayType string

const (
	// NationalHoliday represents "Libur Nasional" 
	NationalHoliday HolidayType = "national"
	// CollectiveLeave represents "Cuti Bersama"
	CollectiveLeave HolidayType = "collective_leave"
)

// Holiday represents a holiday record
type Holiday struct {
	ID          int         `json:"id" db:"id"`
	Name        string      `json:"name" db:"name" validate:"required,min=3,max=255"`
	Date        time.Time   `json:"date" db:"date" validate:"required"`
	Type        HolidayType `json:"type" db:"type" validate:"required,oneof=national collective_leave"`
	Description string      `json:"description" db:"description" validate:"max=1000"`
	IsActive    bool        `json:"is_active" db:"is_active"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateHolidayRequest represents request to create a holiday
type CreateHolidayRequest struct {
	Name        string      `json:"name" validate:"required,min=3,max=255"`
	Date        string      `json:"date" validate:"required"` // Format: YYYY-MM-DD
	Type        HolidayType `json:"type" validate:"required,oneof=national collective_leave"`
	Description string      `json:"description" validate:"max=1000"`
}

// UpdateHolidayRequest represents request to update a holiday
type UpdateHolidayRequest struct {
	Name        *string      `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Date        *string      `json:"date,omitempty"` // Format: YYYY-MM-DD
	Type        *HolidayType `json:"type,omitempty" validate:"omitempty,oneof=national collective_leave"`
	Description *string      `json:"description,omitempty" validate:"omitempty,max=1000"`
	IsActive    *bool        `json:"is_active,omitempty"`
}

// HolidayFilter represents filters for querying holidays
type HolidayFilter struct {
	Year      *int         `json:"year,omitempty"`
	Month     *int         `json:"month,omitempty"`
	Day       *int         `json:"day,omitempty"`
	Type      *HolidayType `json:"type,omitempty"`
	IsActive  *bool        `json:"is_active,omitempty"`
	StartDate *time.Time   `json:"start_date,omitempty"`
	EndDate   *time.Time   `json:"end_date,omitempty"`
	Limit     int          `json:"limit,omitempty"`
	Offset    int          `json:"offset,omitempty"`
}

// HolidayResponse represents the API response for holidays
type HolidayResponse struct {
	Data       []Holiday `json:"data"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	TotalPages int       `json:"total_pages"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
