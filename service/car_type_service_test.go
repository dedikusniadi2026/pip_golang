package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCarTypeRepository struct {
	mock.Mock
}

func (m *MockCarTypeRepository) FindAll() ([]model.CarType, error) {
	args := m.Called()
	return args.Get(0).([]model.CarType), args.Error(1)
}

func (m *MockCarTypeRepository) GetByID(id int) (*model.CarType, error) {
	args := m.Called(id)
	return args.Get(0).(*model.CarType), args.Error(1)
}

func (m *MockCarTypeRepository) Create(ct model.CarType) error {
	args := m.Called(ct)
	return args.Error(0)
}

func TestCarTypeService_GetAll(t *testing.T) {
	mockRepo := new(MockCarTypeRepository)
	svc := service.NewCarTypeService(mockRepo)

	expected := []model.CarType{
		{ID: "1", TypeName: "Sedan"},
		{ID: "2", TypeName: "SUV"},
	}
	mockRepo.On("FindAll").Return(expected, nil)

	result, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCarTypeService_GetByID(t *testing.T) {
	mockRepo := new(MockCarTypeRepository)
	svc := service.NewCarTypeService(mockRepo)

	expected := &model.CarType{ID: "1", TypeName: "Sedan"}
	mockRepo.On("GetByID", 1).Return(expected, nil)

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCarTypeService_Create(t *testing.T) {
	mockRepo := new(MockCarTypeRepository)
	svc := service.NewCarTypeService(mockRepo)

	ct := model.CarType{TypeName: "Sedan"}
	mockRepo.On("Create", ct).Return(nil)

	err := svc.Create(ct)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
