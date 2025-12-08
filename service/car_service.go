package service

import (
	"auth-service/model"
	"auth-service/repository"
)

type CarService struct {
	Repo repository.CarRepositoryInterface
}

func NewCarService(repo repository.CarRepositoryInterface) *CarService {
	return &CarService{Repo: repo}
}

func (s *CarService) GetAll() ([]model.Car, error) {
	return s.Repo.GetAll()
}

func (s *CarService) GetByID(id int) (*model.Car, error) {
	return s.Repo.GetByID(id)
}

func (s *CarService) Create(v model.Car) error {
	return s.Repo.Create(v)
}

func (s *CarService) Update(id int, v model.Car) error {
	return s.Repo.Update(id, v)
}

func (s *CarService) Delete(id int) error {
	return s.Repo.Delete(id)
}
