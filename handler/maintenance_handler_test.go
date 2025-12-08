package handler_test

import (
	"auth-service/handler"
	"auth-service/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockMaintenanceService struct{}

func (m *MockMaintenanceService) Create(maintenance *model.VehicleMaintenance) error {
	if maintenance.VehicleID == 0 {
		return errors.New("vehicle id required")
	}
	maintenance.ID = 1
	return nil
}

func (m *MockMaintenanceService) FindByVehicle(vehicleID int) ([]model.VehicleMaintenance, error) {
	if vehicleID == 0 {
		return nil, errors.New("invalid vehicle id")
	}
	return []model.VehicleMaintenance{
		{
			ID:          1,
			VehicleID:   uint(vehicleID),
			ServiceDate: time.Now(),
			ServiceType: "Oil Change",
			Description: "Regular maintenance",
			Mileage:     10000,
			Cost:        50.0,
		},
	}, nil
}

func (m *MockMaintenanceService) Update(maintenance *model.VehicleMaintenance) error {
	if maintenance.ID == 0 {
		return errors.New("id required")
	}
	return nil
}

func (m *MockMaintenanceService) Delete(id uint) error {
	if id == 0 {
		return errors.New("id required")
	}
	return nil
}

func setupMaintenanceRouter(h *handler.MaintenanceHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/maintenance", h.Create)
	r.GET("/maintenance/vehicle/:vehicle_id", h.FindByVehicle)
	r.PUT("/maintenance/:id", h.Update)
	r.DELETE("/maintenance/:id", h.Delete)

	return r
}

func TestCreateMaintenance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	payload := model.VehicleMaintenance{
		VehicleID:   1,
		ServiceDate: time.Now(),
		ServiceType: "Oil Change",
		Description: "Regular maintenance",
		Mileage:     10000,
		Cost:        50.0,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/maintenance", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.VehicleMaintenance
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, uint(1), response.ID)
}

func TestCreateMaintenance_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	req, _ := http.NewRequest("POST", "/maintenance", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateMaintenance_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	payload := model.VehicleMaintenance{
		VehicleID: 0,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/maintenance", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestFindByVehicle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	req, _ := http.NewRequest("GET", "/maintenance/vehicle/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.VehicleMaintenance
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 1)
	assert.Equal(t, uint(1), response[0].VehicleID)
}

func TestFindByVehicle_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	req, _ := http.NewRequest("GET", "/maintenance/vehicle/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFindByVehicle_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	req, _ := http.NewRequest("GET", "/maintenance/vehicle/0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateMaintenance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	payload := model.VehicleMaintenance{
		VehicleID:   1,
		ServiceDate: time.Now(),
		ServiceType: "Brake Check",
		Description: "Updated maintenance",
		Mileage:     15000,
		Cost:        75.0,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/maintenance/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateMaintenance_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	payload := model.VehicleMaintenance{}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/maintenance/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateMaintenance_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	req, _ := http.NewRequest("PUT", "/maintenance/1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateMaintenance_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	payload := model.VehicleMaintenance{
		VehicleID:   1,
		ServiceDate: time.Now(),
		ServiceType: "Oil Change",
		Description: "Regular maintenance",
		Mileage:     10000,
		Cost:        50.0,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/maintenance/0", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteMaintenance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	req, _ := http.NewRequest("DELETE", "/maintenance/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "deleted", response["message"])
}

func TestDeleteMaintenance_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	req, _ := http.NewRequest("DELETE", "/maintenance/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteMaintenance_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockMaintenanceService{}
	h := handler.NewMaintenanceHandler(mockSvc)

	router := setupMaintenanceRouter(h)

	req, _ := http.NewRequest("DELETE", "/maintenance/0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
