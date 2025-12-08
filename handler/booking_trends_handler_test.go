package handler_test

import (
	"auth-service/handler"
	"auth-service/model"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockBookingTrendService struct {
	ReturnError bool
}

func (m *MockBookingTrendService) GetTrends(year int) ([]model.BookingTrend, error) {
	if m.ReturnError {
		return nil, errors.New("no trip found")
	}
	return []model.BookingTrend{
		{Month: "January", Booking: 12},
		{Month: "February", Booking: 15},
	}, nil
}

func TestBookingTrendsHandler_GetTrends_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockBookingTrendService{ReturnError: false}
	h := handler.NewBookingTrendsHandler(mockSvc)

	router := gin.Default()
	router.GET("/trends", h.GetTrends)

	req, _ := http.NewRequest("GET", "/trends?year=2024", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBookingTrendsHandler_GetTrends_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockBookingTrendService{ReturnError: true}
	h := handler.NewBookingTrendsHandler(mockSvc)

	router := gin.Default()
	router.GET("/trends", h.GetTrends)

	req, _ := http.NewRequest("GET", "/trends?year=2024", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBookingTrendsHandler_GetTrends_InvalidYear(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockBookingTrendService{ReturnError: false}
	h := handler.NewBookingTrendsHandler(mockSvc)

	router := gin.Default()
	router.GET("/trends", h.GetTrends)

	req, _ := http.NewRequest("GET", "/trends?year=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
