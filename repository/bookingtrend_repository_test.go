package repository_test

import (
	"auth-service/repository"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestBookingTrendsRepository_GetTrends_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewBookingTrendsRepository(db)

	year := 2023
	rows := sqlmock.NewRows([]string{"month", "booking_count", "year"}).
		AddRow(1, 100, 2023).
		AddRow(2, 150, 2023)

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT month, booking_count, year
        FROM booking_trends
        WHERE year = $1
        ORDER BY id ASC
    `)).WithArgs(year).WillReturnRows(rows)

	trends, err := repo.GetTrends(year)

	assert.NoError(t, err)
	assert.Len(t, trends, 2)
	assert.Equal(t, "1", trends[0].Month)
	assert.Equal(t, 100, trends[0].Booking)
	assert.Equal(t, 2023, trends[0].Year)
	assert.Equal(t, "2", trends[1].Month)
	assert.Equal(t, 150, trends[1].Booking)
	assert.Equal(t, 2023, trends[1].Year)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBookingTrendsRepository_GetTrends_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewBookingTrendsRepository(db)

	year := 2023

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT month, booking_count, year
        FROM booking_trends
        WHERE year = $1
        ORDER BY id ASC
    `)).WithArgs(year).WillReturnError(sql.ErrConnDone)

	trends, err := repo.GetTrends(year)

	assert.Error(t, err)
	assert.Nil(t, trends)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
