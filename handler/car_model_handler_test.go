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

type MockCarModelService struct {
	ReturnError bool
}

func TestNewCarModelRepositoryHandler(t *testing.T) {
	mockSvc := &MockCarModelService{}
	h := handler.NewCarModelRepositoryHandler(mockSvc)
	assert.NotNil(t, h)
}

func (m *MockCarModelService) GetAll() ([]model.CarModel, error) {
	if m.ReturnError {
		return nil, errors.New("failed to fetch car models")
	}
	return []model.CarModel{
		{ID: 1, ModelName: "Model A"},
		{ID: 2, ModelName: "Model B"},
	}, nil
}

func (m *MockCarModelService) GetByID(id int) (*model.CarModel, error) {
	if m.ReturnError {
		return nil, errors.New("car model not found")
	}
	if id != 1 {
		return nil, nil
	}
	return &model.CarModel{ID: 1, ModelName: "Model A"}, nil
}

func (m *MockCarModelService) Create(cm model.CarModel) error {
	if m.ReturnError {
		return errors.New("failed to create car model")
	}
	return nil
}

func TestCarModelHandler_GetAll(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarModelService{}

	h := &handler.CarModelHandler{Service: mockSvc}

	router := gin.Default()
	router.GET("/car-models", h.GetAll)

	req, _ := http.NewRequest("GET", "/car-models", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestCarModelHandler_GetByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarModelService{}
	h := &handler.CarModelHandler{Service: mockSvc}

	router := gin.Default()
	router.GET("/car-models/:id", h.GetByID)

	req, _ := http.NewRequest("GET", "/car-models/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCarModelHandler_GetByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarModelService{}
	h := &handler.CarModelHandler{Service: mockSvc}

	router := gin.Default()
	router.GET("/car-models/:id", h.GetByID)

	req, _ := http.NewRequest("GET", "/car-models/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCarModelHandler_GetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarModelService{}
	h := &handler.CarModelHandler{Service: mockSvc}

	router := gin.Default()
	router.GET("/car-models/:id", h.GetByID)

	req, _ := http.NewRequest("GET", "/car-models/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCarModelHandler_GetByID_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarModelService{}
	mockSvc.ReturnError = true
	h := &handler.CarModelHandler{Service: mockSvc}

	router := gin.Default()
	router.GET("/car-models/:id", h.GetByID)

	req, _ := http.NewRequest("GET", "/car-models/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCarModelHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarModelService{}
	h := &handler.CarModelHandler{Service: mockSvc}

	router := gin.Default()
	router.POST("/car-models", h.Create)

	payload := `{"model_name":"New Model"}`
	req, _ := http.NewRequest("POST", "/car-models", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCarModelHandler_Create_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarModelService{}
	mockSvc.ReturnError = true
	h := &handler.CarModelHandler{Service: mockSvc}

	router := gin.Default()
	router.POST("/car-models", h.Create)

	payload := `{"model_name":"New Model"}`
	req, _ := http.NewRequest("POST", "/car-models", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCarModelHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarModelService{}
	h := &handler.CarModelHandler{Service: mockSvc}

	router := gin.Default()
	router.POST("/car-models", h.Create)

	req, _ := http.NewRequest("POST", "/car-models", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCarModelHandler_GetAll_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarModelService{
		ReturnError: true,
	}

	h := &handler.CarModelHandler{Service: mockSvc}

	r := gin.New()
	r.GET("/car-models", h.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/car-models", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to fetch car models")
}

func TestCarModelHandler_GetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockCarModelService{
		ReturnError: false,
	}

	h := &handler.CarModelHandler{
		Service: mockSvc,
	}

	r := gin.New()
	r.GET("/car-models", h.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/car-models", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
