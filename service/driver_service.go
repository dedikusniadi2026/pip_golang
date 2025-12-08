package service

import (
	"auth-service/model"
	"auth-service/repository"
)

type DriverService struct {
	Repo repository.DriverRepositoryInterface
}

func (s *DriverService) GetAll() ([]model.Driver, error) {
	return s.Repo.GetAll()
}

func (s *DriverService) Create(driver *model.Driver) error {
	return s.Repo.Create(driver)
}

func (s *DriverService) GetByID(id string) (*model.Driver, error) {
	return s.Repo.GetByID(id)
}

func (s *DriverService) Update(id string, driver *model.Driver) error {
	return s.Repo.Update(id, driver)
}

func (s *DriverService) Delete(id string) error {
	return s.Repo.Delete(id)
}
