package service

import (
	"auth-service/model"
	"auth-service/utils"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockUserRepo struct {
	FindByUsernameFn func(username string) (*model.User, error)
	SaveFn           func(user *model.User) error
}

func (m *MockUserRepo) FindByUsername(username string) (*model.User, error) {
	if m.FindByUsernameFn != nil {
		return m.FindByUsernameFn(username)
	}
	return nil, nil
}

func (m *MockUserRepo) Save(user *model.User) error {
	if m.SaveFn != nil {
		return m.SaveFn(user)
	}
	return nil
}

type MockTokenRepo struct {
	SaveFn func(token *model.RefreshToken) error
}

func (m *MockTokenRepo) Save(token *model.RefreshToken) error {
	if m.SaveFn != nil {
		return m.SaveFn(token)
	}
	return nil
}

type MockTokenRepoError struct{}

func (m *MockTokenRepoError) Save(token *model.RefreshToken) error {
	return errors.New("db error")
}

func TestAuthService_Login_Success(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password123")

	userRepo := &MockUserRepo{
		FindByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				ID:       1,
				Username: "admin",
				Password: hashedPassword,
				Role:     "ADMIN",
			}, nil
		},
	}

	tokenRepo := &MockTokenRepo{
		SaveFn: func(token *model.RefreshToken) error {
			assert.Equal(t, int64(1), token.UserID)
			assert.NotEmpty(t, token.TokenHash)
			assert.True(t, token.ExpiresAt.After(time.Now()))
			return nil
		},
	}

	service := NewAuthService(userRepo, tokenRepo)

	accessToken, refreshToken, err := service.Login("admin", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
}

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := &MockUserRepo{
		SaveFn: func(user *model.User) error {
			assert.Equal(t, "admin", user.Username)
			assert.NotEmpty(t, user.Password)
			assert.Equal(t, "ADMIN", user.Role)
			return nil
		},
	}

	service := &AuthService{
		UserRepo:       userRepo,
		HashPasswordFn: utils.HashPassword,
	}

	err := service.Register("admin", "password123", "ADMIN")
	assert.NoError(t, err)
}

func TestAuthService_Register_HashPasswordError(t *testing.T) {
	userRepo := &MockUserRepo{}

	service := &AuthService{
		UserRepo: userRepo,
		HashPasswordFn: func(password string) (string, error) {
			return "", errors.New("hash error")
		},
	}

	err := service.Register("admin", "password123", "ADMIN")
	assert.EqualError(t, err, "hash error")
}

func TestAuthService_Register_SaveUserError(t *testing.T) {
	userRepo := &MockUserRepo{
		SaveFn: func(user *model.User) error {
			return errors.New("db error")
		},
	}

	service := &AuthService{
		UserRepo:       userRepo,
		HashPasswordFn: utils.HashPassword,
	}

	err := service.Register("admin", "password123", "ADMIN")
	assert.EqualError(t, err, "db error")
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	userRepo := &MockUserRepo{
		FindByUsernameFn: func(username string) (*model.User, error) {
			return nil, errors.New("user not found")
		},
	}
	tokenRepo := &MockTokenRepo{}

	service := NewAuthService(userRepo, tokenRepo)
	access, refresh, err := service.Login("admin", "password123")

	assert.Empty(t, access)
	assert.Empty(t, refresh)
	assert.EqualError(t, err, "invalid credentials")
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("correct_password")

	userRepo := &MockUserRepo{
		FindByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				ID:       1,
				Username: "admin",
				Password: hashedPassword,
				Role:     "ADMIN",
			}, nil
		},
	}
	tokenRepo := &MockTokenRepo{}

	service := NewAuthService(userRepo, tokenRepo)
	access, refresh, err := service.Login("admin", "wrong_password")

	assert.Empty(t, access)
	assert.Empty(t, refresh)
	assert.EqualError(t, err, "invalid credentials password")
}

func TestAuthService_Login_GenerateAccessTokenError(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password123")

	userRepo := &MockUserRepo{
		FindByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				ID:       1,
				Username: "admin",
				Password: hashedPassword,
				Role:     "ADMIN",
			}, nil
		},
	}
	tokenRepo := &MockTokenRepo{}

	service := NewAuthService(userRepo, tokenRepo)
	service.GenerateAccessToken = func(userID int64, role string) (string, error) {
		return "", errors.New("token generation error")
	}

	access, refresh, err := service.Login("admin", "password123")

	assert.Empty(t, access)
	assert.Empty(t, refresh)
	assert.EqualError(t, err, "token generation error")
}

func TestAuthService_Login_SaveRefreshTokenError(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password123")

	userRepo := &MockUserRepo{
		FindByUsernameFn: func(username string) (*model.User, error) {
			return &model.User{
				ID:       1,
				Username: "admin",
				Password: hashedPassword,
				Role:     "ADMIN",
			}, nil
		},
	}
	tokenRepo := &MockTokenRepoError{}

	service := NewAuthService(userRepo, tokenRepo)

	access, refresh, err := service.Login("admin", "password123")

	assert.Empty(t, access)
	assert.Empty(t, refresh)
	assert.EqualError(t, err, "db error")
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	service := &AuthService{}

	access, err := service.RefreshToken("any_refresh_token")

	assert.NoError(t, err)
	assert.Equal(t, "NEW_ACCESS_TOKEN", access)
}
