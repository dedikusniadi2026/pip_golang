package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"auth-service/handler"
	"auth-service/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) GetPayments(page, pageSize int) ([]model.Payment, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Payment), args.Error(1)
}

func (m *MockPaymentService) GetPaymentByID(ctx context.Context, id int) (*model.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentService) CreatePayment(ctx context.Context, p *model.Payment) (int, error) {
	args := m.Called(ctx, p)
	return args.Int(0), args.Error(1)
}

func (m *MockPaymentService) UpdatePayment(ctx context.Context, p *model.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPaymentService) DeletePayment(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentService) GetPaymentStats(ctx context.Context) (*model.PaymentStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PaymentStats), args.Error(1)
}

func setupPaymentRouter(h *handler.PaymentHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/payments", h.GetPayments)
	r.GET("/payments/:id", h.GetPaymentByID)
	r.POST("/payments", h.CreatePayment)
	r.PUT("/payments/:id", h.UpdatePayment)
	r.DELETE("/payments/:id", h.DeletePayment)
	r.GET("/payments/stats", h.GetPaymentStats)

	return r
}

func TestGetPaymentByID_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockSvc}

	mockSvc.
		On("GetPaymentByID", mock.Anything, 1).
		Return(nil, errors.New("db error"))

	r := gin.New()
	r.GET("/payments/:id", h.GetPaymentByID)

	req := httptest.NewRequest(http.MethodGet, "/payments/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "db error")

	mockSvc.AssertExpectations(t)
}

func TestCreatePayment_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockSvc}

	mockSvc.
		On("CreatePayment", mock.Anything, mock.Anything).
		Return(0, errors.New("insert failed"))

	r := gin.New()
	r.POST("/payments", h.CreatePayment)

	payload := `{"amount":1000}`
	req := httptest.NewRequest(
		http.MethodPost,
		"/payments",
		strings.NewReader(payload),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "insert failed")

	mockSvc.AssertExpectations(t)
}

func TestGetPayments(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	payments := []model.Payment{
		{PaymentID: 1, BookingID: 1, Customer: "John", Amount: 100.0},
	}
	mockService.On("GetPayments", 1, 5).Return(payments, nil)

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("GET", "/payments?page=1&pageSize=5", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []model.Payment
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, payments, response)
	mockService.AssertExpectations(t)
}

func TestGetPayments_InvalidPage(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("GET", "/payments?page=abc&pageSize=5", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid page")
}

func TestGetPayments_InvalidPageSize(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("GET", "/payments?page=1&pageSize=abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid pageSize")
}

func TestGetPaymentByID_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	payment := &model.Payment{PaymentID: 1, BookingID: 1, Customer: "John", Amount: 100.0}
	mockService.On("GetPaymentByID", mock.Anything, 1).Return(payment, nil)

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("GET", "/payments/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response model.Payment
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, *payment, response)
	mockService.AssertExpectations(t)
}

func TestGetPaymentByID_NotFound(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	mockService.On("GetPaymentByID", mock.Anything, 1).Return(nil, nil)

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("GET", "/payments/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "payment not found")
	mockService.AssertExpectations(t)
}

func TestGetPaymentByID_InvalidID(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("GET", "/payments/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid payment id")
}

func TestCreatePayment_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	payment := model.Payment{BookingID: 1, Customer: "John", Amount: 100.0}
	mockService.On("CreatePayment", mock.Anything, &payment).Return(1, nil)

	router := setupPaymentRouter(h)

	body, _ := json.Marshal(payment)
	req, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "payment created successfully")
	mockService.AssertExpectations(t)
}

func TestCreatePayment_InvalidJSON(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("POST", "/payments", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePayment_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	payment := model.Payment{PaymentID: 1, BookingID: 1, Customer: "John", Amount: 100.0}
	mockService.On("UpdatePayment", mock.Anything, &payment).Return(nil)

	router := setupPaymentRouter(h)

	body, _ := json.Marshal(payment)
	req, _ := http.NewRequest("PUT", "/payments/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "payment updated successfully")
	mockService.AssertExpectations(t)
}

func TestUpdatePayment_InvalidID(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("PUT", "/payments/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid payment id")
}

func TestGetPaymentStats_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	stats := &model.PaymentStats{TotalPayment: 1000, PendingPayment: 100, TotalTransactions: 10}
	mockService.On("GetPaymentStats", mock.Anything).Return(stats, nil)

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("GET", "/payments/stats", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(1000), response["total_payment"])
	mockService.AssertExpectations(t)
}

func TestDeletePayment_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	mockService.On("DeletePayment", mock.Anything, 1).Return(nil)

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("DELETE", "/payments/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "payment deleted successfully")
	mockService.AssertExpectations(t)
}

func TestDeletePayment_InvalidID(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("DELETE", "/payments/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid payment id")
}

func TestGetPayments_Error(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	mockService.On("GetPayments", 1, 5).Return(nil, errors.New("get payments error"))

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("GET", "/payments?page=1&pageSize=5", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdatePayment_Error(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	payment := model.Payment{PaymentID: 1, BookingID: 1, Customer: "John", Amount: 100.0}
	mockService.On("UpdatePayment", mock.Anything, &payment).Return(errors.New("update error"))

	router := setupPaymentRouter(h)

	body, _ := json.Marshal(payment)
	req, _ := http.NewRequest("PUT", "/payments/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestPaymentHandler_GetPaymentByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	r := gin.New()
	r.GET("/payments/:id", h.GetPaymentByID)

	req := httptest.NewRequest(http.MethodGet, "/payments/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid payment id")
}

func TestPaymentHandler_CreatePayment_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	r := gin.New()
	r.POST("/payments", h.CreatePayment)

	req := httptest.NewRequest(
		http.MethodPost,
		"/payments",
		strings.NewReader("{invalid json"),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPaymentHandler_UpdatePayment_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	r := gin.New()
	r.PUT("/payments/:id", h.UpdatePayment)

	payload := `{"amount":1000}`
	req := httptest.NewRequest(
		http.MethodPut,
		"/payments/xyz",
		strings.NewReader(payload),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPaymentStats_Error(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	mockService.On("GetPaymentStats", mock.Anything).Return(nil, errors.New("stats error"))

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("GET", "/payments/stats", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestPaymentHandler_UpdatePayment_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	r := gin.New()
	r.PUT("/payments/:id", h.UpdatePayment)

	req := httptest.NewRequest(
		http.MethodPut,
		"/payments/1",
		strings.NewReader("invalid json"),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeletePayment_Error(t *testing.T) {
	mockService := new(MockPaymentService)
	h := &handler.PaymentHandler{Service: mockService}

	mockService.On("DeletePayment", mock.Anything, 1).Return(errors.New("delete error"))

	router := setupPaymentRouter(h)

	req, _ := http.NewRequest("DELETE", "/payments/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
