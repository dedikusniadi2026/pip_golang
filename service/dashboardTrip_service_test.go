package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDashboardTripRepository struct {
	mock.Mock
}

func (m *MockDashboardTripRepository) GetDashboardSummary(ctx context.Context) (*model.DashboardSummary, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.DashboardSummary), args.Error(1)
}

func TestDashboardTripService_GetDashboardSummary(t *testing.T) {
	mockRepo := new(MockDashboardTripRepository)
	svc := service.NewDashboardTripService(mockRepo)

	expected := &model.DashboardSummary{
		TotalTrips:    100,
		TotalDistance: 1500,
	}
	ctx := context.Background()
	mockRepo.On("GetDashboardSummary", ctx).Return(expected, nil)

	result, err := svc.GetDashboardSummary(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
