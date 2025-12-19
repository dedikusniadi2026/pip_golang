package service

import (
	"auth-service/model"
	"auth-service/repository"
	"auth-service/utils"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"time"
)

type AuthServiceInterface interface {
	Register(username, password, role string) error
	Login(username, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
}

type AuthService struct {
	UserRepo            repository.UserRepository
	TokenRepo           repository.TokenRepository
	GenerateAccessToken func(userID int64, role string) (string, error)
	HashPasswordFn      func(string) (string, error)
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
) *AuthService {
	return &AuthService{
		UserRepo:            userRepo,
		TokenRepo:           tokenRepo,
		GenerateAccessToken: utils.GenerateAccessToken,
		HashPasswordFn:      utils.HashPassword,
	}
}

func (s *AuthService) Register(username, password, role string) error {
	hashedPassword, err := s.HashPasswordFn(password) // <- pakai yang di-inject
	if err != nil {
		return err
	}

	user := &model.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	return s.UserRepo.Save(user)
}

func (s *AuthService) Login(username, password string) (string, string, error) {
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	log.Println("DB PASSWORD HASH:", user.Password)
	log.Println("INPUT PASSWORD:", password)

	err = utils.CheckPassword(user.Password, password)
	if err != nil {
		log.Println("CHECK PASSWORD ERROR:", err)
		return "", "", errors.New("invalid credentials password")
	}

	accessToken, err := s.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken := generateRandomToken()
	refreshTokenHash := hashToken(refreshToken)

	err = s.TokenRepo.Save(&model.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	return "NEW_ACCESS_TOKEN", nil
}

func generateRandomToken() string {
	return hex.EncodeToString([]byte(time.Now().String()))
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
