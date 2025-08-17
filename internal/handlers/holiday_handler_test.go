package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"holidayapi/internal/models"
	"holidayapi/internal/services"
)

// MockHolidayService is a mock implementation of HolidayService
type MockHolidayService struct {
	mock.Mock
}

func (m *MockHolidayService) CreateHoliday(req models.CreateHolidayRequest) (*models.Holiday, error) {
	args := m.Called(req)
	return args.Get(0).(*models.Holiday), args.Error(1)
}

func (m *MockHolidayService) GetHolidayByID(id int) (*models.Holiday, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Holiday), args.Error(1)
}

func (m *MockHolidayService) GetHolidays(filter models.HolidayFilter) (*models.HolidayResponse, error) {
	args := m.Called(filter)
	return args.Get(0).(*models.HolidayResponse), args.Error(1)
}

func (m *MockHolidayService) UpdateHoliday(id int, req models.UpdateHolidayRequest) (*models.Holiday, error) {
	args := m.Called(id, req)
	return args.Get(0).(*models.Holiday), args.Error(1)
}

func (m *MockHolidayService) DeleteHoliday(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockHolidayService) GetHolidaysThisYear() ([]models.Holiday, error) {
	args := m.Called()
	return args.Get(0).([]models.Holiday), args.Error(1)
}

func (m *MockHolidayService) GetHolidaysThisMonth() ([]models.Holiday, error) {
	args := m.Called()
	return args.Get(0).([]models.Holiday), args.Error(1)
}

func (m *MockHolidayService) GetHolidayToday() (*models.Holiday, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Holiday), args.Error(1)
}

func (m *MockHolidayService) GetUpcomingHolidays(limit int) ([]models.Holiday, error) {
	args := m.Called(limit)
	return args.Get(0).([]models.Holiday), args.Error(1)
}

func (m *MockHolidayService) GetHolidaysByYear(year int) ([]models.Holiday, error) {
	args := m.Called(year)
	return args.Get(0).([]models.Holiday), args.Error(1)
}

func (m *MockHolidayService) GetHolidaysByMonth(year, month int) ([]models.Holiday, error) {
	args := m.Called(year, month)
	return args.Get(0).([]models.Holiday), args.Error(1)
}

func (m *MockHolidayService) GetHolidaysByType(holidayType models.HolidayType) ([]models.Holiday, error) {
	args := m.Called(holidayType)
	return args.Get(0).([]models.Holiday), args.Error(1)
}

func TestHolidayHandler_GetHolidayToday(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*MockHolidayService)
		expectedStatus int
		expectedData   interface{}
	}{
		{
			name: "holiday found today",
			setupMock: func(m *MockHolidayService) {
				holiday := &models.Holiday{
					ID:   1,
					Name: "Christmas",
					Date: time.Now(),
					Type: models.NationalHoliday,
				}
				m.On("GetHolidayToday").Return(holiday, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "no holiday today",
			setupMock: func(m *MockHolidayService) {
				m.On("GetHolidayToday").Return((*models.Holiday)(nil), nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockHolidayService)
			tt.setupMock(mockService)

			handler := NewHolidayHandler(mockService)

			router := gin.New()
			router.GET("/holidays/today", handler.GetHolidayToday)

			req, _ := http.NewRequest("GET", "/holidays/today", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHolidayHandler_GetHolidaysByYear(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockHolidayService)
	expectedHolidays := []models.Holiday{
		{ID: 1, Name: "New Year", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	mockService.On("GetHolidaysByYear", 2024).Return(expectedHolidays, nil)

	handler := NewHolidayHandler(mockService)

	router := gin.New()
	router.GET("/holidays/year/:year", handler.GetHolidaysByYear)

	req, _ := http.NewRequest("GET", "/holidays/year/2024", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockService.AssertExpectations(t)
}
