package repository

import (
	"auth-service/model"
	"database/sql"
)

type BookingTrendsRepositoryInterface interface {
	GetTrends(year int) ([]model.BookingTrend, error)
}

type BookingTrendsRepository struct {
	DB *sql.DB
}

func NewBookingTrendsRepository(db *sql.DB) *BookingTrendsRepository {
	return &BookingTrendsRepository{DB: db}
}

func (r *BookingTrendsRepository) GetTrends(year int) ([]model.BookingTrend, error) {
	rows, err := r.DB.Query(`
        SELECT month, booking_count, year 
        FROM booking_trends 
        WHERE year = $1
        ORDER BY id ASC
    `, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []model.BookingTrend
	for rows.Next() {
		var t model.BookingTrend
		if err := rows.Scan(&t.Month, &t.Booking, &t.Year); err != nil {
			return nil, err
		}
		trends = append(trends, t)
	}

	return trends, nil
}
