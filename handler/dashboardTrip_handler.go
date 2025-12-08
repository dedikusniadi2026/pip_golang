package handler

import (
	"auth-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardTripHandler struct {
	service service.DashboardTripServiceInterface
}

func NewDashboardTripHandler(service service.DashboardTripServiceInterface) *DashboardTripHandler {
	return &DashboardTripHandler{service: service}
}

func (h *DashboardTripHandler) GetDashboard(c *gin.Context) {
	data, err := h.service.GetDashboardSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
