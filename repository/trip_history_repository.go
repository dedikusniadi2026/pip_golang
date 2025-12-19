package repository

import (
	"auth-service/model"
	"database/sql"
)

type TripHistoryRepositoryInterface interface {
	GetTripHistory() ([]model.TripHistory, error)
}

type TripHistoryRepository struct {
	DB *sql.DB
}

func NewTripHistoryRepository(db *sql.DB) *TripHistoryRepository {
	return &TripHistoryRepository{DB: db}
}

func (r *TripHistoryRepository) GetTripHistory() ([]model.TripHistory, error) {
	query := `
        SELECT id, booking_code, customer_name, booking_date,
               duration_minutes, distance_km, pickup_location, destination,
               driver_name, vehicle_name, amount, rating, feedback
        FROM trips
    `
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trips := []model.TripHistory{}

	for rows.Next() {
		var b model.TripHistory
		err := rows.Scan(
			&b.ID, &b.BookingCode, &b.CustomerName, &b.BookingDate,
			&b.DurationMinutes, &b.DistanceKM, &b.PickupLocation,
			&b.Destination, &b.DriverName, &b.VehicleName,
			&b.Amount, &b.Rating, &b.Feedback,
		)
		if err != nil {
			return nil, err
		}
		trips = append(trips, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trips, nil
}
