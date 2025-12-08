package model

type BookingTrend struct {
	Month   string `json:"month"`
	Booking int    `json:"booking"`
	Year    int    `json:"year"`
}
