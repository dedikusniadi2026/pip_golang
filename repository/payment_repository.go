package repository

import (
	"auth-service/model"
	"context"
	"database/sql"
	"errors"
)

type PaymentRepositoryInterface interface {
	GetPayments(page, pageSize int) ([]model.Payment, error)
	GetPaymentStats(ctx context.Context) (*model.PaymentStats, error)
	GetAll(ctx context.Context) ([]model.Payment, error)
	GetByID(ctx context.Context, id int) (*model.Payment, error)
	Create(ctx context.Context, p *model.Payment) (int, error)
	Update(ctx context.Context, p *model.Payment) error
	Delete(ctx context.Context, id int) error
}

type PaymentRepository struct {
	DB *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

func (r *PaymentRepository) GetPaymentStats(ctx context.Context) (*model.PaymentStats, error) {
	stats := &model.PaymentStats{}

	err := r.DB.QueryRowContext(ctx, `
	SELECT 
		CAST(COALESCE(SUM(amount), 0) AS BIGINT) 
		FROM payment 
	WHERE status = 'paid';
	`).Scan(&stats.TotalPayment)
	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, `
        SELECT 
			CAST(COALESCE(SUM(amount), 0) AS BIGINT) 
			FROM payment 
		WHERE status = 'pending';

    `).Scan(&stats.PendingPayment)
	if err != nil {
		return nil, err
	}

	err = r.DB.QueryRowContext(ctx, `
        SELECT COUNT(*)
        FROM payment
    `).Scan(&stats.TotalTransactions)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *PaymentRepository) GetPayments(page, pageSize int) ([]model.Payment, error) {
	offset := (page - 1) * pageSize
	query := `SELECT payment_id, booking_id, customer, driver, amount, method, status, payment_date FROM payment LIMIT $1 OFFSET $2`
	rows, err := r.DB.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(
			&p.PaymentID,
			&p.BookingID,
			&p.Customer,
			&p.Driver,
			&p.Amount,
			&p.Method,
			&p.Status,
			&p.PaymentDate,
		); err != nil {
			return nil, err
		}

		payments = append(payments, p)
	}

	return payments, nil
}

func (r *PaymentRepository) GetAll(ctx context.Context) ([]model.Payment, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT payment_id, booking_id, customer, driver, amount, method, status, payment_date 
		 FROM payment`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment

	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(
			&p.PaymentID,
			&p.BookingID,
			&p.Customer,
			&p.Driver,
			&p.Amount,
			&p.Method,
			&p.Status,
			&p.PaymentDate,
		); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id int) (*model.Payment, error) {
	var p model.Payment

	err := r.DB.QueryRowContext(ctx,
		`SELECT payment_id, booking_id, customer, driver, amount, method, status, payment_date 
		 FROM payment 
		 WHERE payment_id=$1`,
		id,
	).Scan(
		&p.PaymentID,
		&p.BookingID,
		&p.Customer,
		&p.Driver,
		&p.Amount,
		&p.Method,
		&p.Status,
		&p.PaymentDate,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PaymentRepository) Create(ctx context.Context, p *model.Payment) (int, error) {
	var id int

	err := r.DB.QueryRowContext(ctx,
		`INSERT INTO payment (booking_id, customer, driver, amount, method, status) 
         VALUES ($1,$2,$3,$4,$5,$6) 
         RETURNING payment_id`,
		p.BookingID,
		p.Customer,
		p.Driver,
		p.Amount,
		p.Method,
		p.Status,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PaymentRepository) Update(ctx context.Context, p *model.Payment) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE payment 
		 SET booking_id=$1, customer=$2, driver=$3, amount=$4, method=$5, status=$6
		 WHERE payment_id=$7`,
		p.BookingID,
		p.Customer,
		p.Driver,
		p.Amount,
		p.Method,
		p.Status,
		p.PaymentID,
	)

	return err
}

func (r *PaymentRepository) Delete(ctx context.Context, id int) error {
	_, err := r.DB.ExecContext(ctx,
		`DELETE FROM payment WHERE payment_id=$1`, id)
	return err
}
