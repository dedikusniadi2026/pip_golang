package model

type DashboardSummary struct {
	TotalTrips    int  `json:"totalTrips"`
	TotalDistance int  `json:"totalDistance"`
	TotalDuration int  `json:"totalDuration"`
	LastTrip      Trip `json:"lastTrip"`
}
