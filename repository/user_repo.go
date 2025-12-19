package repository

import (
	"auth-service/model"
	"database/sql"
)

type UserRepositoryImpl struct {
	DB *sql.DB
}

func (r *UserRepositoryImpl) Save(user *model.User) error {
	_, err := r.DB.Exec(`
		INSERT INTO users (username, password, role)
		VALUES ($1, $2, $3)
	`, user.Username, user.Password, user.Role)
	return err
}

func (r *UserRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, username, password, role
		FROM users
		WHERE username = $1
	`, username)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type tokenRepository struct {
	DB *sql.DB
}

func (r *tokenRepository) Save(token *model.RefreshToken) error {
	_, err := r.DB.Exec(`
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, token.UserID, token.TokenHash, token.ExpiresAt)
	return err
}
