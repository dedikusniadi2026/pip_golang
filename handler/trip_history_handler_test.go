package handler

import (
	"auth-service/model"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockTripHistoryService struct {
	ReturnError bool
}

func (m *MockTripHistoryService) GetTripHistory() ([]model.TripHistory, error) {
	if m.ReturnError {
		return nil, errors.New("no trip found")
	}

	return []model.TripHistory{
		{ID: 1, CustomerName: "123", DriverName: "Trip A"},
		{ID: 2, CustomerName: "123", DriverName: "Trip B"},
	}, nil
}

func TestGetTripHistory_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockTripHistoryService{ReturnError: false}

	h := NewTripHistoryHandler(mockSvc)

	router := gin.Default()
	router.GET("/trips", h.GetTripHistory)

	req, _ := http.NewRequest("GET", "/trips", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTripHistory_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockTripHistoryService{ReturnError: true}
	h := NewTripHistoryHandler(mockSvc)

	router := gin.Default()
	router.GET("/trips", h.GetTripHistory)

	req, _ := http.NewRequest("GET", "/trips", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
