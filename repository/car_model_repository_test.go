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

func TestCarModelRepository_FindAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarModelRepository(db)

	mock.ExpectQuery(`SELECT id, model_name, created_at, updated_at FROM car_model`).
		WillReturnError(sql.ErrConnDone)

	carModels, err := repo.FindAll()

	assert.Error(t, err)
	assert.Nil(t, carModels)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarModelRepository_FindAll_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository.CarModelRepository{DB: db}

	rows := sqlmock.NewRows([]string{
		"id",
		"model_name",
		"created_at",
	}).AddRow(
		1,
		"Avanza",
		time.Now(),
	)

	mock.ExpectQuery("SELECT id, model_name, created_at, updated_at FROM car_model").
		WillReturnRows(rows)

	result, err := repo.FindAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCarModelRepository_FindAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository.CarModelRepository{DB: db}

	rows := sqlmock.NewRows([]string{
		"id",
		"model_name",
		"created_at",
		"updated_at",
	}).AddRow(
		1,
		"Avanza",
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery("SELECT id, model_name, created_at, updated_at FROM car_model").
		WillReturnRows(rows)

	result, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Avanza", result[0].ModelName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCarModelRepository_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarModelRepository(db)

	id := 1
	rows := sqlmock.NewRows([]string{"id", "model_name", "created_at", "updated_at"}).
		AddRow(1, "Model A", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT id, model_name, created_at, updated_at FROM car_model WHERE id=\$1`).
		WithArgs(id).
		WillReturnRows(rows)

	carModel, err := repo.GetByID(id)

	assert.NoError(t, err)
	assert.NotNil(t, carModel)
	assert.Equal(t, 1, carModel.ID)
	assert.Equal(t, "Model A", carModel.ModelName)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarModelRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarModelRepository(db)

	id := 1

	mock.ExpectQuery(`SELECT id, model_name, created_at, updated_at FROM car_model WHERE id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	carModel, err := repo.GetByID(id)

	assert.NoError(t, err)
	assert.Nil(t, carModel)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarModelRepository_GetByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarModelRepository(db)

	id := 1

	mock.ExpectQuery(`SELECT id, model_name, created_at, updated_at FROM car_model WHERE id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	carModel, err := repo.GetByID(id)

	assert.Error(t, err)
	assert.Nil(t, carModel)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarModelRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarModelRepository(db)

	carModel := &model.CarModel{
		ModelName: "Model A",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`INSERT INTO car_model \(model_name, created_at, updated_at\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(carModel.ModelName, carModel.CreatedAt, carModel.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.Create(carModel)

	assert.NoError(t, err)
	assert.Equal(t, 1, carModel.ID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarModelRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarModelRepository(db)

	carModel := &model.CarModel{
		ModelName: "Model A",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`INSERT INTO car_model \(model_name, created_at, updated_at\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(carModel.ModelName, carModel.CreatedAt, carModel.UpdatedAt).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(carModel)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
