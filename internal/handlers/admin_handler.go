package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"holidayapi/internal/models"
	"holidayapi/internal/services"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	service   services.HolidayService
	validator *validator.Validate
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(service services.HolidayService) *AdminHandler {
	return &AdminHandler{
		service:   service,
		validator: validator.New(),
	}
}

// CreateHoliday godoc
// @Summary Create a new holiday (Admin only)
// @Description Create a new holiday entry
// @Tags admin
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Admin API Key"
// @Param holiday body models.CreateHolidayRequest true "Holiday data"
// @Success 201 {object} models.APIResponse{data=models.Holiday}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/admin/holidays [post]
func (h *AdminHandler) CreateHoliday(c *gin.Context) {
	var req models.CreateHolidayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	holiday, err := h.service.CreateHoliday(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to create holiday",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Holiday created successfully",
		Data:    holiday,
	})
}

// GetHoliday godoc
// @Summary Get holiday by ID (Admin only)
// @Description Get a specific holiday by ID
// @Tags admin
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Admin API Key"
// @Param id path int true "Holiday ID"
// @Success 200 {object} models.APIResponse{data=models.Holiday}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/admin/holidays/{id} [get]
func (h *AdminHandler) GetHoliday(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid holiday ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	holiday, err := h.service.GetHolidayByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Message: "Holiday not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Holiday retrieved successfully",
		Data:    holiday,
	})
}

// UpdateHoliday godoc
// @Summary Update holiday (Admin only)
// @Description Update an existing holiday
// @Tags admin
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Admin API Key"
// @Param id path int true "Holiday ID"
// @Param holiday body models.UpdateHolidayRequest true "Holiday update data"
// @Success 200 {object} models.APIResponse{data=models.Holiday}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/admin/holidays/{id} [put]
func (h *AdminHandler) UpdateHoliday(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid holiday ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	var req models.UpdateHolidayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	holiday, err := h.service.UpdateHoliday(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to update holiday",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Holiday updated successfully",
		Data:    holiday,
	})
}

// DeleteHoliday godoc
// @Summary Delete holiday (Admin only)
// @Description Soft delete a holiday
// @Tags admin
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Admin API Key"
// @Param id path int true "Holiday ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/admin/holidays/{id} [delete]
func (h *AdminHandler) DeleteHoliday(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Message: "Invalid holiday ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	if err := h.service.DeleteHoliday(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Message: "Failed to delete holiday",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Holiday deleted successfully",
	})
}
