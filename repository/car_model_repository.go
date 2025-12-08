package repository

import (
	"auth-service/model"
	"database/sql"
)

type CarModelRepositoryInterface interface {
	FindAll() ([]model.CarModel, error)
	GetByID(id int) (*model.CarModel, error)
	Create(cm *model.CarModel) error
}

type CarModelRepository struct {
	DB *sql.DB
}

func NewCarModelRepository(db *sql.DB) *CarModelRepository {
	return &CarModelRepository{DB: db}
}

func (r *CarModelRepository) FindAll() ([]model.CarModel, error) {
	rows, err := r.DB.Query("SELECT id, model_name, created_at, updated_at FROM car_model")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.CarModel
	for rows.Next() {
		var cm model.CarModel
		if err := rows.Scan(&cm.ID, &cm.ModelName, &cm.CreatedAt, &cm.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, cm)
	}
	return list, nil
}

func (r *CarModelRepository) GetByID(id int) (*model.CarModel, error) {
	var cm model.CarModel
	err := r.DB.QueryRow(`
		SELECT id, model_name, created_at, updated_at
		FROM car_model
		WHERE id=$1
	`, id).Scan(&cm.ID, &cm.ModelName, &cm.CreatedAt, &cm.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &cm, nil
}

func (r *CarModelRepository) Create(cm *model.CarModel) error {
	query := `INSERT INTO car_model (model_name, created_at, updated_at)
	          VALUES ($1, $2, $3) RETURNING id`
	return r.DB.QueryRow(query, cm.ModelName, cm.CreatedAt, cm.UpdatedAt).Scan(&cm.ID)
}
