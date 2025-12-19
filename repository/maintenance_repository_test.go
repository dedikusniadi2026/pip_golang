package repository_test

import (
	"auth-service/model"
	"auth-service/repository"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMaintenanceRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewMaintenanceRepository(db)

	maintenance := &model.VehicleMaintenance{
		VehicleID:   1,
		ServiceDate: time.Now(),
		Description: "Oil change",
		Cost:        100.0,
		Mileage:     5000,
		ServiceType: "Maintenance",
	}

	mock.ExpectExec(`INSERT INTO vehicle_maintenance \(vehicle_id, service_date, description, cost, mileage, service_type\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
		WithArgs(maintenance.VehicleID, maintenance.ServiceDate, maintenance.Description, maintenance.Cost, maintenance.Mileage, maintenance.ServiceType).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(maintenance)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestMaintenanceRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewMaintenanceRepository(db)

	maintenance := &model.VehicleMaintenance{
		VehicleID:   1,
		ServiceDate: time.Now(),
		Description: "Oil change",
		Cost:        100.0,
		Mileage:     5000,
		ServiceType: "Maintenance",
	}

	mock.ExpectExec(`INSERT INTO vehicle_maintenance \(vehicle_id, service_date, description, cost, mileage, service_type\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
		WithArgs(maintenance.VehicleID, maintenance.ServiceDate, maintenance.Description, maintenance.Cost, maintenance.Mileage, maintenance.ServiceType).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(maintenance)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestMaintenanceRepository_FindByVehicle_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewMaintenanceRepository(db)

	vehicleID := 1
	rows := sqlmock.NewRows([]string{"id", "vehicle_id", "service_date", "description", "cost", "mileage", "service_type"}).
		AddRow(1, 1, time.Now(), "Oil change", 100.0, 5000, "Maintenance")

	mock.ExpectQuery(`SELECT id, vehicle_id, service_date, description, cost , mileage , service_type FROM vehicle_maintenance WHERE vehicle_id = \$1 ORDER BY service_date DESC`).
		WithArgs(vehicleID).
		WillReturnRows(rows)

	maintenances, err := repo.FindByVehicle(vehicleID)

	assert.NoError(t, err)
	assert.Len(t, maintenances, 1)
	assert.Equal(t, uint(1), maintenances[0].ID)
	assert.Equal(t, "Oil change", maintenances[0].Description)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestMaintenanceRepository_FindByVehicle_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewMaintenanceRepository(db)

	vehicleID := 1

	mock.ExpectQuery(`SELECT id, vehicle_id, service_date, description, cost , mileage , service_type FROM vehicle_maintenance WHERE vehicle_id = \$1 ORDER BY service_date DESC`).
		WithArgs(vehicleID).
		WillReturnError(sql.ErrConnDone)

	maintenances, err := repo.FindByVehicle(vehicleID)

	assert.Error(t, err)
	assert.Nil(t, maintenances)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestMaintenanceRepository_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewMaintenanceRepository(db)

	maintenance := &model.VehicleMaintenance{
		ID:          1,
		VehicleID:   1,
		ServiceDate: time.Now(),
		Description: "Brake check",
		Cost:        200.0,
		Mileage:     6000,
		ServiceType: "Repair",
	}

	mock.ExpectExec(`UPDATE vehicle_maintenance SET vehicle_id = \$1, service_date = \$2, description = \$3, cost = \$4, mileage = \$5, service_type = \$6 WHERE id = \$7`).
		WithArgs(maintenance.VehicleID, maintenance.ServiceDate, maintenance.Description, maintenance.Cost, maintenance.Mileage, maintenance.ServiceType, maintenance.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(maintenance)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestMaintenanceRepository_Update_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewMaintenanceRepository(db)

	maintenance := &model.VehicleMaintenance{
		ID:          1,
		VehicleID:   1,
		ServiceDate: time.Now(),
		Description: "Brake check",
		Cost:        200.0,
		Mileage:     6000,
		ServiceType: "Repair",
	}

	mock.ExpectExec(`UPDATE vehicle_maintenance SET vehicle_id = \$1, service_date = \$2, description = \$3, cost = \$4, mileage = \$5, service_type = \$6 WHERE id = \$7`).
		WithArgs(maintenance.VehicleID, maintenance.ServiceDate, maintenance.Description, maintenance.Cost, maintenance.Mileage, maintenance.ServiceType, maintenance.ID).
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(maintenance)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestMaintenanceRepository_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewMaintenanceRepository(db)

	id := uint(1)

	mock.ExpectExec(`DELETE FROM vehicle_maintenance WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(id)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestMaintenanceRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewMaintenanceRepository(db)

	id := uint(1)

	mock.ExpectExec(`DELETE FROM vehicle_maintenance WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(id)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestMaintenanceRepository_FindByVehicle_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.MaintenanceRepository{DB: db}

	mock.ExpectQuery("SELECT .* FROM vehicle_maintenance WHERE vehicle_id = .* ORDER BY service_date DESC").
		WithArgs(1).
		WillReturnError(fmt.Errorf("query failed"))

	list, err := repo.FindByVehicle(1)

	assert.Error(t, err)
	assert.Nil(t, list)
	assert.Contains(t, err.Error(), "query failed")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaintenanceRepository_FindByVehicle_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository.MaintenanceRepository{DB: db}

	rows := sqlmock.NewRows([]string{
		"id",
		"vehicle_id",
		"service_date",
		"description",
		"cost",
		"mileage",
		"service_type",
	}).AddRow(
		1,
		1,
		time.Now(),
		"Routine service",
		"INVALID_COST",
		12000,
		"OIL_CHANGE",
	)

	mock.ExpectQuery(`FROM vehicle_maintenance`).
		WithArgs(1).
		WillReturnRows(rows)

	result, err := repo.FindByVehicle(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
