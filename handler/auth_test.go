package handler_test

import (
	"auth-service/handler"
	"auth-service/service/mock_auth_service"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{
		LoginFn: func(username, password string) (string, string, error) {
			return "access_token_123", "refresh_token_456", nil
		},
	}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/login", h.Login)

	payload := `{"username":"admin","password":"password123"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "access_token_123")
	assert.Contains(t, w.Body.String(), "refresh_token_456")
}

func TestAuthHandler_Login_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/login", h.Login)

	payload := `{"username":"admin"}` // missing password
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/login", h.Login)

	payload := `{"username":"admin","password":}` // invalid JSON
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestAuthHandler_Login_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{
		LoginFn: func(username, password string) (string, string, error) {
			return "", "", errors.New("invalid credentials")
		},
	}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/login", h.Login)

	payload := `{"username":"admin","password":"wrongpassword"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
}

func TestAuthHandler_Refresh_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{
		RefreshTokenFn: func(refreshToken string) (string, error) {
			return "new_access_token_789", nil
		},
	}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/refresh", h.Refresh)

	payload := `{"refresh_token":"refresh_token_456"}`
	req, _ := http.NewRequest("POST", "/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "new_access_token_789")
}

func TestAuthHandler_Refresh_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/refresh", h.Refresh)

	payload := `{"refresh_token":}`
	req, _ := http.NewRequest("POST", "/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestAuthHandler_Refresh_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/refresh", h.Refresh)

	payload := `{}`
	req, _ := http.NewRequest("POST", "/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestAuthHandler_Refresh_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{
		RefreshTokenFn: func(refreshToken string) (string, error) {
			return "", errors.New("invalid refresh token")
		},
	}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/refresh", h.Refresh)

	payload := `{"refresh_token":"invalid_token"}`
	req, _ := http.NewRequest("POST", "/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid refresh token")
}

func TestAuthHandler_Register_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{
		RegisterFn: func(username, password, role string) error {
			return nil
		},
	}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/register", h.Register)

	payload := `{"username":"dedi","password":"12345","role":"user"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}

func TestAuthHandler_Register_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/register", h.Register)

	payload := `{"username":"dedi"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestAuthHandler_Login_Unauthorized_Coverage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{
		LoginFn: func(username, password string) (string, string, error) {
			return "", "", errors.New("invalid credentials")
		},
	}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/login", h.Login)

	payload := `{"username":"admin","password":"wrongpassword"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
}

func TestAuthHandler_Register_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mock_auth_service.MockAuthService{
		RegisterFn: func(username, password, role string) error {
			return errors.New("service error")
		},
	}
	h := handler.NewAuthHandler(mockSvc)

	router := gin.Default()
	router.POST("/register", h.Register)

	payload := `{"username":"dedi","password":"12345","role":"user"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "service error")
}
