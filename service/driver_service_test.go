package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDriverRepository struct {
	mock.Mock
}

func (m *MockDriverRepository) GetAll() ([]model.Driver, error) {
	args := m.Called()
	return args.Get(0).([]model.Driver), args.Error(1)
}

func (m *MockDriverRepository) Create(d *model.Driver) error {
	args := m.Called(d)
	return args.Error(0)
}

func (m *MockDriverRepository) GetByID(id string) (*model.Driver, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Driver), args.Error(1)
}

func (m *MockDriverRepository) Update(id string, d *model.Driver) error {
	args := m.Called(id, d)
	return args.Error(0)
}

func (m *MockDriverRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestDriverService_GetAll(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	svc := &service.DriverService{Repo: mockRepo}

	expected := []model.Driver{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
	}
	mockRepo.On("GetAll").Return(expected, nil)

	result, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestDriverService_Create(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	svc := &service.DriverService{Repo: mockRepo}

	driver := &model.Driver{Name: "John Doe", Email: "john@example.com"}
	mockRepo.On("Create", driver).Return(nil)

	err := svc.Create(driver)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDriverService_GetByID(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	svc := &service.DriverService{Repo: mockRepo}

	expected := &model.Driver{ID: 1, Name: "John Doe", Email: "john@example.com"}
	mockRepo.On("GetByID", "1").Return(expected, nil)

	result, err := svc.GetByID("1")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestDriverService_Update(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	svc := &service.DriverService{Repo: mockRepo}

	driver := &model.Driver{ID: 1, Name: "Updated Name", Email: "updated@example.com"}
	mockRepo.On("Update", "1", driver).Return(nil)

	err := svc.Update("1", driver)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDriverService_Delete(t *testing.T) {
	mockRepo := new(MockDriverRepository)
	svc := &service.DriverService{Repo: mockRepo}

	mockRepo.On("Delete", "1").Return(nil)

	err := svc.Delete("1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
