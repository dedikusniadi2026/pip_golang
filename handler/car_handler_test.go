package handler_test

import (
	"auth-service/handler"
	"auth-service/model"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockCarService struct {
	ReturnError bool
}

func (m *MockCarService) GetAll() ([]model.Car, error) {
	if m.ReturnError {
		return nil, errors.New("failed to fetch cars")
	}
	return []model.Car{
		{ID: 1, Brand: "Car A"},
		{ID: 2, Brand: "Car B"},
	}, nil
}

func (m *MockCarService) GetByID(id int) (*model.Car, error) {
	if m.ReturnError {
		return nil, errors.New("service error")
	}
	if id != 1 {
		return nil, nil
	}
	return &model.Car{ID: 1, Brand: "Car A"}, nil
}

func (m *MockCarService) Create(car model.Car) error {
	if m.ReturnError {
		return errors.New("failed to create car")
	}
	return nil
}

func (m *MockCarService) Update(id int, car model.Car) error {
	if m.ReturnError {
		return errors.New("failed to update car")
	}
	if id != 1 {
		return errors.New("failed to update car")
	}
	return nil
}

func (m *MockCarService) Delete(id int) error {
	if m.ReturnError {
		return errors.New("failed to delete car")
	}
	if id != 1 {
		return errors.New("failed to delete car")
	}
	return nil
}

func TestCarHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarService{}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.GET("/cars", h.GetAll)

	req, _ := http.NewRequest("GET", "/cars", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCarHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.GET("/cars/:id", h.GetByID)

	req, _ := http.NewRequest("GET", "/cars/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/cars/2", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCarHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.POST("/cars", h.Create)

	payload := `{"name":"Car X"}`
	req, _ := http.NewRequest("POST", "/cars", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCarHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.PUT("/cars/:id", h.Update)

	payload := `{"name":"Car Updated"}`
	req, _ := http.NewRequest("PUT", "/cars/1", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCarHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.DELETE("/cars/:id", h.Delete)

	req, _ := http.NewRequest("DELETE", "/cars/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCarHandler_GetAll_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarService{ReturnError: true}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.GET("/cars", h.GetAll)

	req, _ := http.NewRequest("GET", "/cars", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCarHandler_GetByID_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{ReturnError: true}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.GET("/cars/:id", h.GetByID)

	req, _ := http.NewRequest("GET", "/cars/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCarHandler_Create_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{ReturnError: true}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.POST("/cars", h.Create)

	payload := `{"name":"Car X"}`
	req, _ := http.NewRequest("POST", "/cars", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCarHandler_Update_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{ReturnError: true}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.PUT("/cars/:id", h.Update)

	payload := `{"name":"Car Updated"}`
	req, _ := http.NewRequest("PUT", "/cars/1", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCarHandler_Update_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{ReturnError: false}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.PUT("/cars/:id", h.Update)

	payload := `{"name":"Car Updated"}`
	req, _ := http.NewRequest("PUT", "/cars/2", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCarHandler_Delete_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{ReturnError: true}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.DELETE("/cars/:id", h.Delete)

	req, _ := http.NewRequest("DELETE", "/cars/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCarHandler_GetByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.GET("/cars/:id", h.GetByID)

	req, _ := http.NewRequest("GET", "/cars/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCarHandler_Update_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{ReturnError: true}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.PUT("/cars/:id", h.Update)

	payload := `{"name":"Car Updated"}`
	req, _ := http.NewRequest("PUT", "/cars/abc", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCarHandler_Delete_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockCarService{ReturnError: true}
	h := handler.NewCarHandler(mockSvc)

	router := gin.Default()
	router.DELETE("/cars/:id", h.Delete)

	req, _ := http.NewRequest("DELETE", "/cars/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
