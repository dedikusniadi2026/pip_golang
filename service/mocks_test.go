package service

import (
	"errors"
	"testing"
	"time"

	"auth-service/model"

	"github.com/stretchr/testify/assert"
)

func TestMockUserRepository_FindByUsername(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindByUsernameFn: func(username string) (*model.User, error) {
			assert.Equal(t, "admin", username)
			return &model.User{
				ID:       1,
				Username: "admin",
				Password: "hashed-pass",
				Role:     "ADMIN",
			}, nil
		},
	}

	user, err := mockRepo.FindByUsername("admin")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "admin", user.Username)
	assert.Equal(t, "hashed-pass", user.Password)
	assert.Equal(t, "ADMIN", user.Role)
}

func TestMockUserRepository_FindByUsername_Error(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindByUsernameFn: func(username string) (*model.User, error) {
			return nil, errors.New("user not found")
		},
	}

	user, err := mockRepo.FindByUsername("unknown")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())
}

func TestMockTokenRepository_Save(t *testing.T) {
	mockRepo := &MockTokenRepository{
		SaveFn: func(token *model.RefreshToken) error {
			assert.Equal(t, int64(1), token.UserID)
			assert.Equal(t, "token-hash", token.TokenHash)
			assert.WithinDuration(t, time.Now().Add(time.Hour), token.ExpiresAt, time.Minute*1)
			return nil
		},
	}

	err := mockRepo.Save(&model.RefreshToken{
		UserID:    1,
		TokenHash: "token-hash",
		ExpiresAt: time.Now().Add(time.Hour),
	})

	assert.NoError(t, err)
}

func TestMockTokenRepository_Save_Error(t *testing.T) {
	mockRepo := &MockTokenRepository{
		SaveFn: func(token *model.RefreshToken) error {
			return errors.New("insert failed")
		},
	}

	err := mockRepo.Save(&model.RefreshToken{
		UserID:    1,
		TokenHash: "token-hash",
		ExpiresAt: time.Now().Add(time.Hour),
	})

	assert.Error(t, err)
	assert.Equal(t, "insert failed", err.Error())
}
