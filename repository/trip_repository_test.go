package repository_test

import (
	"auth-service/model"
	"auth-service/repository"
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTripsRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripRepo(db)

	trip := &model.VehicleTrip{
		VehicleID:     1,
		TripDate:      time.Now(),
		Origin:        "Jakarta",
		Destination:   "Bandung",
		Rating:        5,
		Price:         100000,
		PassengerName: "John Doe",
		DistanceKM:    150,
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO vehicle_trips
		(vehicle_id, trip_date, origin, destination, rating, price, passenger_name, distance_km)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)).
		WithArgs(trip.VehicleID, trip.TripDate, trip.Origin, trip.Destination, trip.Rating, trip.Price, trip.PassengerName, trip.DistanceKM).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(trip)
	assert.NoError(t, err)
}

func TestTripsRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripRepo(db)

	trip := &model.VehicleTrip{
		ID:            1,
		VehicleID:     1,
		TripDate:      time.Now(),
		Origin:        "Jakarta",
		Destination:   "Bandung",
		Rating:        4,
		Price:         120000,
		PassengerName: "Jane Doe",
		DistanceKM:    160,
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE vehicle_trips
		SET vehicle_id = $1,
			trip_date = $2,
			origin = $3,
			destination = $4,
			rating = $5,
			price = $6
			passenger_name = $7
			distance_km = $8
		WHERE id = $9
	`)).
		WithArgs(trip.VehicleID, trip.TripDate, trip.Origin, trip.Destination, trip.Rating, trip.Price, trip.PassengerName, trip.DistanceKM, trip.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(trip)
	assert.NoError(t, err)
}

func TestTripsRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM vehicle_trips WHERE id = $1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(1)
	assert.NoError(t, err)
}

func TestTripsRepository_FindByVehicle(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripRepo(db)

	tripDate := time.Now()

	rows := sqlmock.NewRows([]string{"id", "vehicle_id", "trip_date", "origin", "destination", "rating", "price", "passenger_name", "distance_km"}).
		AddRow(1, 1, tripDate, "Jakarta", "Bandung", 5, 100000, "John Doe", 150)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, vehicle_id, trip_date, origin, destination, rating, price, passenger_name, distance_km
		FROM vehicle_trips
		WHERE vehicle_id = $1
	`)).WithArgs(1).WillReturnRows(rows)

	result, err := repo.FindByVehicle(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "John Doe", result[0].PassengerName)
}

func TestTripsRepository_FindByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripRepo(db)

	id := uint(1)
	tripDate := time.Now()

	rows := sqlmock.NewRows([]string{"id", "vehicle_id", "trip_date", "origin", "destination", "rating", "price", "passenger_name", "distance_km"}).
		AddRow(1, 1, tripDate, "Jakarta", "Bandung", 5, 100000, "John Doe", 150)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, vehicle_id, trip_date, origin, destination, rating, price, passenger_name, distance_km
		FROM vehicle_trips
		WHERE id = $1
	`)).WithArgs(id).WillReturnRows(rows)

	result, err := repo.FindByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, "John Doe", result.PassengerName)
}

func TestTripsRepository_FindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripRepo(db)

	id := uint(1)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, vehicle_id, trip_date, origin, destination, rating, price, passenger_name, distance_km
		FROM vehicle_trips
		WHERE id = $1
	`)).WithArgs(id).WillReturnError(sql.ErrNoRows)

	result, err := repo.FindByID(id)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestTripsRepository_FindByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripRepo(db)

	id := uint(1)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, vehicle_id, trip_date, origin, destination, rating, price, passenger_name, distance_km
		FROM vehicle_trips
		WHERE id = $1
	`)).WithArgs(id).WillReturnError(sql.ErrConnDone)

	result, err := repo.FindByID(id)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestTripsRepository_GetTripTotal(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTripRepo(db)

	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(*)
		FROM vehicle_trips;
	`)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COALESCE(SUM(distance_km), 0)
		FROM vehicle_trips;
	`)).WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(1500))

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COALESCE(SUM(price), 0)
		FROM vehicle_trips;
	`)).WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(1000000.0))

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COALESCE(AVG(rating), 0)
		FROM vehicle_trips;
	`)).WillReturnRows(sqlmock.NewRows([]string{"avg"}).AddRow(4.5))

	result, err := repo.GetTripTotal(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), result.TotalTrips)
	assert.Equal(t, int64(1500), result.TotalDistance)
	assert.Equal(t, 1000000.0, result.TotalRevenue)
	assert.Equal(t, 4.5, result.AverageRating)
}
