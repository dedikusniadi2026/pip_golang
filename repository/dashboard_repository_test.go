package repository_test

import (
	"auth-service/repository"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDashboardRepository_GetTotalBookings_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewDashboardRepository(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM booking`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	count, err := repo.GetTotalBookings()

	assert.NoError(t, err)
	assert.Equal(t, 10, count)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDashboardRepository_GetTotalBookings_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewDashboardRepository(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM booking`).
		WillReturnError(sql.ErrConnDone)

	count, err := repo.GetTotalBookings()

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDashboardRepository_GetActiveDrivers_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewDashboardRepository(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM drivers WHERE status = 'active'`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := repo.GetActiveDrivers()

	assert.NoError(t, err)
	assert.Equal(t, 5, count)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDashboardRepository_GetActiveDrivers_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewDashboardRepository(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM drivers WHERE status = 'active'`).
		WillReturnError(sql.ErrConnDone)

	count, err := repo.GetActiveDrivers()

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDashboardRepository_GetTotalRevenue_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewDashboardRepository(db)

	mock.ExpectQuery(`SELECT COALESCE\(SUM\(amount\), 0\)::FLOAT FROM payment`).
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1000.0))

	total, err := repo.GetTotalRevenue()

	assert.NoError(t, err)
	assert.Equal(t, 1000.0, total)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDashboardRepository_GetTotalRevenue_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewDashboardRepository(db)

	mock.ExpectQuery(`SELECT COALESCE\(SUM\(amount\), 0\)::FLOAT FROM payment`).
		WillReturnError(sql.ErrConnDone)

	total, err := repo.GetTotalRevenue()

	assert.Error(t, err)
	assert.Equal(t, 0.0, total)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
