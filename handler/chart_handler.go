package handler

import (
	"auth-service/service"

	"github.com/gin-gonic/gin"
)

type PopularDestinationHandler struct {
	Service service.PopularDestinationServiceInterface
}

func (h *PopularDestinationHandler) GetAll(c *gin.Context) {
	data, err := h.Service.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var resp []map[string]interface{}
	for _, d := range data {
		resp = append(resp, map[string]interface{}{
			"destination": d.Destination,
			"bookings":    d.Bookings,
		})
	}

	c.JSON(200, resp)
}
