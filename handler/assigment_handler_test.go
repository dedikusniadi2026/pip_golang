package handler_test

import (
	"auth-service/handler"
	"auth-service/model"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAssignmentsService struct {
	mock.Mock
}

func (m *MockAssignmentsService) Create(a *model.DriverAssignment) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m *MockAssignmentsService) FindByVehicle(vehicleID uint) ([]model.DriverAssignment, error) {
	args := m.Called(vehicleID)
	return args.Get(0).([]model.DriverAssignment), args.Error(1)
}

func (m *MockAssignmentsService) Update(a *model.DriverAssignment) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m *MockAssignmentsService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAssignmentHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	assignment := &model.DriverAssignment{
		VehicleID:  1,
		DriverID:   sql.NullInt64{Int64: 1, Valid: true},
		DriverName: "John Doe",
		StartDate:  time.Now(),
		EndDate:    time.Now().Add(24 * time.Hour),
		TotalTrips: 5,
		Status:     "active",
	}
	mockService.On("Create", mock.Anything).Return(nil)

	body, _ := json.Marshal(assignment)
	req, _ := http.NewRequest("POST", "/assignments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestAssignmentHandler_Create_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	assignment := &model.DriverAssignment{
		VehicleID: 1,
	}
	mockService.On("Create", assignment).Return(errors.New("create error"))

	body, _ := json.Marshal(assignment)
	req, _ := http.NewRequest("POST", "/assignments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestAssignmentHandler_Create_InvalidJSON(t *testing.T) {
	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	req, _ := http.NewRequest("POST", "/assignments", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAssignmentHandler_FindByVehicle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	assignments := []model.DriverAssignment{
		{
			ID:         1,
			VehicleID:  1,
			DriverName: "John Doe",
			Status:     "active",
		},
	}
	mockService.On("FindByVehicle", uint(1)).Return(assignments, nil)

	req, _ := http.NewRequest("GET", "/assignments/vehicle/1", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "vehicle_id", Value: "1"}}

	h.FindByVehicle(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAssignmentHandler_FindByVehicle_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	mockService.On("FindByVehicle", uint(1)).Return([]model.DriverAssignment{}, errors.New("find error"))

	req, _ := http.NewRequest("GET", "/assignments/vehicle/1", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "vehicle_id", Value: "1"}}

	h.FindByVehicle(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestAssignmentHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	assignment := &model.DriverAssignment{
		ID:         1,
		VehicleID:  1,
		DriverName: "Jane Doe",
		Status:     "inactive",
	}
	mockService.On("Update", mock.Anything).Return(nil)

	body, _ := json.Marshal(assignment)
	req, _ := http.NewRequest("PUT", "/assignments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAssignmentHandler_Update_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	assignment := &model.DriverAssignment{
		ID: 1,
	}
	mockService.On("Update", assignment).Return(errors.New("update error"))

	body, _ := json.Marshal(assignment)
	req, _ := http.NewRequest("PUT", "/assignments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestAssignmentHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/assignments/1", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAssignmentHandler_Delete_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	mockService.On("Delete", uint(1)).Return(errors.New("delete error"))

	req, _ := http.NewRequest("DELETE", "/assignments/1", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Delete(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestAssignmentHandler_FindByVehicle_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	req, _ := http.NewRequest("GET", "/assignments/vehicle/abc", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "vehicle_id", Value: "abc"}}

	h.FindByVehicle(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAssignmentHandler_Delete_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	mockService.On("Delete", uint(0)).Return(errors.New("delete error"))

	req, _ := http.NewRequest("DELETE", "/assignments/abc", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.Delete(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestAssignmentHandler_Update_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockAssignmentsService)
	h := handler.NewAssignmentHandler(mockService)

	req, _ := http.NewRequest("PUT", "/assignments", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
