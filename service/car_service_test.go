package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCarRepository struct {
	mock.Mock
}

func (m *MockCarRepository) GetAll() ([]model.Car, error) {
	args := m.Called()
	return args.Get(0).([]model.Car), args.Error(1)
}

func (m *MockCarRepository) GetByID(id int) (*model.Car, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Car), args.Error(1)
}

func (m *MockCarRepository) Create(v model.Car) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockCarRepository) Update(id int, v model.Car) error {
	args := m.Called(id, v)
	return args.Error(0)
}

func (m *MockCarRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCarService_GetAll(t *testing.T) {
	mockRepo := new(MockCarRepository)
	svc := service.NewCarService(mockRepo)

	expected := []model.Car{
		{ID: 1, Model: "Model A"},
		{ID: 2, Model: "Model B"},
	}
	mockRepo.On("GetAll").Return(expected, nil)

	result, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCarService_GetByID(t *testing.T) {
	mockRepo := new(MockCarRepository)
	svc := service.NewCarService(mockRepo)

	expected := &model.Car{ID: 1, Model: "Model A"}
	mockRepo.On("GetByID", 1).Return(expected, nil)

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCarService_Create(t *testing.T) {
	mockRepo := new(MockCarRepository)
	svc := service.NewCarService(mockRepo)

	car := model.Car{Model: "Model A"}
	mockRepo.On("Create", car).Return(nil)

	err := svc.Create(car)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCarService_Update(t *testing.T) {
	mockRepo := new(MockCarRepository)
	svc := service.NewCarService(mockRepo)

	car := model.Car{ID: 1, Model: "Updated Model"}
	mockRepo.On("Update", 1, car).Return(nil)

	err := svc.Update(1, car)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCarService_Delete(t *testing.T) {
	mockRepo := new(MockCarRepository)
	svc := service.NewCarService(mockRepo)

	mockRepo.On("Delete", 1).Return(nil)

	err := svc.Delete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
