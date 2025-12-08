package handler_test

import (
	"auth-service/handler"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type DashboardService interface {
	GetDashboardData() (map[string]interface{}, error)
}

type MockDashboardService struct{}

func (s *MockDashboardService) GetDashboardData() (map[string]interface{}, error) {
	data := map[string]interface{}{
		"totalUsers": 10,
		"totalSales": 5000,
	}
	return data, nil
}

type MockDashboardServiceError struct{}

func (m *MockDashboardServiceError) GetDashboardData() (map[string]interface{}, error) {
	return nil, errors.New("failed to fetch dashboard data")
}

func TestGetDashboard_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDashboardService{}
	h := handler.NewDashboardHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetDashboard(c)

	assert.Equal(t, http.StatusOK, w.Code)
	expected := `{"success":true,"data":{"totalSales":5000,"totalUsers":10}}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestGetDashboard_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDashboardServiceError{}
	h := handler.NewDashboardHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetDashboard(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	expected := `{"error":"failed to fetch dashboard data"}`
	assert.JSONEq(t, expected, w.Body.String())

}

func TestDashboardHandler_GetDashboard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockDashboardService{}
	h := handler.NewDashboardHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetDashboard(c)

	assert.Equal(t, http.StatusOK, w.Code)
	expected := `{"success":true,"data":{"totalSales":5000,"totalUsers":10}}`
	assert.JSONEq(t, expected, w.Body.String())
}
