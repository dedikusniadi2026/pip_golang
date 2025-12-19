package repository

import (
	"auth-service/model"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	return db, mock
}

// ================= UserRepository =================

func TestUserRepository_FindByUsername_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &UserRepositoryImpl{DB: db}

	rows := sqlmock.NewRows([]string{"id", "username", "password", "role"}).
		AddRow(1, "admin", "hashed-password", "ADMIN")

	mock.ExpectQuery(`
		SELECT id, username, password, role
		FROM users
		WHERE username = \$1
	`).WithArgs("admin").WillReturnRows(rows)

	user, err := repo.FindByUsername("admin")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "admin", user.Username)
	assert.Equal(t, "hashed-password", user.Password)
	assert.Equal(t, "ADMIN", user.Role)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByUsername_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &UserRepositoryImpl{DB: db}

	mock.ExpectQuery(`
		SELECT id, username, password, role
		FROM users
		WHERE username = \$1
	`).WithArgs("unknown").WillReturnError(sql.ErrNoRows)

	user, err := repo.FindByUsername("unknown")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByUsername_DBError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &UserRepositoryImpl{DB: db}

	mock.ExpectQuery(`
		SELECT id, username, password, role
		FROM users
		WHERE username = \$1
	`).WithArgs("admin").WillReturnError(errors.New("db error"))

	user, err := repo.FindByUsername("admin")

	assert.Error(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Save(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &UserRepositoryImpl{DB: db}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`
			INSERT INTO users \(username, password, role\)
			VALUES \(\$1, \$2, \$3\)
		`).WithArgs("admin", "hashed-password", "ADMIN").
			WillReturnResult(sqlmock.NewResult(1, 1))

		user := &model.User{
			Username: "admin",
			Password: "hashed-password",
			Role:     "ADMIN",
		}

		err := repo.Save(user)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		mock.ExpectExec(`
			INSERT INTO users \(username, password, role\)
			VALUES \(\$1, \$2, \$3\)
		`).WithArgs("admin", "hashed-password", "ADMIN").
			WillReturnError(errors.New("insert failed"))

		user := &model.User{
			Username: "admin",
			Password: "hashed-password",
			Role:     "ADMIN",
		}

		err := repo.Save(user)
		assert.Error(t, err)
		assert.Equal(t, "insert failed", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// ================= tokenRepository =================

func TestTokenRepository_Save(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := &tokenRepository{DB: db}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(`
			INSERT INTO refresh_tokens \(user_id, token_hash, expires_at\)
			VALUES \(\$1, \$2, \$3\)
		`).WithArgs(1, "token-hash", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		token := &model.RefreshToken{
			UserID:    1,
			TokenHash: "token-hash",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		err := repo.Save(token)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db_error", func(t *testing.T) {
		mock.ExpectExec(`
			INSERT INTO refresh_tokens \(user_id, token_hash, expires_at\)
			VALUES \(\$1, \$2, \$3\)
		`).WithArgs(1, "token-hash", sqlmock.AnyArg()).
			WillReturnError(errors.New("insert failed"))

		token := &model.RefreshToken{
			UserID:    1,
			TokenHash: "token-hash",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		err := repo.Save(token)
		assert.Error(t, err)
		assert.Equal(t, "insert failed", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
