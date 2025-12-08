package service

import (
	"auth-service/model"
	"auth-service/repository"
)

type CarModelService struct {
	repo repository.CarModelRepositoryInterface
}

func NewCarModelService(repo repository.CarModelRepositoryInterface) *CarModelService {
	return &CarModelService{repo: repo}
}

func (s *CarModelService) GetAll() ([]model.CarModel, error) {
	return s.repo.FindAll()
}

func (s *CarModelService) GetByID(id int) (*model.CarModel, error) {
	return s.repo.GetByID(id)
}

func (s *CarModelService) Create(cm model.CarModel) error {
	return s.repo.Create(&cm)
}
