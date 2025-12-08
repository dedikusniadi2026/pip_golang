package model

import "time"

type CarModel struct {
	ID        int       `json:"id"`
	ModelName string    `json:"model_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
