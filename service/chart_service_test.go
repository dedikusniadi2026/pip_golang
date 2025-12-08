package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPopularDestinationRepository struct {
	mock.Mock
}

func (m *MockPopularDestinationRepository) GetAll() ([]model.PopularDestination, error) {
	args := m.Called()
	return args.Get(0).([]model.PopularDestination), args.Error(1)
}

func (m *MockPopularDestinationRepository) Add(pd model.PopularDestination) (*model.PopularDestination, error) {
	args := m.Called(pd)
	return args.Get(0).(*model.PopularDestination), args.Error(1)
}

func (m *MockPopularDestinationRepository) UpdateBookings(id int, bookings int) error {
	args := m.Called(id, bookings)
	return args.Error(0)
}

func (m *MockPopularDestinationRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestPopularDestinationService_GetAll(t *testing.T) {
	mockRepo := new(MockPopularDestinationRepository)
	svc := &service.PopularDestinationService{Repo: mockRepo}

	expected := []model.PopularDestination{
		{ID: 1, Destination: "Jakarta", Bookings: 100, CreatedAt: time.Now()},
		{ID: 2, Destination: "Bandung", Bookings: 80, CreatedAt: time.Now()},
	}
	mockRepo.On("GetAll").Return(expected, nil)

	result, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestPopularDestinationService_GetAll_Error(t *testing.T) {
	mockRepo := new(MockPopularDestinationRepository)
	svc := &service.PopularDestinationService{Repo: mockRepo}

	mockRepo.On("GetAll").Return([]model.PopularDestination{}, assert.AnError)

	result, err := svc.GetAll()
	assert.Error(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestPopularDestinationService_Add(t *testing.T) {
	mockRepo := new(MockPopularDestinationRepository)
	svc := &service.PopularDestinationService{Repo: mockRepo}

	input := model.PopularDestination{Destination: "Jakarta", Bookings: 100}
	expected := &model.PopularDestination{ID: 1, Destination: "Jakarta", Bookings: 100, CreatedAt: time.Now()}
	mockRepo.On("Add", input).Return(expected, nil)

	result, err := svc.Add("Jakarta", 100)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestPopularDestinationService_Add_Error(t *testing.T) {
	mockRepo := new(MockPopularDestinationRepository)
	svc := &service.PopularDestinationService{Repo: mockRepo}

	input := model.PopularDestination{Destination: "Jakarta", Bookings: 100}
	mockRepo.On("Add", input).Return((*model.PopularDestination)(nil), assert.AnError)

	result, err := svc.Add("Jakarta", 100)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestPopularDestinationService_UpdateBookings(t *testing.T) {
	mockRepo := new(MockPopularDestinationRepository)
	svc := &service.PopularDestinationService{Repo: mockRepo}

	mockRepo.On("UpdateBookings", 1, 150).Return(nil)

	err := svc.UpdateBookings(1, 150)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPopularDestinationService_UpdateBookings_Error(t *testing.T) {
	mockRepo := new(MockPopularDestinationRepository)
	svc := &service.PopularDestinationService{Repo: mockRepo}

	mockRepo.On("UpdateBookings", 1, 150).Return(assert.AnError)

	err := svc.UpdateBookings(1, 150)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPopularDestinationService_Delete(t *testing.T) {
	mockRepo := new(MockPopularDestinationRepository)
	svc := &service.PopularDestinationService{Repo: mockRepo}

	mockRepo.On("Delete", 1).Return(nil)

	err := svc.Delete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPopularDestinationService_Delete_Error(t *testing.T) {
	mockRepo := new(MockPopularDestinationRepository)
	svc := &service.PopularDestinationService{Repo: mockRepo}

	mockRepo.On("Delete", 1).Return(assert.AnError)

	err := svc.Delete(1)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
