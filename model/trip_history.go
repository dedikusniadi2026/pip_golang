package model

type TripHistory struct {
	ID              int     `json:"id"`
	BookingCode     string  `json:"booking_code"`
	CustomerName    string  `json:"customer_name"`
	BookingDate     string  `json:"booking_date"`
	DurationMinutes int     `json:"duration_minutes"`
	DistanceKM      int     `json:"distance_km"`
	PickupLocation  string  `json:"pickup_location"`
	Destination     string  `json:"destination"`
	DriverName      string  `json:"driver_name"`
	VehicleName     string  `json:"vehicle_name"`
	Amount          int     `json:"amount"`
	Rating          float32 `json:"rating"`
	Feedback        string  `json:"feedback"`
}
