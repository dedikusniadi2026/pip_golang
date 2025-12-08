package handler

import (
	"net/http"

	"auth-service/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}
