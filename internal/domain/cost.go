package domain

import "time"

type Cost struct {
	ID           int       `json:"id" db:"id"`
	ValuePerKm   float64   `json:"value_per_km" db:"value_per_km"`
	InitialValue float64   `json:"initial_value" db:"initial_value"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
