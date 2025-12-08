package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookingTrendsRepository struct {
	mock.Mock
}

func (m *MockBookingTrendsRepository) GetTrends(year int) ([]model.BookingTrend, error) {
	args := m.Called(year)
	return args.Get(0).([]model.BookingTrend), args.Error(1)
}

func TestBookingTrendsService_GetTrends(t *testing.T) {
	mockRepo := new(MockBookingTrendsRepository)
	svc := service.NewBookingTrendsService(mockRepo)

	expected := []model.BookingTrend{
		{Year: 2023, Month: "1", Booking: 100},
		{Year: 2023, Month: "2", Booking: 150},
	}
	mockRepo.On("GetTrends", 2023).Return(expected, nil)

	result, err := svc.GetTrends(2023)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestBookingTrendsService_GetTrends_Error(t *testing.T) {
	mockRepo := new(MockBookingTrendsRepository)
	svc := service.NewBookingTrendsService(mockRepo)

	mockRepo.On("GetTrends", 2023).Return([]model.BookingTrend{}, assert.AnError)

	result, err := svc.GetTrends(2023)
	assert.Error(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}
