package dto

type SlaRequest struct {
	ClientID int     `json:"client_id" binding:"required"`
	Priority int     `json:"priority" binding:"required"`
	Hours    float64 `json:"hours" binding:"required"`
}

type SlaResponse struct {
	ID       int     `json:"id"`
	ClientID int     `json:"client_id"`
	Priority int     `json:"priority"`
	Hours    float64 `json:"hours"`
}

type SlaUpdate struct {
	ID       int     `json:"id"`
	ClientID int     `json:"client_id"`
	Priority int     `json:"priority"`
	Hours    float64 `json:"hours"`
}
