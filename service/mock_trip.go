package service

import (
	"auth-service/model"
	"context"

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
