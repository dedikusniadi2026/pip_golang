package model

import "time"

type PopularDestination struct {
	ID          int       `json:"id"`
	Destination string    `json:"destination"`
	Bookings    int       `json:"bookings"`
	CreatedAt   time.Time `json:"created_at"`
}
