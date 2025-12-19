package repository

import (
	"auth-service/model"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTokenRepository_Save_Success(t *testing.T) {
	// mock db
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TokenRepositoryImpl{
		DB: db,
	}

	token := &model.RefreshToken{
		UserID:    1,
		TokenHash: "hashed-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	mock.ExpectExec(`
		INSERT INTO refresh_tokens
		\(user_id, token_hash, expires_at\)
		VALUES \(\$1, \$2, \$3\)
	`).
		WithArgs(
			token.UserID,
			token.TokenHash,
			token.ExpiresAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Save(token)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenRepository_Save_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &TokenRepositoryImpl{
		DB: db,
	}

	token := &model.RefreshToken{
		UserID:    1,
		TokenHash: "hashed-token",
		ExpiresAt: time.Now(),
	}

	mock.ExpectExec(`
		INSERT INTO refresh_tokens
		\(user_id, token_hash, expires_at\)
		VALUES \(\$1, \$2, \$3\)
	`).
		WithArgs(
			token.UserID,
			token.TokenHash,
			token.ExpiresAt,
		).
		WillReturnError(sql.ErrConnDone)

	err = repo.Save(token)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
