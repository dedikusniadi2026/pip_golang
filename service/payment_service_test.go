package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) GetPayments(page, pageSize int) ([]model.Payment, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetPaymentStats(ctx context.Context) (*model.PaymentStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.PaymentStats), args.Error(1)
}

func (m *MockPaymentRepository) GetAll(ctx context.Context) ([]model.Payment, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) GetByID(ctx context.Context, id int) (*model.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Create(ctx context.Context, p *model.Payment) (int, error) {
	args := m.Called(ctx, p)
	return args.Int(0), args.Error(1)
}

func (m *MockPaymentRepository) Update(ctx context.Context, p *model.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPaymentRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestPaymentService_GetPayments(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)

	expected := []model.Payment{
		{PaymentID: 1, BookingID: 1, Customer: "John Doe", Driver: "Jane Doe", Amount: 100.0, Method: "credit", Status: "paid", PaymentDate: "2023-01-01"},
		{PaymentID: 2, BookingID: 2, Customer: "Alice", Driver: "Bob", Amount: 200.0, Method: "debit", Status: "pending", PaymentDate: "2023-01-02"},
	}
	mockRepo.On("GetPayments", 1, 10).Return(expected, nil)

	result, err := svc.GetPayments(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_GetPayments_Error(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)

	mockRepo.On("GetPayments", 1, 10).Return([]model.Payment{}, assert.AnError)

	result, err := svc.GetPayments(1, 10)
	assert.Error(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_GetPaymentByID(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	expected := &model.Payment{PaymentID: 1, BookingID: 1, Customer: "John Doe", Driver: "Jane Doe", Amount: 100.0, Method: "credit", Status: "paid", PaymentDate: "2023-01-01"}
	mockRepo.On("GetByID", ctx, 1).Return(expected, nil)

	result, err := svc.GetPaymentByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_GetPaymentByID_NotFound(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, 1).Return((*model.Payment)(nil), nil)

	result, err := svc.GetPaymentByID(ctx, 1)
	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_GetPaymentByID_Error(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, 1).Return((*model.Payment)(nil), assert.AnError)

	result, err := svc.GetPaymentByID(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_CreatePayment(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	payment := &model.Payment{BookingID: 1, Customer: "John Doe", Driver: "Jane Doe", Amount: 100.0, Method: "credit", Status: "paid"}
	mockRepo.On("Create", ctx, payment).Return(1, nil)

	id, err := svc.CreatePayment(ctx, payment)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_CreatePayment_Error(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	payment := &model.Payment{BookingID: 1, Customer: "John Doe", Driver: "Jane Doe", Amount: 100.0, Method: "credit", Status: "paid"}
	mockRepo.On("Create", ctx, payment).Return(0, assert.AnError)

	id, err := svc.CreatePayment(ctx, payment)
	assert.Error(t, err)
	assert.Equal(t, 0, id)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_UpdatePayment(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	payment := &model.Payment{PaymentID: 1, BookingID: 1, Customer: "John Doe", Driver: "Jane Doe", Amount: 100.0, Method: "credit", Status: "paid"}
	mockRepo.On("Update", ctx, payment).Return(nil)

	err := svc.UpdatePayment(ctx, payment)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_UpdatePayment_Error(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	payment := &model.Payment{PaymentID: 1, BookingID: 1, Customer: "John Doe", Driver: "Jane Doe", Amount: 100.0, Method: "credit", Status: "paid"}
	mockRepo.On("Update", ctx, payment).Return(assert.AnError)

	err := svc.UpdatePayment(ctx, payment)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_DeletePayment(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, 1).Return(nil)

	err := svc.DeletePayment(ctx, 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_DeletePayment_Error(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, 1).Return(assert.AnError)

	err := svc.DeletePayment(ctx, 1)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_GetPaymentStats(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	expected := &model.PaymentStats{TotalPayment: 1000, PendingPayment: 200, TotalTransactions: 10}
	mockRepo.On("GetPaymentStats", ctx).Return(expected, nil)

	result, err := svc.GetPaymentStats(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_GetPaymentStats_Error(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetPaymentStats", ctx).Return((*model.PaymentStats)(nil), assert.AnError)

	result, err := svc.GetPaymentStats(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_GetAllPayments(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	expected := []model.Payment{
		{PaymentID: 1, BookingID: 1, Customer: "John Doe", Driver: "Jane Doe", Amount: 100.0, Method: "credit", Status: "paid", PaymentDate: "2023-01-01"},
		{PaymentID: 2, BookingID: 2, Customer: "Alice", Driver: "Bob", Amount: 200.0, Method: "debit", Status: "pending", PaymentDate: "2023-01-02"},
	}
	mockRepo.On("GetAll", ctx).Return(expected, nil)

	result, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_GetAllPayments_Error(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetAll", ctx).Return([]model.Payment{}, assert.AnError)

	result, err := svc.GetAll(ctx)
	assert.Error(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_GetAll(t *testing.T) {
	mockRepo := new(MockPaymentRepository)
	svc := service.NewPaymentService(mockRepo)
	ctx := context.Background()

	expected := []model.Payment{
		{PaymentID: 1, BookingID: 1, Customer: "John Doe", Driver: "Jane Doe", Amount: 100.0, Method: "credit", Status: "paid", PaymentDate: "2023-01-01"},
	}
	mockRepo.On("GetAll", ctx).Return(expected, nil)

	result, err := svc.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
