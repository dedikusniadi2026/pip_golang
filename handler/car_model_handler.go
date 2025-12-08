package handler

import (
	"auth-service/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CarModelServiceInterface interface {
	GetAll() ([]model.CarModel, error)
	GetByID(int) (*model.CarModel, error)
	Create(model.CarModel) error
}

type CarModelHandler struct {
	Service CarModelServiceInterface
}

func NewCarModelRepositoryHandler(service CarModelServiceInterface) *CarModelHandler {
	return &CarModelHandler{Service: service}
}

func (h *CarModelHandler) GetAll(c *gin.Context) {
	data, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *CarModelHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	carModel, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if carModel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CarModel not found"})
		return
	}

	c.JSON(http.StatusOK, carModel)
}

func (h *CarModelHandler) Create(c *gin.Context) {
	var body model.CarModel
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.Create(body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}
