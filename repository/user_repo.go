package repository

import (
	"auth-service/model"
	"database/sql"
)

type UserRepositoryInterface interface {
	GetUserWithRole(username string) (*model.User, error)
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) GetUserWithRole(username string) (*model.User, error) {
	var user model.User

	query := `
		SELECT u.id, u.username, u.password, r.role_name
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE u.username = $1
	`

	err := r.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
