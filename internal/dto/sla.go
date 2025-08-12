package dto

type SlaRequest struct {
	ClientID int     `json:"client_id" binding:"required"`
	Priority int     `json:"priority" binding:"required"`
	Hours    float64 `json:"hours" binding:"required"`
}

type SlaResponse struct {
	ID         int     `json:"id"`
	ClientID   int     `json:"client_id"`
	ClientName string  `json:"client_name"`
	Priority   int     `json:"priority"`
	Hours      float64 `json:"hours"`
}

type SlaUpdate struct {
	ID       int     `json:"id"`
	ClientID int     `json:"client_id"`
	Priority int     `json:"priority"`
	Hours    float64 `json:"hours"`
}

// SlaWithClient representa um SLA com informações do cliente via JOIN
type SlaWithClient struct {
	ID         int     `json:"id" db:"id"`
	ClientID   int     `json:"client_id" db:"client_id"`
	ClientName string  `json:"client_name" db:"client_name"`
	Priority   int     `json:"priority" db:"priority"`
	Hours      float64 `json:"hours" db:"hours"`
}

type SlaJustHours struct {
	Hours float64 `json:"hours" db:"hours"`
}
