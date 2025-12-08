package handler_test

import (
	"auth-service/handler"
	"auth-service/model"
	"auth-service/service"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockDashboardTripService struct {
	service.DashboardTripService
}

func (m *MockDashboardTripService) GetDashboardSummary(ctx context.Context) (*model.DashboardSummary, error) {
	return &model.DashboardSummary{
		TotalTrips:    10,
		TotalDistance: 100,
		TotalDuration: 200,
		LastTrip: model.Trip{
			ID: 1,
		},
	}, nil
}

type MockDashboardTripServiceError struct {
	service.DashboardTripService
}

func (m *MockDashboardTripServiceError) GetDashboardSummary(ctx context.Context) (*model.DashboardSummary, error) {
	return nil, errors.New("service error")
}

func setupRouterDashboardTrip(h *handler.DashboardTripHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/dashboard", h.GetDashboard)

	return r
}

func TestGetDashboardSuccess(t *testing.T) {
	mockSvc := &MockDashboardTripService{}
	h := handler.NewDashboardTripHandler(mockSvc)

	router := setupRouterDashboardTrip(h)

	req, _ := http.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"totalTrips":10`)
}

func TestGetDashboardError(t *testing.T) {
	mockSvc := &MockDashboardTripServiceError{}
	h := handler.NewDashboardTripHandler(mockSvc)

	router := setupRouterDashboardTrip(h)

	req, _ := http.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"service error"`)
}

func TestDashboardTripHandler_GetDashboard(t *testing.T) {
	mockSvc := &MockDashboardTripService{}
	h := handler.NewDashboardTripHandler(mockSvc)

	router := setupRouterDashboardTrip(h)

	req, _ := http.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"totalTrips":10`)
}
