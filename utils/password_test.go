package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {
	password := "super-secret-password"

	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPasswordHash_ValidPassword(t *testing.T) {
	password := "my-password"

	hash, err := HashPassword(password)
	assert.NoError(t, err)

	ok := CheckPasswordHash(password, hash)
	assert.True(t, ok)
}

func TestCheckPasswordHash_InvalidPassword(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"

	hash, err := HashPassword(password)
	assert.NoError(t, err)

	ok := CheckPasswordHash(wrongPassword, hash)
	assert.False(t, ok)
}

func TestCheckPassword_ValidPassword(t *testing.T) {
	password := "admin123"

	hash, err := HashPassword(password)
	assert.NoError(t, err)

	err = CheckPassword(hash, password)
	assert.NoError(t, err)
}

func TestCheckPassword_InvalidPassword(t *testing.T) {
	password := "admin123"
	wrongPassword := "admin321"

	hash, err := HashPassword(password)
	assert.NoError(t, err)

	err = CheckPassword(hash, wrongPassword)
	assert.Error(t, err)
}
