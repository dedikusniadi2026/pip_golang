package repository_test

import (
	"auth-service/model"
	"auth-service/repository"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPaymentRepository_GetPaymentStats_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	rows := sqlmock.NewRows([]string{"total_payment"}).AddRow(1000)
	mock.ExpectQuery(`SELECT CAST\(COALESCE\(SUM\(amount\), 0\) AS BIGINT\) FROM payment WHERE status = 'paid'`).
		WillReturnRows(rows)

	rows2 := sqlmock.NewRows([]string{"pending_payment"}).AddRow(500)
	mock.ExpectQuery(`SELECT CAST\(COALESCE\(SUM\(amount\), 0\) AS BIGINT\) FROM payment WHERE status = 'pending'`).
		WillReturnRows(rows2)

	rows3 := sqlmock.NewRows([]string{"total_transactions"}).AddRow(10)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM payment`).
		WillReturnRows(rows3)

	stats, err := repo.GetPaymentStats(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, int64(1000), stats.TotalPayment)
	assert.Equal(t, int64(500), stats.PendingPayment)
	assert.Equal(t, int64(10), stats.TotalTransactions)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_GetPaymentStats_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	mock.ExpectQuery(`SELECT CAST\(COALESCE\(SUM\(amount\), 0\) AS BIGINT\) FROM payment WHERE status = 'paid'`).
		WillReturnError(sql.ErrConnDone)

	stats, err := repo.GetPaymentStats(context.Background())

	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_GetPayments_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	page := 1
	pageSize := 10
	rows := sqlmock.NewRows([]string{"payment_id", "booking_id", "customer", "driver", "amount", "method", "status", "payment_date"}).
		AddRow(1, 1, "Customer1", "Driver1", 100.0, "Credit", "paid", "2023-01-01").
		AddRow(2, 2, "Customer2", "Driver2", 200.0, "Cash", "pending", "2023-01-02")

	mock.ExpectQuery(`SELECT payment_id, booking_id, customer, driver, amount, method, status, payment_date FROM payment LIMIT \$1 OFFSET \$2`).
		WithArgs(pageSize, 0).
		WillReturnRows(rows)

	payments, err := repo.GetPayments(page, pageSize)

	assert.NoError(t, err)
	assert.Len(t, payments, 2)
	assert.Equal(t, 1, payments[0].PaymentID)
	assert.Equal(t, "Customer1", payments[0].Customer)
	assert.Equal(t, 2, payments[1].PaymentID)
	assert.Equal(t, "Customer2", payments[1].Customer)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_GetPayments_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	page := 1
	pageSize := 10

	mock.ExpectQuery(`SELECT payment_id, booking_id, customer, driver, amount, method, status, payment_date FROM payment LIMIT \$1 OFFSET \$2`).
		WithArgs(pageSize, 0).
		WillReturnError(sql.ErrConnDone)

	payments, err := repo.GetPayments(page, pageSize)

	assert.Error(t, err)
	assert.Nil(t, payments)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_GetAll_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	rows := sqlmock.NewRows([]string{"payment_id", "booking_id", "customer", "driver", "amount", "method", "status", "payment_date"}).
		AddRow(1, 1, "Customer1", "Driver1", 100.0, "Credit", "paid", "2023-01-01").
		AddRow(2, 2, "Customer2", "Driver2", 200.0, "Cash", "pending", "2023-01-02")

	mock.ExpectQuery(`SELECT payment_id, booking_id, customer, driver, amount, method, status, payment_date FROM payment`).
		WillReturnRows(rows)

	payments, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, payments, 2)
	assert.Equal(t, 1, payments[0].PaymentID)
	assert.Equal(t, "Customer1", payments[0].Customer)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_GetAll_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	mock.ExpectQuery(`SELECT payment_id, booking_id, customer, driver, amount, method, status, payment_date FROM payment`).
		WillReturnError(sql.ErrConnDone)

	payments, err := repo.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, payments)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	id := 1
	rows := sqlmock.NewRows([]string{"payment_id", "booking_id", "customer", "driver", "amount", "method", "status", "payment_date"}).
		AddRow(1, 1, "Customer1", "Driver1", 100.0, "Credit", "paid", "2023-01-01")

	mock.ExpectQuery(`SELECT payment_id, booking_id, customer, driver, amount, method, status, payment_date FROM payment WHERE payment_id=\$1`).
		WithArgs(id).
		WillReturnRows(rows)

	payment, err := repo.GetByID(context.Background(), id)

	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, 1, payment.PaymentID)
	assert.Equal(t, "Customer1", payment.Customer)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	id := 1

	mock.ExpectQuery(`SELECT payment_id, booking_id, customer, driver, amount, method, status, payment_date FROM payment WHERE payment_id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	payment, err := repo.GetByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Nil(t, payment)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_GetByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	id := 1

	mock.ExpectQuery(`SELECT payment_id, booking_id, customer, driver, amount, method, status, payment_date FROM payment WHERE payment_id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	payment, err := repo.GetByID(context.Background(), id)

	assert.Error(t, err)
	assert.Nil(t, payment)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	payment := &model.Payment{
		BookingID: 1,
		Customer:  "Customer1",
		Driver:    "Driver1",
		Amount:    100.0,
		Method:    "Credit",
		Status:    "paid",
	}

	mock.ExpectQuery(`INSERT INTO payment \(booking_id, customer, driver, amount, method, status\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING payment_id`).
		WithArgs(payment.BookingID, payment.Customer, payment.Driver, payment.Amount, payment.Method, payment.Status).
		WillReturnRows(sqlmock.NewRows([]string{"payment_id"}).AddRow(1))

	id, err := repo.Create(context.Background(), payment)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	payment := &model.Payment{
		BookingID: 1,
		Customer:  "Customer1",
		Driver:    "Driver1",
		Amount:    100.0,
		Method:    "Credit",
		Status:    "paid",
	}

	mock.ExpectQuery(`INSERT INTO payment \(booking_id, customer, driver, amount, method, status\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING payment_id`).
		WithArgs(payment.BookingID, payment.Customer, payment.Driver, payment.Amount, payment.Method, payment.Status).
		WillReturnError(sql.ErrConnDone)

	id, err := repo.Create(context.Background(), payment)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	payment := &model.Payment{
		PaymentID: 1,
		BookingID: 1,
		Customer:  "Customer1",
		Driver:    "Driver1",
		Amount:    100.0,
		Method:    "Credit",
		Status:    "paid",
	}

	mock.ExpectExec(`UPDATE payment SET booking_id=\$1, customer=\$2, driver=\$3, amount=\$4, method=\$5, status=\$6 WHERE payment_id=\$7`).
		WithArgs(payment.BookingID, payment.Customer, payment.Driver, payment.Amount, payment.Method, payment.Status, payment.PaymentID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(context.Background(), payment)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_Update_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	payment := &model.Payment{
		PaymentID: 1,
		BookingID: 1,
		Customer:  "Customer1",
		Driver:    "Driver1",
		Amount:    100.0,
		Method:    "Credit",
		Status:    "paid",
	}

	mock.ExpectExec(`UPDATE payment SET booking_id=\$1, customer=\$2, driver=\$3, amount=\$4, method=\$5, status=\$6 WHERE payment_id=\$7`).
		WithArgs(payment.BookingID, payment.Customer, payment.Driver, payment.Amount, payment.Method, payment.Status, payment.PaymentID).
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(context.Background(), payment)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_Delete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	id := 1

	mock.ExpectExec(`DELETE FROM payment WHERE payment_id=\$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), id)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPaymentRepository_Delete_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)

	id := 1

	mock.ExpectExec(`DELETE FROM payment WHERE payment_id=\$1`).
		WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(context.Background(), id)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
