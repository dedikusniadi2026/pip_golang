package repository

import (
	"auth-service/model"
	"database/sql"
)

type CarRepositoryInterface interface {
	GetAll() ([]model.Car, error)
	GetByID(id int) (*model.Car, error)
	Create(v model.Car) error
	Update(id int, v model.Car) error
	Delete(id int) error
}

type CarRepository struct {
	DB *sql.DB
}

func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{DB: db}
}

func (r *CarRepository) GetAll() ([]model.Car, error) {
	rows, err := r.DB.Query(`
		SELECT id, brand, model, year, plate_number, capacity, color,
		       driver_id, last_maintenance_date, current_km
		FROM vehicles
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []model.Car

	for rows.Next() {
		var v model.Car
		err := rows.Scan(
			&v.ID,
			&v.Brand,
			&v.Model,
			&v.Year,
			&v.PlateNumber,
			&v.Capacity,
			&v.Color,
			&v.DriverID,
			&v.LastMaintenanceDate,
			&v.CurrentKM,
		)
		if err != nil {
			return nil, err
		}
		cars = append(cars, v)
	}

	return cars, nil
}

func (r *CarRepository) GetByID(id int) (*model.Car, error) {
	var v model.Car
	err := r.DB.QueryRow(`
		SELECT id, brand, model, year, plate_number, capacity, color,
		       driver_id, last_maintenance_date, current_km
		FROM vehicles WHERE id=$1
	`, id).Scan(
		&v.ID, &v.Brand, &v.Model, &v.Year,
		&v.PlateNumber, &v.Capacity, &v.Color,
		&v.DriverID, &v.LastMaintenanceDate, &v.CurrentKM,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &v, nil
}

func (r *CarRepository) Create(v model.Car) error {
	_, err := r.DB.Exec(`
		INSERT INTO vehicles 
		(brand, model, year, plate_number, capacity, color, driver_id, last_maintenance_date, current_km)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`,
		v.Brand, v.Model, v.Year, v.PlateNumber,
		v.Capacity, v.Color, v.DriverID,
		v.LastMaintenanceDate, v.CurrentKM,
	)
	return err
}

func (r *CarRepository) Update(id int, v model.Car) error {
	_, err := r.DB.Exec(`
		UPDATE vehicles SET 
			brand=$1, model=$2, year=$3, plate_number=$4,
			capacity=$5, color=$6, driver_id=$7,
			last_maintenance_date=$8, current_km=$9
		WHERE id=$10
	`,
		v.Brand, v.Model, v.Year, v.PlateNumber,
		v.Capacity, v.Color, v.DriverID,
		v.LastMaintenanceDate, v.CurrentKM,
		id,
	)

	return err
}

func (r *CarRepository) Delete(id int) error {
	_, err := r.DB.Exec(`DELETE FROM vehicles WHERE id=$1`, id)
	return err
}
