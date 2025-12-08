package repository

import (
	"auth-service/model"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type BookingRepositoryInterface interface {
	Create(*model.Booking) error
	GetAll() ([]model.Booking, error)
	Update(*model.Booking) error
	Delete(id string) error
}

type BookingRepository struct {
	DB *sql.DB
}

func generateBookingID() string {
	t := time.Now()
	return fmt.Sprintf("BK%04d%02d%02d%02d%02d%02d%03d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
		t.Nanosecond()/1e6,
	)
}

func (r *BookingRepository) Create(b *model.Booking) error {
	if b.ID == "" {
		b.ID = generateBookingID()
	}

	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	b.Payment = strings.Title(strings.ToLower(strings.TrimSpace(b.Payment)))
	b.Status = strings.Title(strings.ToLower(strings.TrimSpace(b.Status)))

	_, err := r.DB.Exec(
		`INSERT INTO booking
        (id, customer, driver, place, date, price, status, payment, phone_number, pickup_location, drop_location, pickup_time, amount, notes, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		b.ID,
		b.Customer,
		b.Driver,
		b.Place,
		b.Date,
		b.Price,
		b.Status,
		b.Payment,
		b.PhoneNumber,
		b.PickupLocation,
		b.DropLocation,
		b.PickupTime,
		b.Amount,
		b.Notes,
		b.CreatedAt,
		b.UpdatedAt,
	)

	return err
}

func (r *BookingRepository) GetAll() ([]model.Booking, error) {
	var bookings []model.Booking

	rows, err := r.DB.Query(`SELECT id, customer, driver, place, date, price, status, payment, phone_number, pickup_location, drop_location, pickup_time, amount, notes, created_at, updated_at FROM booking`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b model.Booking
		if err := rows.Scan(
			&b.ID,
			&b.Customer,
			&b.Driver,
			&b.Place,
			&b.Date,
			&b.Price,
			&b.Status,
			&b.Payment,
			&b.PhoneNumber,
			&b.PickupLocation,
			&b.DropLocation,
			&b.PickupTime,
			&b.Amount,
			&b.Notes,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	return bookings, nil
}

func (r *BookingRepository) Update(b *model.Booking) error {
	b.UpdatedAt = time.Now()

	b.Payment = strings.Title(strings.ToLower(strings.TrimSpace(b.Payment)))
	b.Status = strings.Title(strings.ToLower(strings.TrimSpace(b.Status)))

	_, err := r.DB.Exec(
		`UPDATE booking SET
			customer = $1,
			driver = $2,
			place = $3,
			date = $4,
			price = $5,
			status = $6,
			payment = $7,
			phone_number = $8,
			pickup_location = $9,
			drop_location = $10,
			pickup_time = $11,
			amount = $12,
			notes = $13,
			updated_at = $14
		WHERE id = $15`,
		b.Customer,
		b.Driver,
		b.Place,
		b.Date,
		b.Price,
		b.Status,
		b.Payment,
		b.PhoneNumber,
		b.PickupLocation,
		b.DropLocation,
		b.PickupTime,
		b.Amount,
		b.Notes,
		b.UpdatedAt,
		b.ID,
	)

	return err
}

func (r *BookingRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM booking WHERE id = $1`, id)
	return err
}
