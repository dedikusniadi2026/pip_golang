package model

import "time"

type VehicleMaintenance struct {
	ID          uint `gorm:"primaryKey"`
	VehicleID   uint
	ServiceDate time.Time
	ServiceType string
	Description string
	Mileage     int
	Cost        float64
}
