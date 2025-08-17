package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"holidayapi/internal/models"
	"holidayapi/internal/services"
)

// HolidayHandler handles holiday-related HTTP requests
type HolidayHandler struct {
	service   services.HolidayService
	validator *validator.Validate
}

// NewHolidayHandler creates a new holiday handler
func NewHolidayHandler(service services.HolidayService) *HolidayHandler {
	return &HolidayHandler{
		service:   service,
		validator: validator.New(),
	}
}

// GetHolidays godoc
// @Summary Get holidays with filters
// @Description Get holidays with optional filters (year, month, type, etc.)
// @Tags holidays
// @Accept json
// @Produce json
// @Param year query int false "Year filter"
// @Param month query int false "Month filter (1-12)"
// @Param type query string false "Holiday type" Enums(national, collective_leave)
// @Param limit query int false "Limit results (max 100)" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} models.HolidayResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/holidays [get]
func (h *HolidayHandler) GetHolidays(c *gin.Context) {
	filter := models.HolidayFilter{}

	// Parse query parameters
	if yearStr := c.Query("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			filter.Year = &year
		}
	}

	if monthStr := c.Query("month"); monthStr != "" {
		if month, err := strconv.Atoi(monthStr); err == nil && month >= 1 && month <= 12 {
			filter.Month = &month
		}
	}

	if typeStr := c.Query("type"); typeStr != "" {
		holidayType := models.HolidayType(typeStr)
		if holidayType == models.NationalHoliday || holidayType == models.CollectiveLeave {
			filter.Type = &holidayType
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	response, err := h.service.GetHolidays(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to get holidays",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Holidays retrieved successfully",
		Data:    response,
	})
}

// GetHolidaysByYear godoc
// @Summary Get holidays by year
// @Description Get all holidays for a specific year
// @Tags holidays
// @Accept json
// @Produce json
// @Param year path int true "Year"
// @Param type query string false "Holiday type" Enums(national, collective_leave)
// @Success 200 {object} models.APIResponse{data=[]models.Holiday}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/holidays/year/{year} [get]
func (h *HolidayHandler) GetHolidaysByYear(c *gin.Context) {
	yearStr := c.Param("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid year parameter",
			Error:   "Year must be a valid integer",
		})
		return
	}

	holidays, err := h.service.GetHolidaysByYear(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to get holidays",
			Error:   err.Error(),
		})
		return
	}

	// Filter by type if specified
	if typeStr := c.Query("type"); typeStr != "" {
		holidayType := models.HolidayType(typeStr)
		if holidayType == models.NationalHoliday || holidayType == models.CollectiveLeave {
			var filtered []models.Holiday
			for _, holiday := range holidays {
				if holiday.Type == holidayType {
					filtered = append(filtered, holiday)
				}
			}
			holidays = filtered
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Holidays retrieved successfully",
		Data:    holidays,
	})
}

// GetHolidaysByMonth godoc
// @Summary Get holidays by month
// @Description Get all holidays for a specific month and year
// @Tags holidays
// @Accept json
// @Produce json
// @Param year path int true "Year"
// @Param month path int true "Month (1-12)"
// @Param type query string false "Holiday type" Enums(national, collective_leave)
// @Success 200 {object} models.APIResponse{data=[]models.Holiday}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/holidays/month/{year}/{month} [get]
func (h *HolidayHandler) GetHolidaysByMonth(c *gin.Context) {
	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid year parameter",
			Error:   "Year must be a valid integer",
		})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid month parameter",
			Error:   "Month must be between 1 and 12",
		})
		return
	}

	holidays, err := h.service.GetHolidaysByMonth(year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to get holidays",
			Error:   err.Error(),
		})
		return
	}

	// Filter by type if specified
	if typeStr := c.Query("type"); typeStr != "" {
		holidayType := models.HolidayType(typeStr)
		if holidayType == models.NationalHoliday || holidayType == models.CollectiveLeave {
			var filtered []models.Holiday
			for _, holiday := range holidays {
				if holiday.Type == holidayType {
					filtered = append(filtered, holiday)
				}
			}
			holidays = filtered
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Holidays retrieved successfully",
		Data:    holidays,
	})
}

