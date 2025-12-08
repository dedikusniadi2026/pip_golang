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

func TestTripHistoryService_GetTripHistory(t *testing.T) {
	mockRepo := new(MockTripHistoryRepository)
	service := service.TripHistoryService{
		Repo: mockRepo,
	}

	expectedData := []model.TripHistory{
		{ID: 1, CustomerName: "Trip A"},
	}

	mockRepo.On("GetTripHistory").Return(expectedData, nil)

	result, err := service.GetTripHistory()

	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
	mockRepo.AssertExpectations(t)
}

func TestTripHistoryService_GetTripHistory_Error(t *testing.T) {
	mockRepo := new(MockTripHistoryRepository)
	svc := service.NewTripHistoryService(mockRepo)

	expectedError := errors.New("database error")
	mockRepo.On("GetTripHistory").Return(([]model.TripHistory)(nil), expectedError)

	result, err := svc.GetTripHistory()
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMockTripHistoryRepository(t *testing.T) {
	mockRepo := &MockTripHistoryRepository{}
	mockRepo.On("GetTripHistory").Return([]model.TripHistory{}, nil)
	mockRepo.GetTripHistory()
}
