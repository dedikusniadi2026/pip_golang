package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken_Success(t *testing.T) {
	userID := int64(1)
	role := "admin"

	tokenString, err := GenerateAccessToken(userID, role)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTclaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(*JWTclaims)
	assert.True(t, ok)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)

	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
}

func TestGenerateAccessToken_ExpiryTime(t *testing.T) {
	start := time.Now()

	tokenString, err := GenerateAccessToken(99, "user")
	assert.NoError(t, err)

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTclaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	assert.NoError(t, err)

	claims := token.Claims.(*JWTclaims)

	diff := claims.ExpiresAt.Time.Sub(start)
	assert.True(t, diff > 14*time.Minute && diff < 16*time.Minute)
}
