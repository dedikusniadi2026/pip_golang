package model

import (
	"time"
)

type Car struct {
	ID int `json:"id"`

	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Year        int    `json:"year"`
	PlateNumber string `json:"plate_number"`
	Capacity    int    `json:"capacity"`
	Color       string `json:"color"`
	DriverID    int    `json:"driver_id"`

	LastMaintenanceDate *time.Time `json:"last_maintenance_date"`
	CurrentKM           int        `json:"current_km"`
	Status              string     `json:"status"`

	Maintenance []VehicleMaintenance `gorm:"foreignKey:VehicleID"`
	Assignments []DriverAssignment   `gorm:"foreignKey:VehicleID"`
	Trips       []VehicleTrip        `gorm:"foreignKey:VehicleID"`
}

type Trip struct {
	ID       int     `json:"id"`
	CarID    int     `json:"car_id"`
	Distance int     `json:"distance_km"`
	Revenue  int     `json:"revenue"`
	Rating   float64 `json:"rating"`
}

type TripSummary struct {
	TotalTrips    int     `json:"total_trips"`
	TotalRevenue  int     `json:"total_revenue"`
	TotalDistance int     `json:"total_distance"`
	AverageRating float64 `json:"average_rating"`
	Car           Car     `json:"car"`
}
