package domain

import "time"

// Solution representa uma solução do catálogo
type Solution struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	UnitPrice   float64   `json:"unit_price" db:"unit_price"`
	ProblemID   int       `json:"problem_id" db:"problem_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
