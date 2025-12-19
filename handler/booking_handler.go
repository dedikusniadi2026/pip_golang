package handler

import (
	"auth-service/model"
	"auth-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingRepoInterface interface {
	Create(*model.Booking) error
	GetAll() ([]model.Booking, error)
	Update(*model.Booking) error
	Delete(id string) error
}

type BookingHandler struct {
	BookingService service.BookingServiceInterface
}

func NewBookingHandler(s service.BookingServiceInterface) *BookingHandler {
	return &BookingHandler{BookingService: s}
}

func (h *BookingHandler) Create(c *gin.Context) {
	var booking model.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.BookingService.Create(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) GetAll(c *gin.Context) {
	bookings, err := h.BookingService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var booking model.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking.ID = id
	if err := h.BookingService.Update(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.BookingService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}