// GetHolidayToday godoc
// @Summary Get today's holiday
// @Description Get holiday for today's date if any
// @Tags holidays
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=models.Holiday}
// @Success 204 "No holiday today"
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/holidays/today [get]
func (h *HolidayHandler) GetHolidayToday(c *gin.Context) {
	holiday, err := h.service.GetHolidayToday()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to get today's holiday",
			Error:   err.Error(),
		})
		return
	}

	if holiday == nil {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "No holiday today",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Today's holiday retrieved successfully",
		Data:    holiday,
	})
}

// GetUpcomingHolidays godoc
// @Summary Get upcoming holidays
// @Description Get upcoming holidays from today
// @Tags holidays
// @Accept json
// @Produce json
// @Param limit query int false "Limit results" default(10)
// @Param type query string false "Holiday type" Enums(national, collective_leave)
// @Success 200 {object} models.APIResponse{data=[]models.Holiday}
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/holidays/upcoming [get]
func (h *HolidayHandler) GetUpcomingHolidays(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	holidays, err := h.service.GetUpcomingHolidays(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to get upcoming holidays",
			Error:   err.Error(),
		})
		return
	}

	// Filter by type if specified
	if typeStr := c.Query("type"); typeStr != "" {
		holidayType := models.HolidayType(typeStr)
		if holidayType == models.NationalHoliday || holidayType == models.CollectiveLeave {
			var filtered []models.Holiday
			for _, holiday := range holidays {
				if holiday.Type == holidayType {
					filtered = append(filtered, holiday)
				}
			}
			holidays = filtered
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Upcoming holidays retrieved successfully",
		Data:    holidays,
	})
}

// GetHolidaysThisYear godoc
// @Summary Get holidays for current year
// @Description Get all holidays for the current year
// @Tags holidays
// @Accept json
// @Produce json
// @Param type query string false "Holiday type" Enums(national, collective_leave)
// @Success 200 {object} models.APIResponse{data=[]models.Holiday}
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/holidays/this-year [get]
func (h *HolidayHandler) GetHolidaysThisYear(c *gin.Context) {
	holidays, err := h.service.GetHolidaysThisYear()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to get holidays for this year",
			Error:   err.Error(),
		})
		return
	}

	// Filter by type if specified
	if typeStr := c.Query("type"); typeStr != "" {
		holidayType := models.HolidayType(typeStr)
		if holidayType == models.NationalHoliday || holidayType == models.CollectiveLeave {
			var filtered []models.Holiday
			for _, holiday := range holidays {
				if holiday.Type == holidayType {
					filtered = append(filtered, holiday)
				}
			}
			holidays = filtered
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Holidays for this year retrieved successfully",
		Data:    holidays,
	})
}

// GetHolidaysThisMonth godoc
// @Summary Get holidays for current month
// @Description Get all holidays for the current month
// @Tags holidays
// @Accept json
// @Produce json
// @Param type query string false "Holiday type" Enums(national, collective_leave)
// @Success 200 {object} models.APIResponse{data=[]models.Holiday}
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/holidays/this-month [get]
func (h *HolidayHandler) GetHolidaysThisMonth(c *gin.Context) {
	holidays, err := h.service.GetHolidaysThisMonth()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to get holidays for this month",
			Error:   err.Error(),
		})
		return
	}

	// Filter by type if specified
	if typeStr := c.Query("type"); typeStr != "" {
		holidayType := models.HolidayType(typeStr)
		if holidayType == models.NationalHoliday || holidayType == models.CollectiveLeave {
			var filtered []models.Holiday
			for _, holiday := range holidays {
				if holiday.Type == holidayType {
					filtered = append(filtered, holiday)
				}
			}
			holidays = filtered
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Holidays for this month retrieved successfully",
		Data:    holidays,
	})
}
