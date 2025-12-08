package model

import "time"

type Pdf struct {
	ID              string    `json:"trip_id"`
	CustomerName    string    `json:"customer_name"`
	BookingDate     time.Time `json:"booking_date"`
	DurationMinutes int       `json:"duration_minutes"`
	DistanceKM      int       `json:"distance_km"`
	PickupLocation  string    `json:"pickup_location"`
	Destination     string    `json:"destination"`
	DriverName      string    `json:"driver_name"`
	VehicleName     string    `json:"vehicle_name"`
	Amount          int       `json:"amount"`
	Rating          float32   `json:"rating"`
	Feedback        string    `json:"feedback"`
}

type PDFTemplateData struct {
	Pdf
	AmountFormatted string
}
