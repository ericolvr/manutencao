package dto

import "time"

// TicketProblemRequest representa a requisição para associar um problema a um ticket
type TicketProblemRequest struct {
	ProblemID int `json:"problem_id" binding:"required"`
}

// TicketProblemResponse representa a resposta de um problema associado a um ticket
type TicketProblemResponse struct {
	ID        int       `json:"id"`
	TicketID  int       `json:"ticket_id"`
	ProblemID int       `json:"problem_id"`
	Problem   *ProblemResponse `json:"problem,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}


