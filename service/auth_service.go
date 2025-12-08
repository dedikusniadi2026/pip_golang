package service

import (
	"auth-service/model"
	"auth-service/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")
var refreshKey = []byte("refresh_secret_key")

type AuthService struct {
	Repo repository.UserRepositoryInterface
}

func (a *AuthService) Login(email string, password string) (string, string, error) {
	user, err := a.Repo.GetUserWithRole(email)
	if err != nil {
		return "", "", err
	}

	if !CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	token, err := GenerateJWT(*user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateRefreshToken(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshKey)
}
