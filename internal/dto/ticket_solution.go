package dto

import "time"

// TicketSolutionRequest representa a requisição para associar uma solution a um ticket
type TicketSolutionRequest struct {
	SolutionID int `json:"solution_id" binding:"required"`
	Quantity   int `json:"quantity" binding:"required"`
}

// TicketSolutionResponse representa a resposta de uma solution associada a um ticket
type TicketSolutionResponse struct {
	ID         int              `json:"id"`
	TicketID   int              `json:"ticket_id"`
	SolutionID int              `json:"solution_id"`
	Solution   *SolutionResponse `json:"solution,omitempty"`
	Quantity   int              `json:"quantity"`
	UnitPrice  float64          `json:"unit_price"`
	Subtotal   float64          `json:"subtotal"`
	CreatedAt  time.Time        `json:"created_at"`
}


