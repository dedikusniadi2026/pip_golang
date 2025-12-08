package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) Create(b *model.Booking) error {
	args := m.Called(b)
	return args.Error(0)
}

func (m *MockBookingRepository) GetAll() ([]model.Booking, error) {
	args := m.Called()
	return args.Get(0).([]model.Booking), args.Error(1)
}

func (m *MockBookingRepository) Update(b *model.Booking) error {
	args := m.Called(b)
	return args.Error(0)
}

func (m *MockBookingRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBookingService_Create(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	svc := &service.BookingService{Repo: mockRepo}

	booking := &model.Booking{ID: "1"}

	mockRepo.On("Create", booking).Return(nil)

	err := svc.Create(booking)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestBookingService_GetAll(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	svc := &service.BookingService{Repo: mockRepo}

	expected := []model.Booking{
		{ID: "1", Customer: "Booking 1"},
	}

	mockRepo.On("GetAll").Return(expected, nil)

	result, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestBookingService_Update(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	svc := &service.BookingService{Repo: mockRepo}

	booking := &model.Booking{ID: "1"}

	mockRepo.On("Update", booking).Return(nil)

	err := svc.Update(booking)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestBookingService_Delete(t *testing.T) {
	mockRepo := new(MockBookingRepository)
	svc := &service.BookingService{Repo: mockRepo}

	id := "1"

	mockRepo.On("Delete", id).Return(nil)

	err := svc.Delete(id)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
