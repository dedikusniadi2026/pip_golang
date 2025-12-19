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
	if m.ReturnError {
		return errors.New("service error")
	}
	if b.ID == "" {
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

func TestDeleteBooking_Final(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockSvc := &MockBookingService{}
		h := handler.NewBookingHandler(mockSvc)

		router := gin.New()
		router.DELETE("/bookings/:id", h.Delete)

		req, _ := http.NewRequest("DELETE", "/bookings/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Booking deleted successfully")
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := &MockBookingService{ReturnError: true}
		h := handler.NewBookingHandler(mockSvc)

		router := gin.New()
		router.DELETE("/bookings/:id", h.Delete)

		req, _ := http.NewRequest("DELETE", "/bookings/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "id required")
	})

	t.Run("Empty ID (direct service call)", func(t *testing.T) {
		mockSvc := &MockBookingService{}
		err := mockSvc.Delete("")
		assert.Error(t, err)
		assert.Equal(t, "id required", err.Error())
	})
}

func TestCreateBooking_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.New()
	router.POST("/bookings", h.Create)

	payload := model.Booking{ID: "123", Customer: "Test Booking"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateBooking_BadRequest_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.New()
	router.POST("/bookings", h.Create)

	body := []byte(`{"customer":`)
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBooking_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{ReturnError: true}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.New()
	router.POST("/bookings", h.Create)

	payload := model.Booking{ID: "123", Customer: "Test Booking"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAllBooking_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.New()
	router.GET("/bookings", h.GetAll)

	req, _ := http.NewRequest("GET", "/bookings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAllBooking_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{ReturnError: true}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.New()
	router.GET("/bookings", h.GetAll)

	req, _ := http.NewRequest("GET", "/bookings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateBooking_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.New()
	router.PUT("/bookings/:id", h.Update)

	payload := model.Booking{Customer: "Updated Booking"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/bookings/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateBooking_BadRequest_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.New()
	router.PUT("/bookings/:id", h.Update)

	body := []byte(`{"customer":`)
	req, _ := http.NewRequest("PUT", "/bookings/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateBooking_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockBookingService{ReturnError: true}
	h := handler.NewBookingHandler(mockSvc)

	router := gin.New()
	router.PUT("/bookings/:id", h.Update)

	payload := model.Booking{Customer: "Updated Booking"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/bookings/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
