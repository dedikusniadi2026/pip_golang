package handler

import (
	"auth-service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TripHistoryServiceInterface interface {
	GetTripHistory() ([]model.TripHistory, error)
}

type TripHistoryHandler struct {
	Service TripHistoryServiceInterface
}

func NewTripHistoryHandler(s TripHistoryServiceInterface) *TripHistoryHandler {
	return &TripHistoryHandler{Service: s}
}

func (h *TripHistoryHandler) GetTripHistory(c *gin.Context) {
	result, err := h.Service.GetTripHistory()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trip not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}
