package service

import "auth-service/repository"

type DashboardServiceInterface interface {
	GetDashboardData() (map[string]interface{}, error)
}

type DashboardService struct {
	repo repository.DashboardRepositoryInterface
}

func NewDashboardService(repo repository.DashboardRepositoryInterface) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetDashboardData() (map[string]interface{}, error) {
	bookings, err := s.repo.GetTotalBookings()
	if err != nil {
		return nil, err
	}

	drivers, err := s.repo.GetActiveDrivers()
	if err != nil {
		return nil, err
	}

	revenue, err := s.repo.GetTotalRevenue()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"totalBookings": bookings,
		"activeDrivers": drivers,
		"totalRevenue":  revenue,
	}, nil
}
