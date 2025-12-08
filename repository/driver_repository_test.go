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

func TestDriverRepository_GetAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "phone", "address", "driver_license_number", "car_model_id", "car_type_id", "plate_number", "status", "created_at", "updated_at"}).
		AddRow("1", "John Doe", "john@example.com", "1234567890", "Address 1", "DL123", 1, 1, "ABC123", "active", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT id, name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number, status, created_at, updated_at FROM drivers`).
		WillReturnRows(rows)

	drivers, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, drivers, 1)
	assert.Equal(t, 1, drivers[0].ID)
	assert.Equal(t, "John Doe", drivers[0].Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDriverRepository_GetAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	mock.ExpectQuery(`SELECT id, name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number, status, created_at, updated_at FROM drivers`).
		WillReturnError(sql.ErrConnDone)

	drivers, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, drivers)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDriverRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	driver := &model.Driver{
		Name:                "John Doe",
		Email:               "john@example.com",
		Phone:               "1234567890",
		Address:             "Address 1",
		DriverLicenseNumber: "DL123",
		CarModelID:          "1",
		CarTypeID:           "1",
		PlateNumber:         "ABC123",
		Status:              "active",
	}

	mock.ExpectQuery(`INSERT INTO drivers \(name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number,status, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, NOW\(\), NOW\(\)\) RETURNING id, created_at, updated_at`).
		WithArgs(driver.Name, driver.Email, driver.Phone, driver.Address, driver.DriverLicenseNumber, driver.CarModelID, driver.CarTypeID, driver.PlateNumber, driver.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("1", time.Now(), time.Now()))

	err = repo.Create(driver)

	assert.NoError(t, err)
	assert.Equal(t, 1, driver.ID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDriverRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	driver := &model.Driver{
		Name: "John Doe",
	}

	mock.ExpectQuery(`INSERT INTO drivers \(name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number,status, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, NOW\(\), NOW\(\)\) RETURNING id, created_at, updated_at`).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(driver)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDriverRepository_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	id := "1"
	rows := sqlmock.NewRows([]string{"id", "name", "email", "phone", "address", "driver_license_number", "car_model_id", "car_type_id", "plate_number", "status", "created_at", "updated_at"}).
		AddRow("1", "John Doe", "john@example.com", "1234567890", "Address 1", "DL123", 1, 1, "ABC123", "active", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT id, name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number, status, created_at, updated_at FROM drivers WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(rows)

	driver, err := repo.GetByID(id)

	assert.NoError(t, err)
	assert.NotNil(t, driver)
	assert.Equal(t, 1, driver.ID)
	assert.Equal(t, "John Doe", driver.Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDriverRepository_GetByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	id := "1"

	mock.ExpectQuery(`SELECT id, name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number, status, created_at, updated_at FROM drivers WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	driver, err := repo.GetByID(id)

	assert.Error(t, err)
	assert.Nil(t, driver)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDriverRepository_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	id := "1"
	driver := &model.Driver{
		Name: "Jane Doe",
	}

	mock.ExpectExec(`UPDATE drivers SET name=\$1, email=\$2, phone=\$3, address=\$4, driver_license_number=\$5, car_model_id=\$6, car_type_id=\$7, plate_number=\$8, status=\$9, updated_at=NOW\(\) WHERE id=\$10`).
		WithArgs(driver.Name, driver.Email, driver.Phone, driver.Address, driver.DriverLicenseNumber, driver.CarModelID, driver.CarTypeID, driver.PlateNumber, driver.Status, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(id, driver)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDriverRepository_Update_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	id := "1"
	driver := &model.Driver{}

	mock.ExpectExec(`UPDATE drivers SET name=\$1, email=\$2, phone=\$3, address=\$4, driver_license_number=\$5, car_model_id=\$6, car_type_id=\$7, plate_number=\$8, status=\$9, updated_at=NOW\(\) WHERE id=\$10`).
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(id, driver)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDriverRepository_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	id := "1"

	mock.ExpectExec(`DELETE FROM drivers WHERE id=\$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(id)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDriverRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	id := "1"

	mock.ExpectExec(`DELETE FROM drivers WHERE id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(id)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
