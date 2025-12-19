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

type MockDriverRepository struct {
	ErrToReturn error
}

func (m *MockDriverRepository) GetAll() ([]model.Driver, error) {
	if m.ErrToReturn != nil {
		return nil, m.ErrToReturn
	}
	return []model.Driver{{ID: 1, Name: "John"}}, nil
}

func TestDriverRepository_GetAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "address",
		"driver_license_number", "car_model_id", "car_type_id", "plate_number",
		"status", "created_at", "updated_at",
	}).AddRow(
		1, "John Doe", "john@example.com", "1234567890", "Address 1",
		"DL123", "1", "1", "ABC123", "active", time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT id, name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number, status, created_at, updated_at FROM drivers").
		WillReturnRows(rows)

	drivers, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, drivers, 1)
	assert.Equal(t, 1, drivers[0].ID)
	assert.Equal(t, "John Doe", drivers[0].Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDriverRepository_GetAll_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	mock.ExpectQuery("(?i)SELECT .* FROM drivers").
		WillReturnError(fmt.Errorf("query failed"))

	drivers, err := repo.GetAll()
	assert.Error(t, err)
	assert.Nil(t, drivers)
	assert.Contains(t, err.Error(), "query failed")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDriverRepository_GetAll_ScanErrorMock(t *testing.T) {
	mockRepo := &MockDriverRepository{ErrToReturn: fmt.Errorf("scan failed")}
	drivers, err := mockRepo.GetAll()
	assert.Error(t, err)
	assert.Nil(t, drivers)
	assert.Contains(t, err.Error(), "scan failed")
}

func TestDriverRepository_GetAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	mock.ExpectQuery("SELECT id, name, email, phone, address, driver_license_number, car_model_id, car_type_id, plate_number, status, created_at, updated_at FROM drivers").
		WillReturnError(sql.ErrConnDone)

	drivers, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, drivers)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
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

	mock.ExpectQuery(`INSERT INTO drivers .* RETURNING id, created_at, updated_at`).
		WithArgs(driver.Name, driver.Email, driver.Phone, driver.Address, driver.DriverLicenseNumber, driver.CarModelID, driver.CarTypeID, driver.PlateNumber, driver.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, time.Now(), time.Now()))

	err = repo.Create(driver)

	assert.NoError(t, err)
	assert.Equal(t, 1, driver.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDriverRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	driver := &model.Driver{Name: "John Doe"}

	mock.ExpectQuery(`INSERT INTO drivers .* RETURNING id, created_at, updated_at`).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(driver)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDriverRepository_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	id := "1"
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "address",
		"driver_license_number", "car_model_id", "car_type_id", "plate_number",
		"status", "created_at", "updated_at",
	}).AddRow(
		1, "John Doe", "john@example.com", "1234567890", "Address 1",
		"DL123", "1", "1", "ABC123", "active", time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT .* FROM drivers WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(rows)

	driver, err := repo.GetByID(id)

	assert.NoError(t, err)
	assert.NotNil(t, driver)
	assert.Equal(t, 1, driver.ID)
	assert.Equal(t, "John Doe", driver.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDriverRepository_GetByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}
	id := "1"

	mock.ExpectQuery(`SELECT .* FROM drivers WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	driver, err := repo.GetByID(id)
	assert.Error(t, err)
	assert.Nil(t, driver)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDriverRepository_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	id := "1"
	driver := &model.Driver{Name: "Jane Doe"}

	mock.ExpectExec(`UPDATE drivers .* WHERE id=\$10`).
		WithArgs(driver.Name, driver.Email, driver.Phone, driver.Address, driver.DriverLicenseNumber, driver.CarModelID, driver.CarTypeID, driver.PlateNumber, driver.Status, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(id, driver)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDriverRepository_Update_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}
	id := "1"
	driver := &model.Driver{}

	mock.ExpectExec(`UPDATE drivers .* WHERE id=\$10`).WillReturnError(sql.ErrConnDone)

	err = repo.Update(id, driver)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
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
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDriverRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}
	id := "1"

	mock.ExpectExec(`DELETE FROM drivers WHERE id=\$1`).WithArgs(id).WillReturnError(sql.ErrConnDone)

	err = repo.Delete(id)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDriverRepository_GetAll_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "address",
		"driver_license_number", "car_model_id", "car_type_id", "plate_number", "status", "created_at", "updated_at",
	}).AddRow(
		"invalid_int", "John Doe", "john@example.com", "1234567890", "Address 1",
		"DL123", 1, 1, "ABC123", "active", time.Now(), time.Now(),
	)

	mock.ExpectQuery("SELECT .* FROM drivers").WillReturnRows(rows)

	drivers, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, drivers)
	assert.Contains(t, err.Error(), "converting driver.Value type string")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDriverRepository_GetAll_RowsErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.DriverRepository{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "phone", "address",
		"driver_license_number", "car_model_id", "car_type_id", "plate_number",
		"status", "created_at", "updated_at",
	}).AddRow(1, "John Doe", "john@example.com", "1234567890", "Address 1",
		"DL123", 1, 1, "ABC123", "active", time.Now(), time.Now(),
	).CloseError(fmt.Errorf("rows iteration error"))

	mock.ExpectQuery("SELECT .* FROM drivers").WillReturnRows(rows)

	drivers, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, drivers)
	assert.Contains(t, err.Error(), "rows iteration error")

	assert.NoError(t, mock.ExpectationsWereMet())
}
