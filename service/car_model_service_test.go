package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCarModelRepository struct {
	mock.Mock
}

func (m *MockCarModelRepository) FindAll() ([]model.CarModel, error) {
	args := m.Called()
	return args.Get(0).([]model.CarModel), args.Error(1)
}

func (m *MockCarModelRepository) GetByID(id int) (*model.CarModel, error) {
	args := m.Called(id)
	return args.Get(0).(*model.CarModel), args.Error(1)
}

func (m *MockCarModelRepository) Create(cm *model.CarModel) error {
	args := m.Called(cm)
	return args.Error(0)
}

func TestCarModelService_GetAll(t *testing.T) {
	mockRepo := new(MockCarModelRepository)
	svc := service.NewCarModelService(mockRepo)

	expected := []model.CarModel{
		{ID: 1, ModelName: "Model A", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, ModelName: "Model B", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	mockRepo.On("FindAll").Return(expected, nil)

	result, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCarModelService_GetAll_Error(t *testing.T) {
	mockRepo := new(MockCarModelRepository)
	svc := service.NewCarModelService(mockRepo)

	mockRepo.On("FindAll").Return([]model.CarModel{}, assert.AnError)

	result, err := svc.GetAll()
	assert.Error(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCarModelService_GetByID(t *testing.T) {
	mockRepo := new(MockCarModelRepository)
	svc := service.NewCarModelService(mockRepo)

	expected := &model.CarModel{ID: 1, ModelName: "Model A", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockRepo.On("GetByID", 1).Return(expected, nil)

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCarModelService_GetByID_Error(t *testing.T) {
	mockRepo := new(MockCarModelRepository)
	svc := service.NewCarModelService(mockRepo)

	mockRepo.On("GetByID", 1).Return((*model.CarModel)(nil), assert.AnError)

	result, err := svc.GetByID(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCarModelService_Create(t *testing.T) {
	mockRepo := new(MockCarModelRepository)
	svc := service.NewCarModelService(mockRepo)

	cm := model.CarModel{ModelName: "Model A", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockRepo.On("Create", &cm).Return(nil)

	err := svc.Create(cm)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCarModelService_Create_Error(t *testing.T) {
	mockRepo := new(MockCarModelRepository)
	svc := service.NewCarModelService(mockRepo)

	cm := model.CarModel{ModelName: "Model A", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockRepo.On("Create", &cm).Return(assert.AnError)

	err := svc.Create(cm)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
