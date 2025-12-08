package repository_test

import (
	"auth-service/model"
	"auth-service/repository"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTripHistoryRepository_GetTripHistory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripHistoryRepository(db)

	rows := sqlmock.NewRows([]string{"id", "booking_code", "customer_name", "booking_date", "duration_minutes", "distance_km", "pickup_location", "destination", "driver_name", "vehicle_name", "amount", "rating", "feedback"}).
		AddRow(1, "ABC123", "John Doe", "2023-01-01", 60, 50, "Location A", "Location B", "Driver X", "Car Y", 100, 4.5, "Good trip")

	mock.ExpectQuery(`
        SELECT id, booking_code, customer_name, booking_date,
               duration_minutes, distance_km, pickup_location, destination,
               driver_name, vehicle_name, amount, rating, feedback
        FROM trips
    `).
		WillReturnRows(rows)

	tripHistories, err := repo.GetTripHistory()

	assert.NoError(t, err)
	assert.NotNil(t, tripHistories)
	assert.Len(t, tripHistories, 1)
	assert.Equal(t, 1, tripHistories[0].ID)
	assert.Equal(t, "ABC123", tripHistories[0].BookingCode)
	assert.Equal(t, "John Doe", tripHistories[0].CustomerName)
	assert.Equal(t, "2023-01-01", tripHistories[0].BookingDate)
	assert.Equal(t, 60, tripHistories[0].DurationMinutes)
	assert.Equal(t, 50, tripHistories[0].DistanceKM)
	assert.Equal(t, "Location A", tripHistories[0].PickupLocation)
	assert.Equal(t, "Location B", tripHistories[0].Destination)
	assert.Equal(t, "Driver X", tripHistories[0].DriverName)
	assert.Equal(t, "Car Y", tripHistories[0].VehicleName)
	assert.Equal(t, 100, tripHistories[0].Amount)
	assert.Equal(t, float32(4.5), tripHistories[0].Rating)
	assert.Equal(t, "Good trip", tripHistories[0].Feedback)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestTripHistoryRepository_GetTripHistory_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripHistoryRepository(db)

	mock.ExpectQuery(`
        SELECT id, booking_code, customer_name, booking_date,
               duration_minutes, distance_km, pickup_location, destination,
               driver_name, vehicle_name, amount, rating, feedback
        FROM trips
    `).
		WillReturnError(sql.ErrNoRows)

	tripHistories, err := repo.GetTripHistory()

	assert.Error(t, err)
	assert.Nil(t, tripHistories)
	assert.Equal(t, sql.ErrNoRows, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestTripHistoryRepository_GetTripHistory_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripHistoryRepository(db)

	rows := sqlmock.NewRows([]string{"id", "booking_code", "customer_name", "booking_date", "duration_minutes", "distance_km", "pickup_location", "destination", "driver_name", "vehicle_name", "amount", "rating", "feedback"})

	mock.ExpectQuery(`
        SELECT id, booking_code, customer_name, booking_date,
               duration_minutes, distance_km, pickup_location, destination,
               driver_name, vehicle_name, amount, rating, feedback
        FROM trips
    `).
		WillReturnRows(rows)

	tripHistories, err := repo.GetTripHistory()

	assert.NoError(t, err)
	assert.Equal(t, []model.TripHistory{}, tripHistories)
	assert.Len(t, tripHistories, 0)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestTripHistoryRepository_GetTripHistory_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripHistoryRepository(db)

	mock.ExpectQuery(`
        SELECT id, booking_code, customer_name, booking_date,
               duration_minutes, distance_km, pickup_location, destination,
               driver_name, vehicle_name, amount, rating, feedback
        FROM trips
    `).
		WillReturnError(sql.ErrConnDone)

	tripHistories, err := repo.GetTripHistory()

	assert.Error(t, err)
	assert.Nil(t, tripHistories)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
