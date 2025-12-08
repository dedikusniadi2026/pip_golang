package service

import (
	"auth-service/model"
	"auth-service/repository"
)

type PopularDestinationService struct {
	Repo repository.PopularDestinationRepositoryInterface
}

type PopularDestinationServiceInterface interface {
	GetAll() ([]model.PopularDestination, error)
}

func (s *PopularDestinationService) GetAll() ([]model.PopularDestination, error) {
	return s.Repo.GetAll()
}

func (s *PopularDestinationService) Add(destination string, bookings int) (*model.PopularDestination, error) {
	pd := model.PopularDestination{
		Destination: destination,
		Bookings:    bookings,
	}
	return s.Repo.Add(pd)
}

func (s *PopularDestinationService) UpdateBookings(id int, bookings int) error {
	return s.Repo.UpdateBookings(id, bookings)
}

func (s *PopularDestinationService) Delete(id int) error {
	return s.Repo.Delete(id)
}
