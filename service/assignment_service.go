package service

import (
	"auth-service/model"
	"auth-service/repository"
)

type AssignmentsServiceInterface interface {
	Create(*model.DriverAssignment) error
	Update(*model.DriverAssignment) error
	Delete(uint) error
	FindByVehicle(uint) ([]model.DriverAssignment, error)
}

type AssignmentsService struct {
	Repo repository.AssignmentsRepositoryInterface
}

func NewAssignmentsService(repo repository.AssignmentsRepositoryInterface) *AssignmentsService {
	return &AssignmentsService{Repo: repo}
}

func (s *AssignmentsService) Create(a *model.DriverAssignment) error {
	return s.Repo.Create(a)
}

func (s *AssignmentsService) FindByVehicle(vehicleID uint) ([]model.DriverAssignment, error) {
	return s.Repo.FindByVehicle(vehicleID)
}

func (s *AssignmentsService) Update(a *model.DriverAssignment) error {
	return s.Repo.Update(a)
}

func (s *AssignmentsService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
