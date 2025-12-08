package service

import (
	"auth-service/model"
	"auth-service/repository"
)

type BookingTrendsService struct {
	Repo repository.BookingTrendsRepositoryInterface
}

func NewBookingTrendsService(repo repository.BookingTrendsRepositoryInterface) *BookingTrendsService {
	return &BookingTrendsService{Repo: repo}
}

func (s *BookingTrendsService) GetTrends(year int) ([]model.BookingTrend, error) {
	return s.Repo.GetTrends(year)
}
