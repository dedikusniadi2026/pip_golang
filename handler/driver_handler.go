package handler

import (
	"auth-service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DriverServiceInterface interface {
	GetAll() ([]model.Driver, error)
	Create(*model.Driver) error
	GetByID(string) (*model.Driver, error)
	Update(string, *model.Driver) error
	Delete(string) error
}

type DriverHandler struct {
	Service DriverServiceInterface
}

func NewDriverHandler(service DriverServiceInterface) *DriverHandler {
	return &DriverHandler{Service: service}
}

func (h *DriverHandler) GetAll(c *gin.Context) {
	drivers, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, drivers)
}

func (h *DriverHandler) Create(c *gin.Context) {
	var req struct {
		Name                string `json:"name"`
		Email               string `json:"email"`
		Phone               string `json:"phone"`
		Address             string `json:"address"`
		DriverLicenseNumber string `json:"driver_license_number"`
		CarModelID          string `json:"car_model_id"`
		CarTypeID           string `json:"car_type_id"`
		PlateNumber         string `json:"plate_number"`
		Status              string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	driver := model.Driver{
		Name:                req.Name,
		Email:               req.Email,
		Phone:               req.Phone,
		Address:             req.Address,
		DriverLicenseNumber: req.DriverLicenseNumber,
		CarModelID:          req.CarModelID,
		CarTypeID:           req.CarTypeID,
		PlateNumber:         req.PlateNumber,
		Status:              req.Status,
	}

	if err := h.Service.Create(&driver); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, driver)
}

func (h *DriverHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	driver, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, driver)
}

func (h *DriverHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name                string `json:"name"`
		Email               string `json:"email"`
		Phone               string `json:"phone"`
		Address             string `json:"address"`
		DriverLicenseNumber string `json:"driver_license_number"`
		CarModelID          string `json:"car_model_id"`
		CarTypeID           string `json:"car_type_id"`
		PlateNumber         string `json:"plate_number"`
		Status              string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	driver := model.Driver{
		Name:                req.Name,
		Email:               req.Email,
		Phone:               req.Phone,
		Address:             req.Address,
		DriverLicenseNumber: req.DriverLicenseNumber,
		CarModelID:          req.CarModelID,
		CarTypeID:           req.CarTypeID,
		PlateNumber:         req.PlateNumber,
		Status:              req.Status,
	}

	if err := h.Service.Update(id, &driver); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Driver updated successfully"})
}

func (h *DriverHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.Service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Driver deleted successfully"})
}
