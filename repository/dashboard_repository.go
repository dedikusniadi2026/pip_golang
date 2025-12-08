package repository

import "database/sql"

type DashboardRepositoryInterface interface {
	GetTotalBookings() (int, error)
	GetActiveDrivers() (int, error)
	GetTotalRevenue() (float64, error)
}

type DashboardRepository struct {
	DB *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{DB: db}
}

func (r *DashboardRepository) GetTotalBookings() (int, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(*) FROM booking`).Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetActiveDrivers() (int, error) {
	var count int
	err := r.DB.QueryRow(`SELECT COUNT(*) FROM drivers WHERE status = 'active'`).Scan(&count)
	return count, err
}

func (r *DashboardRepository) GetTotalRevenue() (float64, error) {
	var total float64
	err := r.DB.QueryRow(`
    SELECT COALESCE(SUM(amount), 0)::FLOAT  FROM payment`).Scan(&total)
	return total, err
}
