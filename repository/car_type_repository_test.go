package repository_test

import (
	"auth-service/model"
	"auth-service/repository"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCarTypeRepository_FindAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarTypeRepository(db)

	rows := sqlmock.NewRows([]string{"id", "type_name"}).
		AddRow(1, "Sedan").
		AddRow(2, "SUV")

	mock.ExpectQuery(`SELECT id, type_name FROM car_type`).
		WillReturnRows(rows)

	carTypes, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, carTypes, 2)
	assert.Equal(t, "1", carTypes[0].ID)
	assert.Equal(t, "Sedan", carTypes[0].TypeName)
	assert.Equal(t, "2", carTypes[1].ID)
	assert.Equal(t, "SUV", carTypes[1].TypeName)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarTypeRepository_FindAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarTypeRepository(db)

	mock.ExpectQuery(`SELECT id, type_name FROM car_type`).
		WillReturnError(sql.ErrConnDone)

	carTypes, err := repo.FindAll()

	assert.Error(t, err)
	assert.Nil(t, carTypes)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarTypeRepository_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarTypeRepository(db)

	id := 1
	rows := sqlmock.NewRows([]string{"id", "type_name"}).
		AddRow(1, "Sedan")

	mock.ExpectQuery(`SELECT id, type_name FROM car_type WHERE id=\$1`).
		WithArgs(id).
		WillReturnRows(rows)

	carType, err := repo.GetByID(id)

	assert.NoError(t, err)
	assert.NotNil(t, carType)
	assert.Equal(t, "1", carType.ID)
	assert.Equal(t, "Sedan", carType.TypeName)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarTypeRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarTypeRepository(db)

	id := 1

	mock.ExpectQuery(`SELECT id, type_name FROM car_type WHERE id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	carType, err := repo.GetByID(id)

	assert.NoError(t, err)
	assert.Nil(t, carType)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarTypeRepository_GetByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarTypeRepository(db)

	id := 1

	mock.ExpectQuery(`SELECT id, type_name FROM car_type WHERE id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	carType, err := repo.GetByID(id)

	assert.Error(t, err)
	assert.Nil(t, carType)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarTypeRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarTypeRepository(db)

	carType := model.CarType{
		TypeName: "Sedan",
	}

	mock.ExpectExec(`INSERT INTO car_type \(type_name\) VALUES \(\$1\)`).
		WithArgs(carType.TypeName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(carType)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarTypeRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarTypeRepository(db)

	carType := model.CarType{
		TypeName: "Sedan",
	}

	mock.ExpectExec(`INSERT INTO car_type \(type_name\) VALUES \(\$1\)`).
		WithArgs(carType.TypeName).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(carType)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCarTypeRepository_FindAll_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarTypeRepository(db)

	mock.ExpectQuery("SELECT id, type_name FROM car_type").
		WillReturnError(errors.New("db error"))

	result, err := repo.FindAll()

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestCarTypeRepository_FindAll_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewCarTypeRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("abc")

	mock.ExpectQuery(`SELECT id, type_name FROM car_type`).
		WillReturnRows(rows)

	result, err := repo.FindAll()

	assert.Nil(t, result)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
