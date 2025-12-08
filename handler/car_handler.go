package handler

import (
	"net/http"
	"strconv"

	"auth-service/model"

	"github.com/gin-gonic/gin"
)

type CarServiceInterface interface {
	GetAll() ([]model.Car, error)
	GetByID(int) (*model.Car, error)
	Create(model.Car) error
	Update(int, model.Car) error
	Delete(int) error
}

type CarHandler struct {
	Service CarServiceInterface
}

func NewCarHandler(service CarServiceInterface) *CarHandler {
	return &CarHandler{Service: service}
}

func (h *CarHandler) GetAll(c *gin.Context) {
	data, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *CarHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	v, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if v == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *CarHandler) Create(c *gin.Context) {
	var v model.Car
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.Create(v); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Car created"})
}

func (h *CarHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var v model.Car
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.Update(id, v); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car updated"})
}

func (h *CarHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.Service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car deleted"})
}
