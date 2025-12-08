package service

import (
	"auth-service/model"
	"auth-service/repository"
	"context"
)

type DashboardTripServiceInterface interface {
	GetDashboardSummary(ctx context.Context) (*model.DashboardSummary, error)
}

type DashboardTripService struct {
	Repo repository.DashboardTripRepositoryInterface
}

func NewDashboardTripService(repo repository.DashboardTripRepositoryInterface) *DashboardTripService {
	return &DashboardTripService{Repo: repo}
}

func (s *DashboardTripService) GetDashboardSummary(ctx context.Context) (*model.DashboardSummary, error) {
	return s.Repo.GetDashboardSummary(ctx)
}
