package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMaintenanceRepository struct {
	mock.Mock
}

func (m *MockMaintenanceRepository) Create(maint *model.VehicleMaintenance) error {
	args := m.Called(maint)
	return args.Error(0)
}

func (m *MockMaintenanceRepository) FindByVehicle(vehicleID int) ([]model.VehicleMaintenance, error) {
	args := m.Called(vehicleID)
	return args.Get(0).([]model.VehicleMaintenance), args.Error(1)
}

func (m *MockMaintenanceRepository) Update(maint *model.VehicleMaintenance) error {
	args := m.Called(maint)
	return args.Error(0)
}

func (m *MockMaintenanceRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestMaintenanceService_Create(t *testing.T) {
	mockRepo := new(MockMaintenanceRepository)
	svc := service.NewMaintenanceService(mockRepo)

	maint := &model.VehicleMaintenance{
		VehicleID:   1,
		ServiceType: "Oil Change",
		Cost:        50000,
	}
	mockRepo.On("Create", maint).Return(nil)

	err := svc.Create(maint)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_FindByVehicle(t *testing.T) {
	mockRepo := new(MockMaintenanceRepository)
	svc := service.NewMaintenanceService(mockRepo)

	expected := []model.VehicleMaintenance{
		{ID: 1, VehicleID: 1, ServiceType: "Oil Change", Cost: 50000},
		{ID: 2, VehicleID: 1, ServiceType: "Tire Replacement", Cost: 200000},
	}
	mockRepo.On("FindByVehicle", 1).Return(expected, nil)

	result, err := svc.FindByVehicle(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_Update(t *testing.T) {
	mockRepo := new(MockMaintenanceRepository)
	svc := service.NewMaintenanceService(mockRepo)

	maint := &model.VehicleMaintenance{
		ID:          1,
		VehicleID:   1,
		ServiceType: "Updated Maintenance",
		Cost:        75000,
	}
	mockRepo.On("Update", maint).Return(nil)

	err := svc.Update(maint)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMaintenanceService_Delete(t *testing.T) {
	mockRepo := new(MockMaintenanceRepository)
	svc := service.NewMaintenanceService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
