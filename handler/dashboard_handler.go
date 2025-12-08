package handler

import (
	"auth-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service service.DashboardServiceInterface
}

func NewDashboardHandler(s service.DashboardServiceInterface) *DashboardHandler {
	return &DashboardHandler{service: s}
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	data, err := h.service.GetDashboardData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
