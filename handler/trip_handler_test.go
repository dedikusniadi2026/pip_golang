package handler_test

import (
	"auth-service/handler"
	"auth-service/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockTripService struct {
	ReturnError bool
}

func (m *MockTripService) Create(t *model.VehicleTrip) error {
	if m.ReturnError {
		return errors.New("mock error")
	}
	if t.VehicleID == 0 {
		return errors.New("vehicle id required")
	}
	t.ID = 1
	return nil
}

func (m *MockTripService) FindByVehicle(vehicleID uint) ([]model.VehicleTrip, error) {
	if m.ReturnError {
		return nil, errors.New("mock error")
	}
	return []model.VehicleTrip{
		{ID: 1, VehicleID: vehicleID, DriverID: 1, Origin: "A", Destination: "B", DistanceKM: 10, Rating: 5, Price: 100.0, PassengerName: "John"},
	}, nil
}

func (m *MockTripService) Update(t *model.VehicleTrip) error {
	if m.ReturnError {
		return errors.New("mock error")
	}
	if t.ID == 0 {
		return errors.New("id required")
	}
	return nil
}

func (m *MockTripService) Delete(id uint) error {
	if m.ReturnError {
		return errors.New("mock error")
	}
	if id == 0 {
		return errors.New("id required")
	}
	return nil
}

func (m *MockTripService) GetTripTotals(ctx context.Context) (*model.TotalTrips, error) {
	if m.ReturnError {
		return nil, errors.New("service error")
	}
	return &model.TotalTrips{
		TotalTrips:    10,
		TotalDistance: 100,
		TotalRevenue:  1000.0,
		AverageRating: 4.5,
	}, nil
}

func TestCreateTrip(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.POST("/trips", h.Create)

	payload := model.VehicleTrip{VehicleID: 1, DriverID: 1, Origin: "A", Destination: "B", DistanceKM: 10, Rating: 5, Price: 100.0, PassengerName: "John"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/trips", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestTripFindByVehicle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.GET("/trips/vehicle/:vehicle_id", h.FindByVehicle)

	req, _ := http.NewRequest("GET", "/trips/vehicle/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateTrip(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.PUT("/trips/:id", h.Update)

	payload := model.VehicleTrip{ID: 1, VehicleID: 1, DriverID: 1, Origin: "A", Destination: "B", DistanceKM: 10, Rating: 5, Price: 100.0, PassengerName: "John"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/trips/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTrip(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.DELETE("/trips/:id", h.Delete)

	req, _ := http.NewRequest("DELETE", "/trips/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTripTotal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.GET("/trips/total", h.GetTripTotal)

	req, _ := http.NewRequest("GET", "/trips/total", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTripHandler_FindByVehicle_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.GET("/trips/vehicle/:vehicle_id", h.FindByVehicle)

	req, _ := http.NewRequest("GET", "/trips/vehicle/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTripHandler_Update_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.PUT("/trips/:id", h.Update)

	payload := model.VehicleTrip{ID: 1, VehicleID: 1, DriverID: 1, Origin: "A", Destination: "B", DistanceKM: 10, Rating: 5, Price: 100.0, PassengerName: "John"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/trips/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTripHandler_Update_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.PUT("/trips/:id", h.Update)

	req, _ := http.NewRequest("PUT", "/trips/1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTripHandler_Delete_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.DELETE("/trips/:id", h.Delete)

	req, _ := http.NewRequest("DELETE", "/trips/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTrip_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{ReturnError: true}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.POST("/trips", h.Create)

	payload := model.VehicleTrip{VehicleID: 1, DriverID: 1, Origin: "A", Destination: "B", DistanceKM: 10, Rating: 5, Price: 100.0, PassengerName: "John"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/trips", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTripFindByVehicle_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{ReturnError: true}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.GET("/trips/vehicle/:vehicle_id", h.FindByVehicle)

	req, _ := http.NewRequest("GET", "/trips/vehicle/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateTrip_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{ReturnError: true}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.PUT("/trips/:id", h.Update)

	payload := model.VehicleTrip{ID: 1, VehicleID: 1, DriverID: 1, Origin: "A", Destination: "B", DistanceKM: 10, Rating: 5, Price: 100.0, PassengerName: "John"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/trips/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteTrip_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{ReturnError: true}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.DELETE("/trips/:id", h.Delete)

	req, _ := http.NewRequest("DELETE", "/trips/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetTripTotal_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{ReturnError: true}
	h := handler.NewTripHandler(mockSvc)

	router := gin.Default()
	router.GET("/trips/total", h.GetTripTotal)

	req, _ := http.NewRequest("GET", "/trips/total", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTripHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockTripService{}
	h := handler.NewTripHandler(mockSvc)

	r := gin.New()
	r.POST("/trips", h.Create)

	req := httptest.NewRequest(
		http.MethodPost,
		"/trips",
		strings.NewReader("{invalid json"),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
