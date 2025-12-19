package handler

import (
	"auth-service/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CarTypeServiceInterface interface {
	GetAll() ([]model.CarType, error)
	GetByID(id int) (*model.CarType, error)
	Create(model.CarType) error
}
type CarTypeHandler struct {
	Service CarTypeServiceInterface
}

func NewCarTypeHandler(s CarTypeServiceInterface) *CarTypeHandler {
	return &CarTypeHandler{Service: s}
}
func (h *CarTypeHandler) GetAll(c *gin.Context) {
	data, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *CarTypeHandler) GetByID(c *gin.Context) {
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

func (h *CarTypeHandler) Create(c *gin.Context) {
	var body model.CarType

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
