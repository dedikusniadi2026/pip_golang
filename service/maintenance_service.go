package service

import (
	"auth-service/model"
	"auth-service/repository"
)

type MaintenanceServiceInterface interface {
	Create(*model.VehicleMaintenance) error
	FindByVehicle(vehicleID int) ([]model.VehicleMaintenance, error)
	Update(*model.VehicleMaintenance) error
	Delete(id uint) error
}

type MaintenanceService struct {
	Repo repository.MaintenanceRepositoryInterface
}

func NewMaintenanceService(repo repository.MaintenanceRepositoryInterface) *MaintenanceService {
	return &MaintenanceService{Repo: repo}
}

func (s *MaintenanceService) Create(m *model.VehicleMaintenance) error {
	return s.Repo.Create(m)
}

func (s *MaintenanceService) FindByVehicle(vehicleID int) ([]model.VehicleMaintenance, error) {
	return s.Repo.FindByVehicle(vehicleID)
}

func (s *MaintenanceService) Update(m *model.VehicleMaintenance) error {
	return s.Repo.Update(m)
}

func (s *MaintenanceService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
