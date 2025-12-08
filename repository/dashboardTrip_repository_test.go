package repository_test

import (
	"auth-service/repository"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDashboardTripRepository_GetDashboardSummary_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewDashboardTripRepository(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM vehicle_trips`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(50))

	mock.ExpectQuery(`SELECT COALESCE\(SUM\(distance_km\), 0\) FROM vehicle_trips`).
		WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(1000))

	mock.ExpectQuery(`SELECT COALESCE\(SUM\(duration\), 0\) FROM vehicle_trips`).
		WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(200))

	ctx := context.Background()
	summary, err := repo.GetDashboardSummary(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 50, summary.TotalTrips)
	assert.Equal(t, 1000, summary.TotalDistance)
	assert.Equal(t, 200, summary.TotalDuration) // Note: Code has bug, assigns to TotalTrips, but testing as TotalDuration

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDashboardTripRepository_GetDashboardSummary_ErrorOnTotalTrips(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewDashboardTripRepository(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM vehicle_trips`).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	summary, err := repo.GetDashboardSummary(ctx)

	assert.Error(t, err)
	assert.Nil(t, summary)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDashboardTripRepository_GetDashboardSummary_ErrorOnTotalDistance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewDashboardTripRepository(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM vehicle_trips`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(50))

	mock.ExpectQuery(`SELECT COALESCE\(SUM\(distance_km\), 0\) FROM vehicle_trips`).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	summary, err := repo.GetDashboardSummary(ctx)

	assert.Error(t, err)
	assert.Nil(t, summary)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDashboardTripRepository_GetDashboardSummary_ErrorOnTotalDuration(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewDashboardTripRepository(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM vehicle_trips`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(50))

	mock.ExpectQuery(`SELECT COALESCE\(SUM\(distance_km\), 0\) FROM vehicle_trips`).
		WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(1000))

	mock.ExpectQuery(`SELECT COALESCE\(SUM\(duration\), 0\) FROM vehicle_trips`).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	summary, err := repo.GetDashboardSummary(ctx)

	assert.Error(t, err)
	assert.Nil(t, summary)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
