package repository

import (
	"auth-service/model"
	"database/sql"
)

type PDFRepositoryInterface interface {
	GetTripByID(tripID string) (model.Pdf, error)
}

type PDFRepository struct {
	DB *sql.DB
}

func NewPDFRepository(db *sql.DB) *PDFRepository {
	return &PDFRepository{DB: db}
}

func (r *PDFRepository) GetTripByID(tripID string) (model.Pdf, error) {
	var trip model.Pdf

	query := `
        SELECT
            id, customer_name, booking_date, duration_minutes, distance_km,
            pickup_location, destination, driver_name, vehicle_name,
            amount, rating, feedback
        FROM trips
        WHERE id = $1
    `

	row := r.DB.QueryRow(query, tripID)

	err := row.Scan(
		&trip.ID,
		&trip.CustomerName,
		&trip.BookingDate,
		&trip.DurationMinutes,
		&trip.DistanceKM,
		&trip.PickupLocation,
		&trip.Destination,
		&trip.DriverName,
		&trip.VehicleName,
		&trip.Amount,
		&trip.Rating,
		&trip.Feedback,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Pdf{}, nil
		}
		return model.Pdf{}, err
	}

	return trip, nil
}
