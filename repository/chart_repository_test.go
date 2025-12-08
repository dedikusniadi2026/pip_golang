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

func TestPopularDestinationRepository_GetAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.PopularDestinationRepository{DB: db}

	rows := sqlmock.NewRows([]string{"id", "destination", "bookings", "created_at"}).
		AddRow(1, "Jakarta", 100, time.Now()).
		AddRow(2, "Bandung", 80, time.Now())

	mock.ExpectQuery(`SELECT id, destination, bookings, created_at FROM popular_destinations ORDER BY bookings DESC`).
		WillReturnRows(rows)

	destinations, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, destinations, 2)
	assert.Equal(t, 1, destinations[0].ID)
	assert.Equal(t, "Jakarta", destinations[0].Destination)
	assert.Equal(t, 100, destinations[0].Bookings)
	assert.Equal(t, 2, destinations[1].ID)
	assert.Equal(t, "Bandung", destinations[1].Destination)
	assert.Equal(t, 80, destinations[1].Bookings)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPopularDestinationRepository_GetAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.PopularDestinationRepository{DB: db}

	mock.ExpectQuery(`SELECT id, destination, bookings, created_at FROM popular_destinations ORDER BY bookings DESC`).
		WillReturnError(sql.ErrConnDone)

	destinations, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, destinations)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPopularDestinationRepository_Add_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.PopularDestinationRepository{DB: db}

	pd := model.PopularDestination{
		Destination: "Jakarta",
		Bookings:    100,
	}

	mock.ExpectQuery(`INSERT INTO popular_destinations \(destination, bookings\) VALUES \(\$1, \$2\) RETURNING id, created_at`).
		WithArgs(pd.Destination, pd.Bookings).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	result, err := repo.Add(pd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Jakarta", result.Destination)
	assert.Equal(t, 100, result.Bookings)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPopularDestinationRepository_Add_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.PopularDestinationRepository{DB: db}

	pd := model.PopularDestination{
		Destination: "Jakarta",
		Bookings:    100,
	}

	mock.ExpectQuery(`INSERT INTO popular_destinations \(destination, bookings\) VALUES \(\$1, \$2\) RETURNING id, created_at`).
		WithArgs(pd.Destination, pd.Bookings).
		WillReturnError(sql.ErrConnDone)

	result, err := repo.Add(pd)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPopularDestinationRepository_UpdateBookings_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.PopularDestinationRepository{DB: db}

	id := 1
	bookings := 150

	mock.ExpectExec(`UPDATE popular_destinations SET bookings = \$1 WHERE id = \$2`).
		WithArgs(bookings, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateBookings(id, bookings)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPopularDestinationRepository_UpdateBookings_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.PopularDestinationRepository{DB: db}

	id := 1
	bookings := 150

	mock.ExpectExec(`UPDATE popular_destinations SET bookings = \$1 WHERE id = \$2`).
		WithArgs(bookings, id).
		WillReturnError(sql.ErrConnDone)

	err = repo.UpdateBookings(id, bookings)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPopularDestinationRepository_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.PopularDestinationRepository{DB: db}

	id := 1

	mock.ExpectExec(`DELETE FROM popular_destinations WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(id)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPopularDestinationRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.PopularDestinationRepository{DB: db}

	id := 1

	mock.ExpectExec(`DELETE FROM popular_destinations WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(id)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
