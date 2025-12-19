package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTripHistoryRepository struct {
	mock.Mock
}

func (m *MockTripHistoryRepository) GetTripHistory() ([]model.TripHistory, error) {
	args := m.Called()
	return args.Get(0).([]model.TripHistory), args.Error(1)
}

func (m *MockTripHistoryRepository) GetBookingByCode() (*model.TripHistory, error) {
	args := m.Called()
	return args.Get(0).(*model.TripHistory), args.Error(1)
}

func TestMockTripHistoryService_GetTripHistory(t *testing.T) {
	mockService := &service.MockTripHistoryService{}

	result, err := mockService.GetTripHistory()

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestTripHistoryService_GetTripHistory(t *testing.T) {
	mockRepo := new(MockTripHistoryRepository)
	svc := service.NewTripHistoryService(mockRepo)

	expected := []model.TripHistory{{ID: 1, CustomerName: "Trip A"}}
	mockRepo.On("GetTripHistory").Return(expected, nil)

	result, err := svc.GetTripHistory()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestTripHistoryService_GetTripHistory_Error(t *testing.T) {
	mockRepo := new(MockTripHistoryRepository)
	svc := service.NewTripHistoryService(mockRepo)

	expectedErr := errors.New("db error")
	mockRepo.On("GetTripHistory").Return([]model.TripHistory(nil), expectedErr)

	result, err := svc.GetTripHistory()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

func TestTripHistoryService_GetTripHistory_Empty(t *testing.T) {
	mockRepo := new(MockTripHistoryRepository)
	svc := service.NewTripHistoryService(mockRepo)

	mockRepo.On("GetTripHistory").Return(([]model.TripHistory)(nil), nil)

	result, err := svc.GetTripHistory()
	assert.NoError(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)

}
