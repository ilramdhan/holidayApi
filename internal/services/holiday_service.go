package services

import (
	"fmt"
	"time"

	"holidayapi/internal/models"
	"holidayapi/internal/repository"
)

// HolidayService interface defines holiday business logic methods
type HolidayService interface {
	CreateHoliday(req models.CreateHolidayRequest) (*models.Holiday, error)
	GetHolidayByID(id int) (*models.Holiday, error)
	GetHolidays(filter models.HolidayFilter) (*models.HolidayResponse, error)
	UpdateHoliday(id int, req models.UpdateHolidayRequest) (*models.Holiday, error)
	DeleteHoliday(id int) error
	GetHolidaysThisYear() ([]models.Holiday, error)
	GetHolidaysThisMonth() ([]models.Holiday, error)
	GetHolidayToday() (*models.Holiday, error)
	GetUpcomingHolidays(limit int) ([]models.Holiday, error)
	GetHolidaysByYear(year int) ([]models.Holiday, error)
	GetHolidaysByMonth(year, month int) ([]models.Holiday, error)
	GetHolidaysByType(holidayType models.HolidayType) ([]models.Holiday, error)
}

// holidayService implements HolidayService
type holidayService struct {
	repo repository.HolidayRepository
}

// NewHolidayService creates a new holiday service
func NewHolidayService(repo repository.HolidayRepository) HolidayService {
	return &holidayService{repo: repo}
}

// CreateHoliday creates a new holiday
func (s *holidayService) CreateHoliday(req models.CreateHolidayRequest) (*models.Holiday, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
	}

	// Check if holiday already exists on this date
	existing, err := s.repo.GetByDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing holiday: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("holiday already exists on date %s", req.Date)
	}

	holiday := &models.Holiday{
		Name:        req.Name,
		Date:        date,
		Type:        req.Type,
		Description: req.Description,
	}

	if err := s.repo.Create(holiday); err != nil {
		return nil, fmt.Errorf("failed to create holiday: %w", err)
	}

	return holiday, nil
}

// GetHolidayByID retrieves a holiday by ID
func (s *holidayService) GetHolidayByID(id int) (*models.Holiday, error) {
	return s.repo.GetByID(id)
}

// GetHolidays retrieves holidays with filters and pagination
func (s *holidayService) GetHolidays(filter models.HolidayFilter) (*models.HolidayResponse, error) {
	// Set default pagination
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	holidays, total, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get holidays: %w", err)
	}

	page := (filter.Offset / filter.Limit) + 1
	totalPages := (total + filter.Limit - 1) / filter.Limit

	return &models.HolidayResponse{
		Data:       holidays,
		Total:      total,
		Page:       page,
		PerPage:    filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateHoliday updates a holiday
func (s *holidayService) UpdateHoliday(id int, req models.UpdateHolidayRequest) (*models.Holiday, error) {
	// Get existing holiday
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD: %w", err)
		}
		existing.Date = date
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := s.repo.Update(id, existing); err != nil {
		return nil, fmt.Errorf("failed to update holiday: %w", err)
	}

	return existing, nil
}

// DeleteHoliday deletes a holiday
func (s *holidayService) DeleteHoliday(id int) error {
	return s.repo.Delete(id)
}

// GetHolidaysThisYear gets holidays for current year
func (s *holidayService) GetHolidaysThisYear() ([]models.Holiday, error) {
	year := time.Now().Year()
	filter := models.HolidayFilter{Year: &year}
	holidays, _, err := s.repo.GetAll(filter)
	return holidays, err
}

// GetHolidaysThisMonth gets holidays for current month
func (s *holidayService) GetHolidaysThisMonth() ([]models.Holiday, error) {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	filter := models.HolidayFilter{Year: &year, Month: &month}
	holidays, _, err := s.repo.GetAll(filter)
	return holidays, err
}

// GetHolidayToday gets today's holiday if any
func (s *holidayService) GetHolidayToday() (*models.Holiday, error) {
	today := time.Now()
	return s.repo.GetByDate(today)
}

// GetUpcomingHolidays gets upcoming holidays
func (s *holidayService) GetUpcomingHolidays(limit int) ([]models.Holiday, error) {
	if limit <= 0 {
		limit = 10
	}
	
	today := time.Now()
	endDate := today.AddDate(1, 0, 0) // Next year
	
	return s.repo.GetByDateRange(today, endDate, nil)
}

// GetHolidaysByYear gets holidays by specific year
func (s *holidayService) GetHolidaysByYear(year int) ([]models.Holiday, error) {
	filter := models.HolidayFilter{Year: &year}
	holidays, _, err := s.repo.GetAll(filter)
	return holidays, err
}

// GetHolidaysByMonth gets holidays by specific month
func (s *holidayService) GetHolidaysByMonth(year, month int) ([]models.Holiday, error) {
	filter := models.HolidayFilter{Year: &year, Month: &month}
	holidays, _, err := s.repo.GetAll(filter)
	return holidays, err
}

// GetHolidaysByType gets holidays by type
func (s *holidayService) GetHolidaysByType(holidayType models.HolidayType) ([]models.Holiday, error) {
	filter := models.HolidayFilter{Type: &holidayType}
	holidays, _, err := s.repo.GetAll(filter)
	return holidays, err
}
