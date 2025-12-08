package handler

import (
	"auth-service/model"
	"auth-service/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TripHandler struct {
	service service.TripServiceInterface
}

func NewTripHandler(tripService service.TripServiceInterface) *TripHandler {
	return &TripHandler{service: tripService}
}

func (h *TripHandler) Create(c *gin.Context) {
	var t model.VehicleTrip
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *TripHandler) FindByVehicle(c *gin.Context) {
	vehicleIDStr := c.Param("vehicle_id")
	vehicleID, err := strconv.Atoi(vehicleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle_id"})
		return
	}
	data, err := h.service.FindByVehicle(uint(vehicleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *TripHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var t model.VehicleTrip
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.ID = uint(id)
	if err := h.service.Update(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *TripHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *TripHandler) GetTripTotal(c *gin.Context) {
	stats, err := h.service.GetTripTotals(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_trips":    stats.TotalTrips,
		"total_distance": stats.TotalDistance,
		"total_revenue":  stats.TotalRevenue,
		"average_Rating": stats.AverageRating,
	})
}
