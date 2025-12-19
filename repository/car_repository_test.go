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

func TestCarRepository_GetAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarRepository(db)

	mock.ExpectQuery(`SELECT id, brand, model, year, plate_number, capacity, color, driver_id, last_maintenance_date, current_km FROM vehicles`).
		WillReturnError(sql.ErrConnDone)

	cars, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, cars)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarRepository_GetAll_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository.CarRepository{DB: db}

	rows := sqlmock.NewRows([]string{
		"id",
		"brand",
		"model",
		"year",
		"plate_number",
		"capacity",
		"color",
		"driver_id",
		"last_maintenance_date",
		"current_km",
	}).AddRow(
		1,
		"Toyota",
		"Avanza",
		"INVALID_YEAR",
		"D 1234 AA",
		7,
		"Black",
		1,
		time.Now(),
		10000,
	)

	mock.ExpectQuery(`FROM vehicles`).
		WillReturnRows(rows)

	result, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCarRepository_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarRepository(db)

	id := 1
	lastMaintenance := time.Now()
	rows := sqlmock.NewRows([]string{"id", "brand", "model", "year", "plate_number", "capacity", "color", "driver_id", "last_maintenance_date", "current_km"}).
		AddRow(1, "Toyota", "Camry", 2020, "ABC123", 5, "Blue", 1, lastMaintenance, 10000)

	mock.ExpectQuery(`SELECT id, brand, model, year, plate_number, capacity, color, driver_id, last_maintenance_date, current_km FROM vehicles WHERE id=\$1`).
		WithArgs(id).
		WillReturnRows(rows)

	car, err := repo.GetByID(id)

	assert.NoError(t, err)
	assert.NotNil(t, car)
	assert.Equal(t, 1, car.ID)
	assert.Equal(t, "Toyota", car.Brand)
	assert.Equal(t, "Camry", car.Model)
	assert.Equal(t, 2020, car.Year)
	assert.Equal(t, "ABC123", car.PlateNumber)
	assert.Equal(t, 5, car.Capacity)
	assert.Equal(t, "Blue", car.Color)
	assert.Equal(t, 1, car.DriverID)
	assert.Equal(t, lastMaintenance, *car.LastMaintenanceDate)
	assert.Equal(t, 10000, car.CurrentKM)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarRepository(db)

	id := 1

	mock.ExpectQuery(`SELECT id, brand, model, year, plate_number, capacity, color, driver_id, last_maintenance_date, current_km FROM vehicles WHERE id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	car, err := repo.GetByID(id)

	assert.NoError(t, err)
	assert.Nil(t, car)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarRepository_GetByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarRepository(db)

	id := 1

	mock.ExpectQuery(`SELECT id, brand, model, year, plate_number, capacity, color, driver_id, last_maintenance_date, current_km FROM vehicles WHERE id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	car, err := repo.GetByID(id)

	assert.Error(t, err)
	assert.Nil(t, car)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarRepository(db)

	lastMaintenance := time.Now()
	car := model.Car{
		Brand:               "Toyota",
		Model:               "Camry",
		Year:                2020,
		PlateNumber:         "ABC123",
		Capacity:            5,
		Color:               "Blue",
		DriverID:            1,
		LastMaintenanceDate: &lastMaintenance,
		CurrentKM:           10000,
	}

	mock.ExpectExec(`INSERT INTO vehicles \(brand, model, year, plate_number, capacity, color, driver_id, last_maintenance_date, current_km\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\)`).
		WithArgs(car.Brand, car.Model, car.Year, car.PlateNumber, car.Capacity, car.Color, car.DriverID, car.LastMaintenanceDate, car.CurrentKM).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(car)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarRepository(db)

	lastMaintenance := time.Now()
	car := model.Car{
		Brand:               "Toyota",
		Model:               "Camry",
		Year:                2020,
		PlateNumber:         "ABC123",
		Capacity:            5,
		Color:               "Blue",
		DriverID:            1,
		LastMaintenanceDate: &lastMaintenance,
		CurrentKM:           10000,
	}

	mock.ExpectExec(`INSERT INTO vehicles \(brand, model, year, plate_number, capacity, color, driver_id, last_maintenance_date, current_km\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\)`).
		WithArgs(car.Brand, car.Model, car.Year, car.PlateNumber, car.Capacity, car.Color, car.DriverID, car.LastMaintenanceDate, car.CurrentKM).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(car)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarRepository_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarRepository(db)

	id := 1
	lastMaintenance := time.Now()
	car := model.Car{
		Brand:               "Toyota",
		Model:               "Camry",
		Year:                2020,
		PlateNumber:         "ABC123",
		Capacity:            5,
		Color:               "Blue",
		DriverID:            1,
		LastMaintenanceDate: &lastMaintenance,
		CurrentKM:           10000,
	}

	mock.ExpectExec(`UPDATE vehicles SET brand=\$1, model=\$2, year=\$3, plate_number=\$4, capacity=\$5, color=\$6, driver_id=\$7, last_maintenance_date=\$8, current_km=\$9 WHERE id=\$10`).
		WithArgs(car.Brand, car.Model, car.Year, car.PlateNumber, car.Capacity, car.Color, car.DriverID, car.LastMaintenanceDate, car.CurrentKM, id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(id, car)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarRepository_Update_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarRepository(db)

	id := 1
	lastMaintenance := time.Now()
	car := model.Car{
		Brand:               "Toyota",
		Model:               "Camry",
		Year:                2020,
		PlateNumber:         "ABC123",
		Capacity:            5,
		Color:               "Blue",
		DriverID:            1,
		LastMaintenanceDate: &lastMaintenance,
		CurrentKM:           10000,
	}

	mock.ExpectExec(`UPDATE vehicles SET brand=\$1, model=\$2, year=\$3, plate_number=\$4, capacity=\$5, color=\$6, driver_id=\$7, last_maintenance_date=\$8, current_km=\$9 WHERE id=\$10`).
		WithArgs(car.Brand, car.Model, car.Year, car.PlateNumber, car.Capacity, car.Color, car.DriverID, car.LastMaintenanceDate, car.CurrentKM, id).
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(id, car)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarRepository_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarRepository(db)

	id := 1

	mock.ExpectExec(`DELETE FROM vehicles WHERE id=\$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(id)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarRepository(db)

	id := 1

	mock.ExpectExec(`DELETE FROM vehicles WHERE id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(id)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarRepository_GetAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository.CarRepository{DB: db}

	rows := sqlmock.NewRows([]string{
		"id",
		"brand",
		"model",
		"year",
		"plate_number",
		"capacity",
		"color",
		"driver_id",
		"last_maintenance_date",
		"current_km",
	}).AddRow(
		1,
		"Toyota",
		"Avanza",
		2022,
		"D 1234 AA",
		7,
		"Black",
		1,
		time.Now(),
		15000,
	)

	mock.ExpectQuery("FROM vehicles").
		WillReturnRows(rows)

	result, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Toyota", result[0].Brand)
	assert.NoError(t, mock.ExpectationsWereMet())
}
