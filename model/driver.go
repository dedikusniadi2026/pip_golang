package model

import "time"

type Driver struct {
	ID                  int       `json:"id"`
	Name                string    `json:"name"`
	Email               string    `json:"email"`
	Phone               string    `json:"phone"`
	Address             string    `json:"address"`
	DriverLicenseNumber string    `json:"driver_license_number"`
	CarModelID          string    `json:"car_model_id"`
	CarTypeID           string    `json:"car_type_id"`
	PlateNumber         string    `json:"plate_number"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
