package handler

import (
	"auth-service/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingTrendsServiceInterface interface {
	GetTrends(year int) ([]model.BookingTrend, error)
}

type BookingTrendsHandler struct {
	Service BookingTrendsServiceInterface
}

func NewBookingTrendsHandler(s BookingTrendsServiceInterface) *BookingTrendsHandler {
	return &BookingTrendsHandler{Service: s}
}

func (h *BookingTrendsHandler) GetTrends(c *gin.Context) {
	yearStr := c.DefaultQuery("year", "2024")
	year, _ := strconv.Atoi(yearStr)

	data, err := h.Service.GetTrends(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var months []string
	var bookings []int

	for _, d := range data {
		months = append(months, d.Month)
		bookings = append(bookings, d.Booking)
	}

	c.JSON(http.StatusOK, gin.H{
		"months":   months,
		"bookings": bookings,
	})
}
