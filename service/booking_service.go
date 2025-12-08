package service

import (
	"auth-service/model"
	"auth-service/repository"
)

type BookingServiceInterface interface {
	Create(*model.Booking) error
	GetAll() ([]model.Booking, error)
	Update(*model.Booking) error
	Delete(id string) error
}

type BookingService struct {
	Repo repository.BookingRepositoryInterface
}

func (s *BookingService) Create(b *model.Booking) error {
	return s.Repo.Create(b)
}

func (s *BookingService) GetAll() ([]model.Booking, error) {
	return s.Repo.GetAll()
}

func (s *BookingService) Update(b *model.Booking) error {
	return s.Repo.Update(b)
}

func (s *BookingService) Delete(id string) error {
	return s.Repo.Delete(id)
}
