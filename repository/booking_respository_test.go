package repository_test

import (
	"auth-service/model"
	"auth-service/repository"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestBookingRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.BookingRepository{DB: db}

	booking := &model.Booking{
		Customer:       "John Doe",
		Driver:         "Driver1",
		Place:          "Location A",
		Date:           "2023-10-01",
		Price:          "100.00",
		Status:         "pending",
		Payment:        "cash",
		PhoneNumber:    stringPtr("1234567890"),
		PickupLocation: stringPtr("Pickup"),
		DropLocation:   stringPtr("Drop"),
		PickupTime:     stringPtr("10:00"),
		Amount:         float64Ptr(100.0),
		Notes:          stringPtr("Test note"),
	}

	mock.ExpectExec(`INSERT INTO booking`).
		WithArgs(sqlmock.AnyArg(), "John Doe", "Driver1", "Location A", "2023-10-01", "100.00", "Pending", "Cash", "1234567890", "Pickup", "Drop", "10:00", 100.0, "Test note", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(booking)

	assert.NoError(t, err)
	assert.NotEmpty(t, booking.ID)
	assert.True(t, booking.CreatedAt.After(time.Time{}))
	assert.True(t, booking.UpdatedAt.After(time.Time{}))

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBookingRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.BookingRepository{DB: db}

	booking := &model.Booking{
		Customer: "John Doe",
		Status:   "pending",
		Payment:  "cash",
	}

	mock.ExpectExec(`INSERT INTO booking`).
		WillReturnError(sql.ErrConnDone)

	err = repo.Create(booking)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBookingRepository_GetAll_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.BookingRepository{DB: db}

	mock.ExpectQuery("FROM booking").
		WillReturnError(errors.New("query error"))

	bookings, err := repo.GetAll()

	assert.Nil(t, bookings)
	assert.Error(t, err)
}

func TestBookingRepository_GetAll_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.BookingRepository{DB: db}

	rows := sqlmock.NewRows([]string{
		"id", "customer",
	}).AddRow(1, "John")

	mock.ExpectQuery("FROM booking").
		WillReturnRows(rows)

	bookings, err := repo.GetAll()

	assert.Nil(t, bookings)
	assert.Error(t, err)
}

func TestBookingRepository_GetAll_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.BookingRepository{DB: db}

	rows := sqlmock.NewRows([]string{
		"id", "customer", "driver", "place", "date",
		"price", "status", "payment", "phone_number",
		"pickup_location", "drop_location", "pickup_time",
		"amount", "notes", "created_at", "updated_at",
	}).AddRow(
		1,
		"John",
		"Driver A",
		"Bandung",
		time.Now(),
		100000,
		"CONFIRMED",
		"CASH",
		"08123456789",
		"A",
		"B",
		time.Now(),
		100000,
		"OK",
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery("FROM booking").
		WillReturnRows(rows)

	bookings, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, bookings, 1)
}

func TestBookingRepository_GetAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.BookingRepository{DB: db}

	mock.ExpectQuery(`SELECT id, customer, driver, place, date, price, status, payment, phone_number, pickup_location, drop_location, pickup_time, amount, notes, created_at, updated_at FROM booking`).
		WillReturnError(sql.ErrConnDone)

	bookings, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, bookings)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBookingRepository_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.BookingRepository{DB: db}

	booking := &model.Booking{
		ID:             "BK123",
		Customer:       "Jane Doe",
		Driver:         "Driver2",
		Place:          "Location B",
		Date:           "2023-10-02",
		Price:          "200.00",
		Status:         "confirmed",
		Payment:        "card",
		PhoneNumber:    stringPtr("0987654321"),
		PickupLocation: stringPtr("New Pickup"),
		DropLocation:   stringPtr("New Drop"),
		PickupTime:     stringPtr("11:00"),
		Amount:         float64Ptr(200.0),
		Notes:          stringPtr("Updated note"),
	}

	mock.ExpectExec(`UPDATE booking SET`).
		WithArgs("Jane Doe", "Driver2", "Location B", "2023-10-02", "200.00", "Confirmed", "Card", "0987654321", "New Pickup", "New Drop", "11:00", 200.0, "Updated note", sqlmock.AnyArg(), "BK123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(booking)

	assert.NoError(t, err)
	assert.True(t, booking.UpdatedAt.After(time.Time{}))

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBookingRepository_Update_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.BookingRepository{DB: db}

	booking := &model.Booking{
		ID:       "BK123",
		Customer: "Jane Doe",
		Status:   "confirmed",
		Payment:  "card",
	}

	mock.ExpectExec(`UPDATE booking SET`).
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(booking)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBookingRepository_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.BookingRepository{DB: db}

	id := "BK123"

	mock.ExpectExec(`DELETE FROM booking WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(id)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBookingRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.BookingRepository{DB: db}

	id := "BK123"

	mock.ExpectExec(`DELETE FROM booking WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(id)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Helper functions for pointers
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
