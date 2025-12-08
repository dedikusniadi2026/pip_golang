package repository

import (
	"auth-service/model"
	"database/sql"
)

type MaintenanceRepositoryInterface interface {
	Create(*model.VehicleMaintenance) error
	FindByVehicle(vehicleID int) ([]model.VehicleMaintenance, error)
	Update(*model.VehicleMaintenance) error
	Delete(id uint) error
}

type MaintenanceRepository struct {
	DB *sql.DB
}

func NewMaintenanceRepository(db *sql.DB) *MaintenanceRepository {
	return &MaintenanceRepository{DB: db}
}

func (r *MaintenanceRepository) Create(m *model.VehicleMaintenance) error {
	query := `
        INSERT INTO vehicle_maintenance (vehicle_id, service_date, description, cost, mileage, service_type)
        VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.DB.Exec(query, m.VehicleID, m.ServiceDate, m.Description, m.Cost, m.Mileage, m.ServiceType)
	return err
}

func (r *MaintenanceRepository) FindByVehicle(vehicleID int) ([]model.VehicleMaintenance, error) {
	query := `
        SELECT id, vehicle_id, service_date, description, cost , mileage , service_type
        FROM vehicle_maintenance
        WHERE vehicle_id = $1
        ORDER BY service_date DESC`

	rows, err := r.DB.Query(query, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.VehicleMaintenance

	for rows.Next() {
		var m model.VehicleMaintenance
		if err := rows.Scan(&m.ID, &m.VehicleID, &m.ServiceDate, &m.Description, &m.Cost, &m.Mileage, &m.ServiceType); err != nil {
			return nil, err
		}
		list = append(list, m)
	}

	return list, nil
}

func (r *MaintenanceRepository) Update(m *model.VehicleMaintenance) error {
	query := `
        UPDATE vehicle_maintenance
        SET vehicle_id = $1,
            service_date = $2,
            description = $3,
            cost = $4,
            mileage = $5,
            service_type = $6
        WHERE id = $7`

	_, err := r.DB.Exec(query,
		m.VehicleID,
		m.ServiceDate,
		m.Description,
		m.Cost,
		m.Mileage,
		m.ServiceType,
		m.ID,
	)

	return err
}

func (r *MaintenanceRepository) Delete(id uint) error {
	query := `DELETE FROM vehicle_maintenance WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
