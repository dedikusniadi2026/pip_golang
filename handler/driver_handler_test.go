package handler_test

import (
	"auth-service/handler"
	"auth-service/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockDriverService struct {
	ReturnError bool
}

func (m *MockDriverService) GetAll() ([]model.Driver, error) {
	if m.ReturnError {
		return nil, errors.New("failed to fetch drivers")
	}
	return []model.Driver{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Phone: "1234567890", Address: "123 Main St", DriverLicenseNumber: "DL123", CarModelID: "1", CarTypeID: "1", PlateNumber: "ABC123", Status: "active"},
	}, nil
}

func (m *MockDriverService) Create(driver *model.Driver) error {
	if m.ReturnError {
		return errors.New("insert failed")
	}
	if driver.Name == "" {
		return errors.New("name required")
	}
	driver.ID = 1
	return nil
}

func (m *MockDriverService) GetByID(id string) (*model.Driver, error) {
	if id == "1" {
		return &model.Driver{ID: 1, Name: "John Doe", Email: "john@example.com", Phone: "1234567890", Address: "123 Main St", DriverLicenseNumber: "DL123", CarModelID: "1", CarTypeID: "1", PlateNumber: "ABC123", Status: "active"}, nil
	}
	return nil, errors.New("driver not found")
}

func (m *MockDriverService) Update(id string, driver *model.Driver) error {
	if id == "1" {
		return nil
	}
	return errors.New("driver not found")
}

func (m *MockDriverService) Delete(id string) error {
	if id == "1" {
		return nil
	}
	return errors.New("driver not found")
}

func setupDriverRouter(h *handler.DriverHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/drivers", h.GetAll)
	r.POST("/drivers", h.Create)
	r.GET("/drivers/:id", h.GetByID)
	r.PUT("/drivers/:id", h.Update)
	r.DELETE("/drivers/:id", h.Delete)

	return r
}

func TestGetAllDrivers(t *testing.T) {
	mockSvc := &MockDriverService{}
	h := handler.NewDriverHandler(mockSvc)

	router := setupDriverRouter(h)

	req, _ := http.NewRequest("GET", "/drivers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []model.Driver
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 1)
	assert.Equal(t, "John Doe", response[0].Name)
}

func TestGetAllDrivers_Error(t *testing.T) {
	mockSvc := &MockDriverService{}
	mockSvc.ReturnError = true
	h := handler.NewDriverHandler(mockSvc)

	router := setupDriverRouter(h)

	req, _ := http.NewRequest("GET", "/drivers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDriverHandler_Create_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockDriverService{
		ReturnError: true,
	}

	h := &handler.DriverHandler{
		Service: mockSvc,
	}

	r := gin.New()
	r.POST("/drivers", h.Create)

	payload := `{
        "name":"John",
        "email":"john@mail.com",
        "phone":"08123",
        "address":"Bandung",
        "driver_license_number":"SIM123",
        "car_model_id":"1",
        "car_type_id":"1",
        "plate_number":"B1234CD",
        "status":"active"
    }`

	req := httptest.NewRequest(
		http.MethodPost,
		"/drivers",
		strings.NewReader(payload),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "insert failed")
}

func TestDriverHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockDriverService{
		ReturnError: false,
	}

	h := &handler.DriverHandler{
		Service: mockSvc,
	}

	r := gin.New()
	r.POST("/drivers", h.Create)

	payload := `{
		"name":"John",
		"email":"john@mail.com",
		"phone":"08123",
		"address":"Bandung",
		"driver_license_number":"SIM123",
		"car_model_id":"1",
		"car_type_id":"1",
		"plate_number":"B1234CD",
		"status":"active"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/drivers",
		strings.NewReader(payload),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDriverHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockDriverService{}
	h := &handler.DriverHandler{Service: mockSvc}

	r := gin.New()
	r.POST("/drivers", h.Create)

	req := httptest.NewRequest(
		http.MethodPost,
		"/drivers",
		strings.NewReader("invalid json"),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateDriver(t *testing.T) {
	mockSvc := &MockDriverService{}
	h := handler.NewDriverHandler(mockSvc)

	router := setupDriverRouter(h)

	payload := map[string]string{
		"name":                  "Jane Doe",
		"email":                 "jane@example.com",
		"phone":                 "0987654321",
		"address":               "456 Elm St",
		"driver_license_number": "DL456",
		"car_model_id":          "2",
		"car_type_id":           "2",
		"plate_number":          "XYZ789",
		"status":                "active",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/drivers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response model.Driver
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Jane Doe", response.Name)
}

func TestCreateDriverInvalidJSON(t *testing.T) {
	mockSvc := &MockDriverService{}
	h := handler.NewDriverHandler(mockSvc)

	router := setupDriverRouter(h)

	req, _ := http.NewRequest("POST", "/drivers", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDriverByID(t *testing.T) {
	mockSvc := &MockDriverService{}
	h := &handler.DriverHandler{Service: mockSvc}

	router := setupDriverRouter(h)

	req, _ := http.NewRequest("GET", "/drivers/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response model.Driver
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "John Doe", response.Name)
}

func TestGetDriverByIDNotFound(t *testing.T) {
	mockSvc := &MockDriverService{}
	h := handler.NewDriverHandler(mockSvc)

	router := setupDriverRouter(h)

	req, _ := http.NewRequest("GET", "/drivers/2", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateDriver(t *testing.T) {
	mockSvc := &MockDriverService{}
	h := handler.NewDriverHandler(mockSvc)

	router := setupDriverRouter(h)

	payload := map[string]string{
		"name":                  "Updated Name",
		"email":                 "updated@example.com",
		"phone":                 "1111111111",
		"address":               "Updated Address",
		"driver_license_number": "DL789",
		"car_model_id":          "3",
		"car_type_id":           "3",
		"plate_number":          "UPD123",
		"status":                "inactive",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/drivers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Driver updated successfully")
}

func TestUpdateDriverInvalidJSON(t *testing.T) {
	mockSvc := &MockDriverService{}
	h := handler.NewDriverHandler(mockSvc)

	router := setupDriverRouter(h)

	req, _ := http.NewRequest("PUT", "/drivers/1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateDriverNotFound(t *testing.T) {
	mockSvc := &MockDriverService{}
	h := handler.NewDriverHandler(mockSvc)

	router := setupDriverRouter(h)

	payload := map[string]string{
		"name":                  "Updated Name",
		"email":                 "updated@example.com",
		"phone":                 "1111111111",
		"address":               "Updated Address",
		"driver_license_number": "DL789",
		"car_model_id":          "3",
		"car_type_id":           "3",
		"plate_number":          "UPD123",
		"status":                "inactive",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/drivers/2", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteDriver(t *testing.T) {
	mockSvc := &MockDriverService{}
	h := handler.NewDriverHandler(mockSvc)

	router := setupDriverRouter(h)

	req, _ := http.NewRequest("DELETE", "/drivers/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Driver deleted successfully")
}

func TestDeleteDriverNotFound(t *testing.T) {
	mockSvc := &MockDriverService{}
	h := handler.NewDriverHandler(mockSvc)

	router := setupDriverRouter(h)

	req, _ := http.NewRequest("DELETE", "/drivers/2", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
