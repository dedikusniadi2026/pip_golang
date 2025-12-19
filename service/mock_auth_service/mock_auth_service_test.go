package mock_auth_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockAuthService_Login_WithFn(t *testing.T) {
	mock := &MockAuthService{
		LoginFn: func(username, password string) (string, string, error) {
			assert.Equal(t, "admin", username)
			assert.Equal(t, "secret", password)
			return "access-token", "refresh-token", nil
		},
	}

	access, refresh, err := mock.Login("admin", "secret")

	assert.NoError(t, err)
	assert.Equal(t, "access-token", access)
	assert.Equal(t, "refresh-token", refresh)
}

func TestMockAuthService_Login_WithoutFn(t *testing.T) {
	mock := &MockAuthService{}

	access, refresh, err := mock.Login("admin", "secret")

	assert.NoError(t, err)
	assert.Equal(t, "", access)
	assert.Equal(t, "", refresh)
}

func TestMockAuthService_RefreshToken_WithFn(t *testing.T) {
	mock := &MockAuthService{
		RefreshTokenFn: func(refreshToken string) (string, error) {
			assert.Equal(t, "refresh-token", refreshToken)
			return "new-access-token", nil
		},
	}

	token, err := mock.RefreshToken("refresh-token")

	assert.NoError(t, err)
	assert.Equal(t, "new-access-token", token)
}

func TestMockAuthService_RefreshToken_WithError(t *testing.T) {
	mock := &MockAuthService{
		RefreshTokenFn: func(refreshToken string) (string, error) {
			return "", errors.New("invalid refresh token")
		},
	}

	token, err := mock.RefreshToken("bad-token")

	assert.Error(t, err)
	assert.Equal(t, "", token)
}

func TestMockAuthService_Register_WithFn(t *testing.T) {
	mock := &MockAuthService{
		RegisterFn: func(username, password, role string) error {
			assert.Equal(t, "user1", username)
			assert.Equal(t, "pass123", password)
			assert.Equal(t, "admin", role)
			return nil
		},
	}

	err := mock.Register("user1", "pass123", "admin")
	assert.NoError(t, err)
}

func TestMockAuthService_Register_WithError(t *testing.T) {
	mock := &MockAuthService{
		RegisterFn: func(username, password, role string) error {
			return errors.New("username already exists")
		},
	}

	err := mock.Register("user1", "pass123", "admin")
	assert.Error(t, err)
	assert.Equal(t, "username already exists", err.Error())
}

func TestMockAuthService_Register_WithoutFn(t *testing.T) {
	mock := &MockAuthService{}

	err := mock.Register("user1", "pass123", "admin")
	assert.NoError(t, err) // default behavior returns nil
}

func TestMockAuthService_RefreshToken_WithoutFn(t *testing.T) {
	mock := &MockAuthService{}

	token, err := mock.RefreshToken("refresh-token")

	assert.NoError(t, err)
	assert.Equal(t, "", token)
}
