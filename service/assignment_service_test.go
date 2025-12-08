package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAssignmentRepo struct {
	mock.Mock
}

func (m *MockAssignmentRepo) Create(a *model.DriverAssignment) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m *MockAssignmentRepo) FindByVehicle(vehicleID uint) ([]model.DriverAssignment, error) {
	args := m.Called(vehicleID)
	return args.Get(0).([]model.DriverAssignment), args.Error(1)
}

func (m *MockAssignmentRepo) Update(a *model.DriverAssignment) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m *MockAssignmentRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAssignmentsService_Create(t *testing.T) {
	mockRepo := new(MockAssignmentRepo)
	svc := service.NewAssignmentsService(mockRepo)

	assignment := &model.DriverAssignment{ID: 1, VehicleID: 1, DriverID: sql.NullInt64{Int64: 1, Valid: true}}
	mockRepo.On("Create", assignment).Return(nil)

	err := svc.Create(assignment)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAssignmentsService_Update(t *testing.T) {
	mockRepo := new(MockAssignmentRepo)
	svc := service.NewAssignmentsService(mockRepo)
	assignment := &model.DriverAssignment{ID: 1, VehicleID: 1, DriverID: sql.NullInt64{Int64: 1, Valid: true}}
	mockRepo.On("Update", assignment).Return(nil)
	err := svc.Update(assignment)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAssignmentsService_FindByVehicle(t *testing.T) {
	mockRepo := new(MockAssignmentRepo)
	svc := service.NewAssignmentsService(mockRepo)
	expected := []model.DriverAssignment{
		{ID: 1, VehicleID: 1, DriverID: sql.NullInt64{Int64: 1, Valid: true}},
		{ID: 2, VehicleID: 1, DriverID: sql.NullInt64{Int64: 2, Valid: true}},
	}

	mockRepo.On("FindByVehicle", uint(1)).Return(expected, nil)

	result, err := svc.FindByVehicle(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestAssignmentsService_Delete(t *testing.T) {
	mockRepo := new(MockAssignmentRepo)
	svc := service.NewAssignmentsService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
