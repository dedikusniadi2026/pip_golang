package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/handler"
	"auth-service/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(h *handler.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/login", h.Login)
	r.POST("/refresh", h.Refresh)

	return r
}

func TestLoginHandler(t *testing.T) {
	authService := &service.AuthService{}
	h := &handler.AuthHandler{AuthService: authService}

	router := setupRouter(h)

	req, _ := http.NewRequest("POST", "/login", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login successful")
}

func TestRefreshHandler(t *testing.T) {
	authService := &service.AuthService{}
	h := &handler.AuthHandler{AuthService: authService}

	router := setupRouter(h)

	req, _ := http.NewRequest("POST", "/refresh", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Token refreshed")
}
