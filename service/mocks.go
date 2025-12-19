package service

import (
	"auth-service/model"
)

type MockUserRepository struct {
	FindByUsernameFn func(username string) (*model.User, error)
}

func (m *MockUserRepository) FindByUsername(username string) (*model.User, error) {
	return m.FindByUsernameFn(username)
}

type MockTokenRepository struct {
	SaveFn func(token *model.RefreshToken) error
}

func (m *MockTokenRepository) Save(token *model.RefreshToken) error {
	return m.SaveFn(token)
}
