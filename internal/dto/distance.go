package dto

type DistanceRequest struct {
	Distance     float64 `json:"distance" binding:"required"`
	TicketNumber string  `json:"ticket_number" binding:"required"`
	ProviderId   int     `json:"provider_id" binding:"required"`
	ProviderName string  `json:"provider_name" binding:"required"`
}

type DistanceResponse struct {
	ID           int     `json:"id"`
	Distance     float64 `json:"distance"`
	TicketNumber string  `json:"ticket_number"`
	ProviderId   int     `json:"provider_id"`
	ProviderName string  `json:"provider_name"`
}
