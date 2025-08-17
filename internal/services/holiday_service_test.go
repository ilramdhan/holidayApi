package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"holidayapi/internal/models"
)

// MockHolidayRepository is a mock implementation of HolidayRepository
type MockHolidayRepository struct {
	mock.Mock
}

func (m *MockHolidayRepository) Create(holiday *models.Holiday) error {
	args := m.Called(holiday)
	return args.Error(0)
}

func (m *MockHolidayRepository) GetByID(id int) (*models.Holiday, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Holiday), args.Error(1)
}

func (m *MockHolidayRepository) GetAll(filter models.HolidayFilter) ([]models.Holiday, int, error) {
	args := m.Called(filter)
	return args.Get(0).([]models.Holiday), args.Int(1), args.Error(2)
}

func (m *MockHolidayRepository) Update(id int, holiday *models.Holiday) error {
	args := m.Called(id, holiday)
	return args.Error(0)
}

func (m *MockHolidayRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockHolidayRepository) GetByDate(date time.Time) (*models.Holiday, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Holiday), args.Error(1)
}

func (m *MockHolidayRepository) GetByDateRange(startDate, endDate time.Time, holidayType *models.HolidayType) ([]models.Holiday, error) {
	args := m.Called(startDate, endDate, holidayType)
	return args.Get(0).([]models.Holiday), args.Error(1)
}

func TestHolidayService_CreateHoliday(t *testing.T) {
	mockRepo := new(MockHolidayRepository)
	service := NewHolidayService(mockRepo)

	tests := []struct {
		name        string
		request     models.CreateHolidayRequest
		setupMock   func()
		expectError bool
	}{
		{
			name: "successful creation",
			request: models.CreateHolidayRequest{
				Name:        "Test Holiday",
				Date:        "2024-12-25",
				Type:        models.NationalHoliday,
				Description: "Test description",
			},
			setupMock: func() {
				date, _ := time.Parse("2006-01-02", "2024-12-25")
				mockRepo.On("GetByDate", date).Return((*models.Holiday)(nil), nil)
				mockRepo.On("Create", mock.AnythingOfType("*models.Holiday")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "invalid date format",
			request: models.CreateHolidayRequest{
				Name: "Test Holiday",
				Date: "invalid-date",
				Type: models.NationalHoliday,
			},
			setupMock:   func() {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock()

			holiday, err := service.CreateHoliday(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, holiday)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, holiday)
				assert.Equal(t, tt.request.Name, holiday.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHolidayService_GetHolidayToday(t *testing.T) {
	mockRepo := new(MockHolidayRepository)
	service := NewHolidayService(mockRepo)

	today := time.Now()
	expectedHoliday := &models.Holiday{
		ID:   1,
		Name: "Today's Holiday",
		Date: today,
		Type: models.NationalHoliday,
	}

	mockRepo.On("GetByDate", mock.MatchedBy(func(date time.Time) bool {
		return date.Format("2006-01-02") == today.Format("2006-01-02")
	})).Return(expectedHoliday, nil)

	holiday, err := service.GetHolidayToday()

	assert.NoError(t, err)
	assert.NotNil(t, holiday)
	assert.Equal(t, expectedHoliday.Name, holiday.Name)
	mockRepo.AssertExpectations(t)
}

func TestHolidayService_GetHolidaysByYear(t *testing.T) {
	mockRepo := new(MockHolidayRepository)
	service := NewHolidayService(mockRepo)

	year := 2024
	expectedHolidays := []models.Holiday{
		{ID: 1, Name: "New Year", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "Christmas", Date: time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)},
	}

	mockRepo.On("GetAll", mock.MatchedBy(func(filter models.HolidayFilter) bool {
		return filter.Year != nil && *filter.Year == year
	})).Return(expectedHolidays, len(expectedHolidays), nil)

	holidays, err := service.GetHolidaysByYear(year)

	assert.NoError(t, err)
	assert.Len(t, holidays, 2)
	assert.Equal(t, expectedHolidays[0].Name, holidays[0].Name)
	mockRepo.AssertExpectations(t)
}
