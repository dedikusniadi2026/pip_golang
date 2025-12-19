package repository

import "auth-service/model"

type UserRepository interface {
	FindByUsername(username string) (*model.User, error)
	Save(user *model.User) error
}

type TokenRepository interface {
	Save(token *model.RefreshToken) error
}
