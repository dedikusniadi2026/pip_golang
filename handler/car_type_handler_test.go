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
	"github.com/stretchr/testify/mock"
)

type MockCarTypeService struct {
	mock.Mock
}

func (m *MockCarTypeService) GetAll() ([]model.CarType, error) {
	args := m.Called()
	return args.Get(0).([]model.CarType), args.Error(1)
}

func (m *MockCarTypeService) GetByID(id int) (*model.CarType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CarType), args.Error(1)
}

func (m *MockCarTypeService) Create(ct model.CarType) error {
	args := m.Called(ct)
	return args.Error(0)
}

func TestNewCarTypeHandler(t *testing.T) {
	mockSvc := new(MockCarTypeService)
	h := handler.NewCarTypeHandler(mockSvc)
	assert.NotNil(t, h)
}

func TestCarTypeHandler_GetAll_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockCarTypeService)
	h := handler.NewCarTypeHandler(mockSvc)

	mockSvc.On("GetAll").Return([]model.CarType{
		{ID: "1", TypeName: "SUV"},
	}, nil)

	r := gin.New()
	r.GET("/car-types", h.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/car-types", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCarTypeHandler_GetAll_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockCarTypeService)
	h := handler.NewCarTypeHandler(mockSvc)

	mockSvc.On("GetAll").Return([]model.CarType(nil), errors.New("db error"))

	r := gin.New()
	r.GET("/car-types", h.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/car-types", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCarTypeHandler_GetByID_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockCarTypeService)
	h := handler.NewCarTypeHandler(mockSvc)

	mockSvc.On("GetByID", 1).
		Return(&model.CarType{ID: "1", TypeName: "SUV"}, nil)

	r := gin.New()
	r.GET("/car-types/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/car-types/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCarTypeHandler_GetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockCarTypeService)
	h := handler.NewCarTypeHandler(mockSvc)

	mockSvc.On("GetByID", 2).Return(nil, nil)

	r := gin.New()
	r.GET("/car-types/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/car-types/2", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCarTypeHandler_GetByID_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockCarTypeService)
	h := handler.NewCarTypeHandler(mockSvc)

	r := gin.New()
	r.GET("/car-types/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/car-types/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCarTypeHandler_GetByID_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockCarTypeService)
	h := handler.NewCarTypeHandler(mockSvc)

	mockSvc.
		On("GetByID", 1).
		Return(nil, errors.New("db error"))

	r := gin.New()
	r.GET("/car-types/:id", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/car-types/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCarTypeHandler_Create_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockCarTypeService)
	h := handler.NewCarTypeHandler(mockSvc)

	mockSvc.On("Create", mock.Anything).Return(nil)

	r := gin.New()
	r.POST("/car-types", h.Create)

	req := httptest.NewRequest(
		http.MethodPost,
		"/car-types",
		strings.NewReader(`{"typeName":"SUV"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCarTypeHandler_Create_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockCarTypeService)
	h := handler.NewCarTypeHandler(mockSvc)

	r := gin.New()
	r.POST("/car-types", h.Create)

	req := httptest.NewRequest(
		http.MethodPost,
		"/car-types",
		strings.NewReader(`invalid json`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCarTypeHandler_Create_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockCarTypeService)
	h := handler.NewCarTypeHandler(mockSvc)

	mockSvc.
		On("Create", mock.Anything).
		Return(errors.New("insert failed"))

	r := gin.New()
	r.POST("/car-types", h.Create)

	req := httptest.NewRequest(
		http.MethodPost,
		"/car-types",
		strings.NewReader(`{"typeName":"SUV"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}
