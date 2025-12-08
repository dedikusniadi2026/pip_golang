package repository

import (
	"auth-service/model"
	"database/sql"
)

type CarTypeRepositoryInterface interface {
	FindAll() ([]model.CarType, error)
	GetByID(id int) (*model.CarType, error)
	Create(ct model.CarType) error
}

type CarTypeRepository struct {
	DB *sql.DB
}

func NewCarTypeRepository(db *sql.DB) *CarTypeRepository {
	return &CarTypeRepository{DB: db}
}

func (r *CarTypeRepository) FindAll() ([]model.CarType, error) {
	rows, err := r.DB.Query("SELECT id, type_name FROM car_type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.CarType
	for rows.Next() {
		var ct model.CarType
		if err := rows.Scan(&ct.ID, &ct.TypeName); err != nil {
			return nil, err
		}
		list = append(list, ct)
	}
	return list, nil
}

func (r *CarTypeRepository) GetByID(id int) (*model.CarType, error) {
	var cm model.CarType
	err := r.DB.QueryRow(`
		SELECT id, type_name
		FROM car_type
		WHERE id=$1
	`, id).Scan(&cm.ID, &cm.TypeName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &cm, nil
}

func (r *CarTypeRepository) Create(ct model.CarType) error {
	_, err := r.DB.Exec(
		"INSERT INTO car_type (type_name) VALUES ($1)",
		ct.TypeName,
	)
	return err
}
