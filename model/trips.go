package model

import "time"

type VehicleTrip struct {
	ID            uint `gorm:"primaryKey"`
	VehicleID     uint
	DriverID      uint
	TripDate      time.Time
	Origin        string
	Destination   string
	DistanceKM    int
	Rating        int
	Price         float64
	PassengerName string
}

type TotalTrips struct {
	TotalTrips    int64   `json:"total_trips"`
	TotalDistance int64   `json:"total_distance"`
	TotalRevenue  float64 `json:"total_revenue"`
	AverageRating float64 `json:"average_rating"`
}
