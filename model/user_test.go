package model_test

import (
	"auth-service/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAdmin(t *testing.T) {
	userAdmin := &model.User{
		ID:       1,
		Username: "dedi",
		Role:     "admin",
	}
	userNonAdmin := &model.User{
		ID:       2,
		Username: "guest",
		Role:     "user",
	}

	assert.True(t, userAdmin.IsAdmin(), "User with role admin should return true")
	assert.False(t, userNonAdmin.IsAdmin(), "User with role user should return false")
}
