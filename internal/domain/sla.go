package domain

type Sla struct {
	ID       int     `json:"id"`
	ClientID int     `json:"client_id"`
	Priority int     `json:"priority"`
	Hours    float64 `json:"hours"`
}
