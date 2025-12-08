package repository

import (
	"auth-service/model"
	"context"
	"database/sql"
)

type TripsRepositoryInterface interface {
	Create(t *model.VehicleTrip) error
	Update(t *model.VehicleTrip) error
	Delete(id uint) error
	FindByVehicle(vehicleID uint) ([]model.VehicleTrip, error)
	GetTripTotal(ctx context.Context) (*model.TotalTrips, error)
}

type TripsRepository struct {
	DB *sql.DB
}

func NewTripRepo(db *sql.DB) *TripsRepository {
	return &TripsRepository{DB: db}
}
func (r *TripsRepository) Create(t *model.VehicleTrip) error {
	query := `
		INSERT INTO vehicle_trips
		(vehicle_id, trip_date, origin, destination, rating, price, passenger_name, distance_km)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.DB.Exec(query,
		t.VehicleID,
		t.TripDate,
		t.Origin,
		t.Destination,
		t.Rating,
		t.Price,
		t.PassengerName,
		t.DistanceKM,
	)
	return err
}

func (r *TripsRepository) Update(t *model.VehicleTrip) error {
	query := `
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
	`
	_, err := r.DB.Exec(query,
		t.VehicleID,
		t.TripDate,
		t.Origin,
		t.Destination,
		t.Rating,
		t.Price,
		t.PassengerName,
		t.DistanceKM,
		t.ID,
	)
	return err
}

func (r *TripsRepository) Delete(id uint) error {
	query := `DELETE FROM vehicle_trips WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *TripsRepository) FindByID(id uint) (*model.VehicleTrip, error) {
	query := `
		SELECT id, vehicle_id, trip_date, origin, destination, rating, price, passenger_name, distance_km
		FROM vehicle_trips
		WHERE id = $1
	`
	row := r.DB.QueryRow(query, id)
	var t model.VehicleTrip
	err := row.Scan(
		&t.ID,
		&t.VehicleID,
		&t.TripDate,
		&t.Origin,
		&t.Destination,
		&t.Rating,
		&t.Price,
		&t.PassengerName,
		&t.DistanceKM,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *TripsRepository) FindByVehicle(vehicleID uint) ([]model.VehicleTrip, error) {
	query := `
		SELECT id, vehicle_id, trip_date, origin, destination, rating, price, passenger_name, distance_km
		FROM vehicle_trips
		WHERE vehicle_id = $1
	`
	rows, err := r.DB.Query(query, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var trips []model.VehicleTrip
	for rows.Next() {
		var t model.VehicleTrip
		err := rows.Scan(
			&t.ID,
			&t.VehicleID,
			&t.TripDate,
			&t.Origin,
			&t.Destination,
			&t.Rating,
			&t.Price,
			&t.PassengerName,
			&t.DistanceKM,
		)
		if err != nil {
			return nil, err
		}
		trips = append(trips, t)
	}
	return trips, nil
}

func (r *TripsRepository) GetTripTotal(ctx context.Context) (*model.TotalTrips, error) {
	stats := &model.TotalTrips{}

	err := r.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM vehicle_trips;
	`).Scan(&stats.TotalTrips)
	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(distance_km), 0)
		FROM vehicle_trips;
	`).Scan(&stats.TotalDistance)
	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(price), 0)
		FROM vehicle_trips;
	`).Scan(&stats.TotalRevenue)
	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, `
		SELECT COALESCE(AVG(rating), 0)
		FROM vehicle_trips;
	`).Scan(&stats.AverageRating)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
