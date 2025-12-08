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

type MockCarTypeService struct {
	ReturnError bool
}

func (m *MockCarTypeService) Create(ct model.CarType) error {
	if m.ReturnError {
		return errors.New("failed to create car type")
	}
	return nil
}

func (m *MockCarTypeService) GetAll() ([]model.CarType, error) {
	if m.ReturnError {
		return nil, errors.New("failed to fetch car types")
	}
	return []model.CarType{
		{ID: "1", TypeName: "Type A"},
		{ID: "2", TypeName: "Type B"},
	}, nil
}

func (m *MockCarTypeService) GetByID(id int) (*model.CarType, error) {
	if m.ReturnError {
		return nil, errors.New("service error")
	}
	if id != 1 {
		return nil, nil
	}
	return &model.CarType{ID: "1", TypeName: "Type A"}, nil
}

func TestCarTypeHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarTypeService{}
	h := &handler.CarTypeHandler{Service: mockSvc}

	router := gin.Default()
	router.GET("/car-types", h.GetAll)

	req, _ := http.NewRequest("GET", "/car-types", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestCarTypeHandler_GetAll_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarTypeService{ReturnError: true}
	h := &handler.CarTypeHandler{Service: mockSvc}

	router := gin.Default()
	router.GET("/car-types", h.GetAll)

	req, _ := http.NewRequest("GET", "/car-types", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

func TestCarTypeHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarTypeService{}

	h := &handler.CarTypeHandler{Service: mockSvc}

	router := gin.Default()
	router.GET("/car-types/:id", h.GetByID)

	req, _ := http.NewRequest("GET", "/car-types/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", "/car-types/2", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	req, _ = http.NewRequest("GET", "/car-types/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCarTypeHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarTypeService{}
	h := &handler.CarTypeHandler{Service: mockSvc}

	router := gin.Default()
	router.POST("/car-types", h.Create)

	payload := `{"name":"New Type"}`
	req, _ := http.NewRequest("POST", "/car-types", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	payload = `invalid`
	req, _ = http.NewRequest("POST", "/car-types", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
