package repository

import (
	"auth-service/model"
	"database/sql"
)

type TokenRepositoryImpl struct {
	DB *sql.DB
}

func (r *TokenRepositoryImpl) Save(token *model.RefreshToken) error {
	_, err := r.DB.Exec(`
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, token.UserID, token.TokenHash, token.ExpiresAt)

	return err
}
