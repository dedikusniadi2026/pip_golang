package service

import (
	"auth-service/model"
	"auth-service/repository"
)

type CarTypeService struct {
	Repo repository.CarTypeRepositoryInterface
}

func NewCarTypeService(repo repository.CarTypeRepositoryInterface) *CarTypeService {
	return &CarTypeService{Repo: repo}
}

func (s *CarTypeService) GetAll() ([]model.CarType, error) {
	return s.Repo.FindAll()
}

func (s *CarTypeService) GetByID(id int) (*model.CarType, error) {
	return s.Repo.GetByID(id)
}

func (s *CarTypeService) Create(ct model.CarType) error {
	return s.Repo.Create(ct)
}
