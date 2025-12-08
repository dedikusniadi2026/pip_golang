package repository

import (
	"auth-service/model"
	"context"
	"database/sql"
)

type DashboardTripRepositoryInterface interface {
	GetDashboardSummary(ctx context.Context) (*model.DashboardSummary, error)
}

type DashboardTripRepository struct {
	DB *sql.DB
}

func NewDashboardTripRepository(db *sql.DB) *DashboardTripRepository {
	return &DashboardTripRepository{DB: db}
}

func (r *DashboardTripRepository) GetDashboardSummary(ctx context.Context) (*model.DashboardSummary, error) {

	summary := &model.DashboardSummary{}

	err := r.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM vehicle_trips
	`).Scan(&summary.TotalTrips)
	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(distance_km), 0) FROM vehicle_trips
	`).Scan(&summary.TotalDistance)
	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(duration), 0) FROM vehicle_trips
	`).Scan(&summary.TotalDuration)
	if err != nil {
		return nil, err
	}

	return summary, nil
}
