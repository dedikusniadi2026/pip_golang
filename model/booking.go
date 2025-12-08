package model

import (
	"time"
)

type Booking struct {
	ID             string    `json:"id"`
	Customer       string    `json:"customer"`
	Driver         string    `json:"driver"`
	Place          string    `json:"place"`
	Date           string    `json:"date"`
	Price          string    `json:"price"`
	Status         string    `json:"status"`
	Payment        string    `json:"payment"`
	PhoneNumber    *string   `json:"phone_number"`
	PickupLocation *string   `json:"pickup_location"`
	DropLocation   *string   `json:"drop_location"`
	PickupTime     *string   `json:"pickup_time"`
	Amount         *float64  `json:"amount"`
	Notes          *string   `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
