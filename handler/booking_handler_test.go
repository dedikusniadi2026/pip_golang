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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockBookingService struct {
	ReturnError bool
}

func (m *MockBookingService) Create(b *model.Booking) error {
	if m.ReturnError || b.ID == "" {
		return errors.New("user id required")
	}
	b.ID = "1"
	return nil
}

func (m *MockBookingService) GetAll() ([]model.Booking, error) {
	if m.ReturnError {
		return nil, errors.New("get all error")
	}
	return []model.Booking{
		{ID: "1", Customer: "123", Driver: "456", Place: "Test Booking"},
	}, nil
}

func (m *MockBookingService) Update(b *model.Booking) error {
	if m.ReturnError || b.ID == "" {
		return errors.New("id required")
	}
	return nil
}

func (m *MockBookingService) Delete(id string) error {
	if m.ReturnError || id == "" {
		return errors.New("id required")
	}
	return nil
}

func TestCreateBooking(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.Default()
	router.POST("/bookings", h.Create)

	payload := model.Booking{ID: "123", Customer: "Test Booking"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetAllBooking(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.Default()
	router.GET("/bookings", h.GetAll)

	req, _ := http.NewRequest("GET", "/bookings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateBooking(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.Default()
	router.PUT("/bookings/:id", h.Update)

	payload := model.Booking{ID: "123", Customer: "Updated Booking"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/bookings/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteBooking(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.Default()
	router.DELETE("/bookings/:id", h.Delete)

	req, _ := http.NewRequest("DELETE", "/bookings/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateBooking_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{ReturnError: true}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.Default()
	router.POST("/bookings", h.Create)

	payload := model.Booking{ID: "123", Customer: "Test Booking"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAllBooking_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{ReturnError: true}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.Default()
	router.GET("/bookings", h.GetAll)

	req, _ := http.NewRequest("GET", "/bookings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateBooking_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{ReturnError: true}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.Default()
	router.PUT("/bookings/:id", h.Update)

	payload := model.Booking{ID: "123", Customer: "Updated Booking"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/bookings/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteBooking_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{ReturnError: true}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.Default()
	router.DELETE("/bookings/:id", h.Delete)

	req, _ := http.NewRequest("DELETE", "/bookings/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateBooking_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.Default()
	router.POST("/bookings", h.Create)

	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
