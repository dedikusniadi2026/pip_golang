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

type MockPopularDestinationService struct{}

func (m *MockPopularDestinationService) GetAll() ([]model.PopularDestination, error) {
	return []model.PopularDestination{
		{Destination: "Bali", Bookings: 100},
		{Destination: "Jakarta", Bookings: 50},
	}, nil
}

type MockErrorService struct{}

func (m *MockErrorService) GetAll() ([]model.PopularDestination, error) {
	return nil, errors.New("service error")
}

func TestGetAll_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &MockPopularDestinationService{}
	handler := &handler.PopularDestinationHandler{Service: service}

	req := httptest.NewRequest(http.MethodGet, "/popular-destinations", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetAll(c)

	assert.Equal(t, 200, w.Code)
	expected := `[{"bookings":100,"destination":"Bali"},{"bookings":50,"destination":"Jakarta"}]`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestGetAll_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &MockErrorService{}
	handler := &handler.PopularDestinationHandler{Service: service}

	req := httptest.NewRequest(http.MethodGet, "/popular-destinations", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetAll(c)

	assert.Equal(t, 500, w.Code)
	expected := `{"error":"service error"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestPopularDestinationHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	service := &MockPopularDestinationService{}
	handler := &handler.PopularDestinationHandler{Service: service}

	req := httptest.NewRequest(http.MethodGet, "/popular-destinations", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetAll(c)

	assert.Equal(t, 200, w.Code)
	expected := `[{"bookings":100,"destination":"Bali"},{"bookings":50,"destination":"Jakarta"}]`
	assert.JSONEq(t, expected, w.Body.String())
}
