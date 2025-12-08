package repository_test

import (
	"auth-service/repository"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserWithRole_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.UserRepository{DB: db}

	username := "dedi"
	rows := sqlmock.NewRows([]string{"id", "username", "password", "role_name"}).
		AddRow(1, username, "hashedpassword", "admin")

	mock.ExpectQuery("SELECT u.id, u.username, u.password, r.role_name").
		WithArgs(username).
		WillReturnRows(rows)

	user, err := repo.GetUserWithRole(username)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, "admin", user.Role)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserWithRole_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.UserRepository{DB: db}

	username := "unknown"

	mock.ExpectQuery("SELECT u.id, u.username, u.password, r.role_name").
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserWithRole(username)

	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
