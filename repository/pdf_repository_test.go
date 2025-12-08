package repository_test

import (
	"auth-service/model"
	"auth-service/repository"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPDFRepository_GetTripByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPDFRepository(db)

	tripID := "123"
	bookingDate := time.Now()
	rows := sqlmock.NewRows([]string{"id", "customer_name", "booking_date", "duration_minutes", "distance_km", "pickup_location", "destination", "driver_name", "vehicle_name", "amount", "rating", "feedback"}).
		AddRow("123", "John Doe", bookingDate, 60, 50, "Location A", "Location B", "Driver X", "Car Y", 100, 4.5, "Good trip")

	mock.ExpectQuery(`SELECT id, customer_name, booking_date, duration_minutes, distance_km, pickup_location, destination, driver_name, vehicle_name, amount, rating, feedback FROM trips WHERE id = \$1`).
		WithArgs(tripID).
		WillReturnRows(rows)

	trip, err := repo.GetTripByID(tripID)

	assert.NoError(t, err)
	assert.NotNil(t, trip)
	assert.Equal(t, "123", trip.ID)
	assert.Equal(t, "John Doe", trip.CustomerName)
	assert.Equal(t, bookingDate, trip.BookingDate)
	assert.Equal(t, 60, trip.DurationMinutes)
	assert.Equal(t, 50, trip.DistanceKM)
	assert.Equal(t, "Location A", trip.PickupLocation)
	assert.Equal(t, "Location B", trip.Destination)
	assert.Equal(t, "Driver X", trip.DriverName)
	assert.Equal(t, "Car Y", trip.VehicleName)
	assert.Equal(t, 100, trip.Amount)
	assert.Equal(t, float32(4.5), trip.Rating)
	assert.Equal(t, "Good trip", trip.Feedback)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPDFRepository_GetTripByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPDFRepository(db)

	tripID := "123"

	mock.ExpectQuery(`SELECT id, customer_name, booking_date, duration_minutes, distance_km, pickup_location, destination, driver_name, vehicle_name, amount, rating, feedback FROM trips WHERE id = \$1`).
		WithArgs(tripID).
		WillReturnError(sql.ErrNoRows)

	trip, err := repo.GetTripByID(tripID)

	assert.NoError(t, err)
	assert.Equal(t, model.Pdf{}, trip)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPDFRepository_GetTripByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPDFRepository(db)

	tripID := "123"

	mock.ExpectQuery(`SELECT id, customer_name, booking_date, duration_minutes, distance_km, pickup_location, destination, driver_name, vehicle_name, amount, rating, feedback FROM trips WHERE id = \$1`).
		WithArgs(tripID).
		WillReturnError(sql.ErrConnDone)

	trip, err := repo.GetTripByID(tripID)

	assert.Error(t, err)
	assert.Equal(t, model.Pdf{}, trip)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
