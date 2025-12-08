package service

import (
	"auth-service/model"
	"auth-service/repository"
	"context"
)

type TripServiceInterface interface {
	Create(*model.VehicleTrip) error
	FindByVehicle(uint) ([]model.VehicleTrip, error)
	Update(*model.VehicleTrip) error
	Delete(uint) error
	GetTripTotals(context.Context) (*model.TotalTrips, error)
}

type TripService struct {
	repo repository.TripsRepositoryInterface
}

func NewTripService(repo repository.TripsRepositoryInterface) *TripService {
	return &TripService{repo: repo}
}

func (s *TripService) Create(t *model.VehicleTrip) error {
	return s.repo.Create(t)
}

func (s *TripService) Update(t *model.VehicleTrip) error {
	return s.repo.Update(t)
}

func (s *TripService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *TripService) FindByVehicle(vehicleID uint) ([]model.VehicleTrip, error) {
	return s.repo.FindByVehicle(vehicleID)
}

func (s *TripService) GetTripTotals(ctx context.Context) (*model.TotalTrips, error) {
	return s.repo.GetTripTotal(ctx)
}
