package service_test

import (
	"auth-service/model"
	"auth-service/service"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTripsRepository struct {
	mock.Mock
}

func (m *MockTripsRepository) Create(t *model.VehicleTrip) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTripsRepository) Update(t *model.VehicleTrip) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTripsRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTripsRepository) FindByVehicle(vehicleID uint) ([]model.VehicleTrip, error) {
	args := m.Called(vehicleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.VehicleTrip), args.Error(1)
}

func (m *MockTripsRepository) GetTripTotal(ctx context.Context) (*model.TotalTrips, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TotalTrips), args.Error(1)
}

func TestTripService_BlackBox(t *testing.T) {
	trip := &model.VehicleTrip{ID: 1}

	t.Run("Create success", func(t *testing.T) {
		mockRepo := new(MockTripsRepository)
		svc := service.NewTripService(mockRepo)

		mockRepo.On("Create", trip).Return(nil)
		err := svc.Create(trip)
		assert.NoError(t, err)
	})

	t.Run("Create error", func(t *testing.T) {
		mockRepo := new(MockTripsRepository)
		svc := service.NewTripService(mockRepo)

		mockRepo.On("Create", mock.MatchedBy(func(t *model.VehicleTrip) bool { return t != nil })).Return(errors.New("db error"))
		err := svc.Create(trip)
		assert.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		mockRepo := new(MockTripsRepository)
		svc := service.NewTripService(mockRepo)

		mockRepo.On("Update", trip).Return(nil)
		err := svc.Update(trip)
		assert.NoError(t, err)
	})

	t.Run("Update error", func(t *testing.T) {
		mockRepo := new(MockTripsRepository)
		svc := service.NewTripService(mockRepo)

		mockRepo.On("Update", mock.MatchedBy(func(t *model.VehicleTrip) bool { return t != nil })).Return(errors.New("db error"))
		err := svc.Update(trip)
		assert.Error(t, err)
	})

	t.Run("Delete success", func(t *testing.T) {
		mockRepo := new(MockTripsRepository)
		svc := service.NewTripService(mockRepo)

		mockRepo.On("Delete", uint(1)).Return(nil)
		err := svc.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Delete error", func(t *testing.T) {
		mockRepo := new(MockTripsRepository)
		svc := service.NewTripService(mockRepo)

		mockRepo.On("Delete", mock.MatchedBy(func(id uint) bool { return id > 0 })).Return(errors.New("db error"))
		err := svc.Delete(1)
		assert.Error(t, err)
	})

	t.Run("FindByVehicle success", func(t *testing.T) {
		mockRepo := new(MockTripsRepository)
		svc := service.NewTripService(mockRepo)

		mockRepo.On("FindByVehicle", uint(1)).Return([]model.VehicleTrip{*trip}, nil)
		result, err := svc.FindByVehicle(1)
		assert.NoError(t, err)
		assert.Equal(t, []model.VehicleTrip{*trip}, result)
	})

	t.Run("FindByVehicle error", func(t *testing.T) {
		mockRepo := new(MockTripsRepository)
		svc := service.NewTripService(mockRepo)

		mockRepo.On("FindByVehicle", mock.MatchedBy(func(id uint) bool { return id > 0 })).Return(nil, errors.New("db error"))
		result, err := svc.FindByVehicle(1)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("GetTripTotals success", func(t *testing.T) {
		mockRepo := new(MockTripsRepository)
		svc := service.NewTripService(mockRepo)

		total := &model.TotalTrips{TotalTrips: 10}
		mockRepo.On("GetTripTotal", mock.Anything).Return(total, nil)

		result, err := svc.GetTripTotals(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, total, result)
	})

	t.Run("GetTripTotals error", func(t *testing.T) {
		mockRepo := new(MockTripsRepository)
		svc := service.NewTripService(mockRepo)

		mockRepo.On("GetTripTotal", mock.Anything).Return(nil, errors.New("db error"))
		result, err := svc.GetTripTotals(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
