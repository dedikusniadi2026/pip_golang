package model

import (
	"database/sql"
	"time"
)

type DriverAssignment struct {
	ID         uint `gorm:"primaryKey"`
	VehicleID  uint
	DriverID   sql.NullInt64
	DriverName string
	StartDate  time.Time
	EndDate    time.Time
	TotalTrips int
	Status     string
}
