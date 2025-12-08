package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserWithRole(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := &service.AuthService{Repo: mockRepo}

	hashedPassword, _ := service.HashPassword("password123")
	user := &model.User{ID: 1, Username: "test@example.com", Password: hashedPassword, Role: "admin"}
	mockRepo.On("GetUserWithRole", "test@example.com").Return(user, nil)

	token, refreshToken, err := svc.Login("test@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refreshToken)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := &service.AuthService{Repo: mockRepo}

	hashedPassword, _ := service.HashPassword("password123")
	user := &model.User{ID: 1, Username: "test@example.com", Password: hashedPassword, Role: "admin"}
	mockRepo.On("GetUserWithRole", "test@example.com").Return(user, nil)

	token, refreshToken, err := svc.Login("test@example.com", "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	assert.Empty(t, token)
	assert.Empty(t, refreshToken)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := &service.AuthService{Repo: mockRepo}

	mockRepo.On("GetUserWithRole", "test@example.com").Return((*model.User)(nil), errors.New("user not found"))

	token, refreshToken, err := svc.Login("test@example.com", "password123")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Empty(t, token)
	assert.Empty(t, refreshToken)

	mockRepo.AssertExpectations(t)
}

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashed, err := service.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestCheckPasswordHash_Valid(t *testing.T) {
	password := "password123"
	hashed, _ := service.HashPassword(password)
	valid := service.CheckPasswordHash(password, hashed)
	assert.True(t, valid)
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
	password := "password123"
	hashed, _ := service.HashPassword(password)
	valid := service.CheckPasswordHash("wrongpassword", hashed)
	assert.False(t, valid)
}

func TestGenerateJWT(t *testing.T) {
	user := model.User{ID: 1, Username: "test@example.com", Password: "hashed", Role: "admin"}
	token, err := service.GenerateJWT(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(1), claims["user_id"])
}

func TestGenerateRefreshToken(t *testing.T) {
	user := model.User{ID: 1, Username: "test@example.com", Password: "hashed", Role: "admin"}
	token, err := service.GenerateRefreshToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("refresh_secret_key"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(1), claims["user_id"])
}
