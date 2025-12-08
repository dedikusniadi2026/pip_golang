package service

import (
	"auth-service/model"
	"auth-service/repository"
)

type TripHistoryServiceInterface interface {
	Create(*model.DriverAssignment) error
	Update(*model.DriverAssignment) error
	Delete(uint) error
	FindByVehicle(uint) ([]model.DriverAssignment, error)
}

type TripHistoryMockService interface {
	GetTripHistory() ([]model.TripHistory, error)
}

type MockTripHistoryService struct {
	ReturnError bool
}

func (m *MockTripHistoryService) GetTripHistory() ([]model.TripHistory, error) {
	return []model.TripHistory{}, nil
}

type TripHistoryService struct {
	Repo repository.TripHistoryRepositoryInterface
}

func NewTripHistoryService(repo repository.TripHistoryRepositoryInterface) *TripHistoryService {
	return &TripHistoryService{Repo: repo}
}

func (s *TripHistoryService) GetTripHistory() ([]model.TripHistory, error) {
	return s.Repo.GetTripHistory()
}
