package domain

import "time"

type Distance struct {
	ID           int       `json:"id" db:"id"`
	Distance     float64   `json:"distance" db:"distance"`
	TicketNumber string    `json:"ticket_number" db:"ticket_number"`
	ProviderId   int       `json:"provider_id" db:"provider_id"`
	ProviderName string    `json:"provider_name" db:"provider_name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
