package service_test

import (
	"auth-service/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) GetTotalBookings() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockDashboardRepository) GetActiveDrivers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockDashboardRepository) GetTotalRevenue() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func TestDashboardService_GetDashboardData(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	svc := service.NewDashboardService(mockRepo)

	mockRepo.On("GetTotalBookings").Return(100, nil)
	mockRepo.On("GetActiveDrivers").Return(50, nil)
	mockRepo.On("GetTotalRevenue").Return(1000000.0, nil)

	result, err := svc.GetDashboardData()
	assert.NoError(t, err)
	assert.Equal(t, 100, result["totalBookings"])
	assert.Equal(t, 50, result["activeDrivers"])
	assert.Equal(t, 1000000.0, result["totalRevenue"])
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetDashboardData_GetTotalBookingsError(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	svc := service.NewDashboardService(mockRepo)

	mockRepo.On("GetTotalBookings").Return(0, errors.New("test error"))

	result, err := svc.GetDashboardData()
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetDashboardData_GetActiveDriversError(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	svc := service.NewDashboardService(mockRepo)

	mockRepo.On("GetTotalBookings").Return(100, nil)
	mockRepo.On("GetActiveDrivers").Return(0, errors.New("test error"))

	result, err := svc.GetDashboardData()
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetDashboardData_GetTotalRevenueError(t *testing.T) {
	mockRepo := new(MockDashboardRepository)
	svc := service.NewDashboardService(mockRepo)

	mockRepo.On("GetTotalBookings").Return(100, nil)
	mockRepo.On("GetActiveDrivers").Return(50, nil)
	mockRepo.On("GetTotalRevenue").Return(0.0, errors.New("test error"))

	result, err := svc.GetDashboardData()
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
