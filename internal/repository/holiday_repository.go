package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"holidayapi/internal/models"
)

// HolidayRepository interface defines holiday data access methods
type HolidayRepository interface {
	Create(holiday *models.Holiday) error
	GetByID(id int) (*models.Holiday, error)
	GetAll(filter models.HolidayFilter) ([]models.Holiday, int, error)
	Update(id int, holiday *models.Holiday) error
	Delete(id int) error
	GetByDate(date time.Time) (*models.Holiday, error)
	GetByDateRange(startDate, endDate time.Time, holidayType *models.HolidayType) ([]models.Holiday, error)
}

// holidayRepository implements HolidayRepository
type holidayRepository struct {
	db *sql.DB
}

// NewHolidayRepository creates a new holiday repository
func NewHolidayRepository(db *sql.DB) HolidayRepository {
	return &holidayRepository{db: db}
}

// Create creates a new holiday
func (r *holidayRepository) Create(holiday *models.Holiday) error {
	query := `
		INSERT INTO holidays (name, date, type, description, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	holiday.CreatedAt = now
	holiday.UpdatedAt = now
	holiday.IsActive = true

	result, err := r.db.Exec(query, holiday.Name, holiday.Date, holiday.Type, 
		holiday.Description, holiday.IsActive, holiday.CreatedAt, holiday.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create holiday: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	holiday.ID = int(id)
	return nil
}

// GetByID retrieves a holiday by ID
func (r *holidayRepository) GetByID(id int) (*models.Holiday, error) {
	query := `
		SELECT id, name, date, type, description, is_active, created_at, updated_at
		FROM holidays
		WHERE id = ? AND is_active = TRUE
	`

	holiday := &models.Holiday{}
	err := r.db.QueryRow(query, id).Scan(
		&holiday.ID, &holiday.Name, &holiday.Date, &holiday.Type,
		&holiday.Description, &holiday.IsActive, &holiday.CreatedAt, &holiday.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("holiday not found")
		}
		return nil, fmt.Errorf("failed to get holiday: %w", err)
	}

	return holiday, nil
}

// GetAll retrieves holidays with filters
func (r *holidayRepository) GetAll(filter models.HolidayFilter) ([]models.Holiday, int, error) {
	// Build WHERE clause
	whereConditions := []string{"is_active = TRUE"}
	args := []interface{}{}

	if filter.Year != nil {
		whereConditions = append(whereConditions, "strftime('%Y', date) = ?")
		args = append(args, fmt.Sprintf("%d", *filter.Year))
	}

	if filter.Month != nil {
		whereConditions = append(whereConditions, "strftime('%m', date) = ?")
		args = append(args, fmt.Sprintf("%02d", *filter.Month))
	}

	if filter.Day != nil {
		whereConditions = append(whereConditions, "strftime('%d', date) = ?")
		args = append(args, fmt.Sprintf("%02d", *filter.Day))
	}

	if filter.Type != nil {
		whereConditions = append(whereConditions, "type = ?")
		args = append(args, string(*filter.Type))
	}

	if filter.StartDate != nil {
		whereConditions = append(whereConditions, "date >= ?")
		args = append(args, filter.StartDate.Format("2006-01-02"))
	}

	if filter.EndDate != nil {
		whereConditions = append(whereConditions, "date <= ?")
		args = append(args, filter.EndDate.Format("2006-01-02"))
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM holidays WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count holidays: %w", err)
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT id, name, date, type, description, is_active, created_at, updated_at
		FROM holidays
		WHERE %s
		ORDER BY date ASC
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
		return nil, 0, fmt.Errorf("failed to query holidays: %w", err)
	}
	defer rows.Close()

	var holidays []models.Holiday
	for rows.Next() {
		var holiday models.Holiday
		err := rows.Scan(
			&holiday.ID, &holiday.Name, &holiday.Date, &holiday.Type,
			&holiday.Description, &holiday.IsActive, &holiday.CreatedAt, &holiday.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan holiday: %w", err)
		}
		holidays = append(holidays, holiday)
	}

	return holidays, total, nil
}

// Update updates a holiday
func (r *holidayRepository) Update(id int, holiday *models.Holiday) error {
	query := `
		UPDATE holidays 
		SET name = ?, date = ?, type = ?, description = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`

	holiday.UpdatedAt = time.Now()

	result, err := r.db.Exec(query, holiday.Name, holiday.Date, holiday.Type,
		holiday.Description, holiday.IsActive, holiday.UpdatedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update holiday: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("holiday not found")
	}

	return nil
}

// Delete soft deletes a holiday
func (r *holidayRepository) Delete(id int) error {
	query := `UPDATE holidays SET is_active = FALSE, updated_at = ? WHERE id = ?`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete holiday: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("holiday not found")
	}

	return nil
}

// GetByDate retrieves holiday by specific date
func (r *holidayRepository) GetByDate(date time.Time) (*models.Holiday, error) {
	query := `
		SELECT id, name, date, type, description, is_active, created_at, updated_at
		FROM holidays
		WHERE date = ? AND is_active = TRUE
	`

	holiday := &models.Holiday{}
	err := r.db.QueryRow(query, date.Format("2006-01-02")).Scan(
		&holiday.ID, &holiday.Name, &holiday.Date, &holiday.Type,
		&holiday.Description, &holiday.IsActive, &holiday.CreatedAt, &holiday.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No holiday found, not an error
		}
		return nil, fmt.Errorf("failed to get holiday by date: %w", err)
	}

	return holiday, nil
}

// GetByDateRange retrieves holidays within date range
func (r *holidayRepository) GetByDateRange(startDate, endDate time.Time, holidayType *models.HolidayType) ([]models.Holiday, error) {
	whereConditions := []string{"is_active = TRUE", "date >= ?", "date <= ?"}
	args := []interface{}{startDate.Format("2006-01-02"), endDate.Format("2006-01-02")}

	if holidayType != nil {
		whereConditions = append(whereConditions, "type = ?")
		args = append(args, string(*holidayType))
	}

	query := fmt.Sprintf(`
		SELECT id, name, date, type, description, is_active, created_at, updated_at
		FROM holidays
		WHERE %s
		ORDER BY date ASC
	`, strings.Join(whereConditions, " AND "))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query holidays by date range: %w", err)
	}
	defer rows.Close()

	var holidays []models.Holiday
	for rows.Next() {
		var holiday models.Holiday
		err := rows.Scan(
			&holiday.ID, &holiday.Name, &holiday.Date, &holiday.Type,
			&holiday.Description, &holiday.IsActive, &holiday.CreatedAt, &holiday.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan holiday: %w", err)
		}
		holidays = append(holidays, holiday)
	}

	return holidays, nil
}
