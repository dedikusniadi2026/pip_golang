package service

import (
	"auth-service/model"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTripsRepositoryImpl struct {
	mock.Mock
}

func (m *MockTripsRepositoryImpl) Create(t *model.VehicleTrip) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTripsRepositoryImpl) Update(t *model.VehicleTrip) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTripsRepositoryImpl) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTripsRepositoryImpl) FindByVehicle(vehicleID uint) ([]model.VehicleTrip, error) {
	args := m.Called(vehicleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.VehicleTrip), args.Error(1)
}

func (m *MockTripsRepositoryImpl) GetTripTotal(ctx context.Context) (*model.TotalTrips, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TotalTrips), args.Error(1)
}

func TestTripService(t *testing.T) {
	trip := &model.VehicleTrip{ID: 1}

	t.Run("Create success", func(t *testing.T) {
		repo := new(MockTripsRepository)
		svc := NewTripService(repo)

		repo.On("Create", trip).Return(nil)

		err := svc.Create(trip)
		assert.NoError(t, err)
	})

	t.Run("Create error", func(t *testing.T) {
		repo := new(MockTripsRepository)
		svc := NewTripService(repo)

		repo.On("Create", mock.Anything).Return(errors.New("db error"))

		err := svc.Create(trip)
		assert.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		repo := new(MockTripsRepository)
		svc := NewTripService(repo)

		repo.On("Update", trip).Return(nil)

		err := svc.Update(trip)
		assert.NoError(t, err)
	})

	t.Run("Update error", func(t *testing.T) {
		repo := new(MockTripsRepository)
		svc := NewTripService(repo)

		repo.On("Update", mock.Anything).Return(errors.New("db error"))

		err := svc.Update(trip)
		assert.Error(t, err)
	})

	t.Run("Delete success", func(t *testing.T) {
		repo := new(MockTripsRepository)
		svc := NewTripService(repo)

		repo.On("Delete", uint(1)).Return(nil)

		err := svc.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("Delete error", func(t *testing.T) {
		repo := new(MockTripsRepository)
		svc := NewTripService(repo)

		repo.On("Delete", uint(1)).Return(errors.New("db error"))

		err := svc.Delete(1)
		assert.Error(t, err)
	})

	t.Run("FindByVehicle success", func(t *testing.T) {
		repo := new(MockTripsRepository)
		svc := NewTripService(repo)

		repo.On("FindByVehicle", uint(1)).
			Return([]model.VehicleTrip{*trip}, nil)

		result, err := svc.FindByVehicle(1)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("FindByVehicle error", func(t *testing.T) {
		repo := new(MockTripsRepository)
		svc := NewTripService(repo)

		repo.On("FindByVehicle", uint(1)).
			Return(nil, errors.New("db error"))

		result, err := svc.FindByVehicle(1)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("GetTripTotals success", func(t *testing.T) {
		repo := new(MockTripsRepository)
		svc := NewTripService(repo)

		total := &model.TotalTrips{TotalTrips: 10}
		repo.On("GetTripTotal", mock.Anything).
			Return(total, nil)

		result, err := svc.GetTripTotals(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, total, result)
	})

	t.Run("GetTripTotals error", func(t *testing.T) {
		repo := new(MockTripsRepository)
		svc := NewTripService(repo)

		repo.On("GetTripTotal", mock.Anything).
			Return(nil, errors.New("db error"))

		result, err := svc.GetTripTotals(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
