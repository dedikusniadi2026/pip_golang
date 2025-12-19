package repository

import (
	"auth-service/model"
	"database/sql"
)

type AssignmentsRepositoryInterface interface {
	Create(a *model.DriverAssignment) error
	FindByVehicle(vehicleID uint) ([]model.DriverAssignment, error)
	Update(a *model.DriverAssignment) error
	Delete(id uint) error
}

type AssignmentsRepository struct {
	DB *sql.DB
}

func NewAssignmentsRepository(db *sql.DB) *AssignmentsRepository {
	return &AssignmentsRepository{DB: db}
}

func (r *AssignmentsRepository) Create(a *model.DriverAssignment) error {
	query := `INSERT INTO driver_assignments
			 (vehicle_id, start_date, end_date, total_trips, driver_name, status)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 RETURNING id
			`
	err := r.DB.QueryRow(query,
		a.VehicleID,
		a.StartDate,
		a.EndDate,
		a.TotalTrips,
		a.DriverName,
		a.Status,
	).Scan(&a.ID)
	return err
}

func (r *AssignmentsRepository) FindByVehicle(vehicleID uint) ([]model.DriverAssignment, error) {
	query := `
    SELECT id, vehicle_id, start_date, end_date, total_trips, driver_name, status
    FROM driver_assignments
    WHERE vehicle_id = $1`

	rows, err := r.DB.Query(query, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.DriverAssignment

	for rows.Next() {
		var a model.DriverAssignment
		err := rows.Scan(
			&a.ID,
			&a.VehicleID,
			&a.StartDate,
			&a.EndDate,
			&a.TotalTrips,
			&a.DriverName,
			&a.Status,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}

	return list, nil
}

func (r *AssignmentsRepository) Update(a *model.DriverAssignment) error {
	query := `
        UPDATE driver_assignments
        SET vehicle_id = $1,
    		start_date = $2,
    		end_date = $3,
    		total_trips = $4,
    		driver_name = $5,
    		status = $6
		WHERE id = $7
    `
	_, err := r.DB.Exec(query,
		a.VehicleID,
		a.StartDate,
		a.EndDate,
		a.TotalTrips,
		a.DriverName,
		a.Status,
		a.ID,
	)
	return err
}

func (r *AssignmentsRepository) Delete(id uint) error {
	query := `DELETE FROM driver_assignments WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
