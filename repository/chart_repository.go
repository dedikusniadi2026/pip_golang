package repository

import (
	"database/sql"

	"auth-service/model"
)

type PopularDestinationRepositoryInterface interface {
	GetAll() ([]model.PopularDestination, error)
	Add(pd model.PopularDestination) (*model.PopularDestination, error)
	UpdateBookings(id int, bookings int) error
	Delete(id int) error
}

type PopularDestinationRepository struct {
	DB *sql.DB
}

func (r *PopularDestinationRepository) GetAll() ([]model.PopularDestination, error) {
	rows, err := r.DB.Query("SELECT id, destination, bookings, created_at FROM popular_destinations ORDER BY bookings DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.PopularDestination
	for rows.Next() {
		var pd model.PopularDestination
		if err := rows.Scan(&pd.ID, &pd.Destination, &pd.Bookings, &pd.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, pd)
	}
	return result, nil
}

func (r *PopularDestinationRepository) Add(pd model.PopularDestination) (*model.PopularDestination, error) {
	err := r.DB.QueryRow(
		"INSERT INTO popular_destinations (destination, bookings) VALUES ($1, $2) RETURNING id, created_at",
		pd.Destination, pd.Bookings,
	).Scan(&pd.ID, &pd.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &pd, nil
}

func (r *PopularDestinationRepository) UpdateBookings(id int, bookings int) error {
	_, err := r.DB.Exec("UPDATE popular_destinations SET bookings = $1 WHERE id = $2", bookings, id)
	return err
}

func (r *PopularDestinationRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM popular_destinations WHERE id = $1", id)
	return err
}
