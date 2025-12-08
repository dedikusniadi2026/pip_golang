package repository

import (
	"auth-service/model"
	"database/sql"
)

type DriverRepositoryInterface interface {
	GetAll() ([]model.Driver, error)
	Create(d *model.Driver) error
	GetByID(id string) (*model.Driver, error)
	Update(id string, d *model.Driver) error
	Delete(id string) error
}

type DriverRepository struct {
	DB *sql.DB
}

func (r *DriverRepository) GetAll() ([]model.Driver, error) {
	rows, err := r.DB.Query("SELECT id, name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number, status, created_at, updated_at FROM drivers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []model.Driver
	for rows.Next() {
		var d model.Driver
		if err := rows.Scan(&d.ID, &d.Name, &d.Email, &d.Phone, &d.Address, &d.DriverLicenseNumber, &d.CarModelID, &d.CarTypeID, &d.PlateNumber,
			&d.Status, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return drivers, nil
}

func (r *DriverRepository) Create(d *model.Driver) error {
	query := `
    INSERT INTO drivers (name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number,status, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
    RETURNING id, created_at, updated_at;
`

	return r.DB.QueryRow(
		query,
		d.Name, d.Email, d.Phone, d.Address, d.DriverLicenseNumber, d.CarModelID, d.CarTypeID, d.PlateNumber, d.Status,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

func (r *DriverRepository) GetByID(id string) (*model.Driver, error) {
	var d model.Driver
	err := r.DB.QueryRow("SELECT id, name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number, status, created_at, updated_at FROM drivers WHERE id = $1", id).Scan(&d.ID, &d.Name, &d.Email, &d.Phone, &d.Address, &d.DriverLicenseNumber, &d.CarModelID, &d.CarTypeID, &d.PlateNumber, &d.Status, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DriverRepository) Update(id string, d *model.Driver) error {
	query := `
        UPDATE drivers
        SET name=$1, email=$2, phone=$3, address=$4, driver_license_number=$5, car_model_id=$6, car_type_id=$7, plate_number=$8,
		status=$9, updated_at=NOW()
        WHERE id=$10
    `
	_, err := r.DB.Exec(query, d.Name, d.Email, d.Phone, d.Address, d.DriverLicenseNumber, d.CarModelID, d.CarTypeID, d.PlateNumber, d.Status, id)
	return err
}

func (r *DriverRepository) Delete(id string) error {
	query := `DELETE FROM drivers WHERE id=$1`
	_, err := r.DB.Exec(query, id)
	return err
}
