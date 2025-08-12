package domain

import "time"

// Ticket representa um ticket de manutenção
type Ticket struct {
	ID          int        `json:"id" db:"id"`
	Number      string     `json:"number" db:"number"`
	Status      int        `json:"status" db:"status"`
	Priority    string     `json:"priority" db:"priority"`
	Description string     `json:"description" db:"description"`
	OpenDate    time.Time  `json:"open_date" db:"open_date"`
	CloseDate   *time.Time `json:"close_date,omitempty" db:"close_date"`
	BranchID    int        `json:"branch_id" db:"branch_id"`
	ProviderID  *int       `json:"provider_id,omitempty" db:"provider_id"`

	// new fields (RC)
	// ClientNumber  string     `json:"client_number" db:"client_number"`         // numero de chamado do cliente
	// SupportPeople string     `json:"support_people" db:"support_people"`       // pessoa do suporte responsavel pelo atendimento
	// SlaEndDate    *time.Time `json:"sla_end_date,omitempty" db:"sla_end_date"` // Data de término do SLA

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TicketCost representa os custos aplicados a um ticket
type TicketCost struct {
	ID           int       `json:"id" db:"id"`
	TicketID     int       `json:"ticket_id" db:"ticket_id"`
	ProblemID    *int      `json:"problem_id,omitempty" db:"problem_id"`
	ProblemName  string    `json:"problem_name" db:"problem_name"`
	SolutionID   *int      `json:"solution_id,omitempty" db:"solution_id"`
	SolutionName string    `json:"solution_name" db:"solution_name"`
	Quantity     int       `json:"quantity" db:"quantity"`
	UnitPrice    float64   `json:"unit_price" db:"unit_price"`
	Subtotal     float64   `json:"subtotal" db:"subtotal"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// TicketProblem representa a relação many-to-many entre Ticket e Problem
type TicketProblem struct {
	ID        int       `json:"id" db:"id"`
	TicketID  int       `json:"ticket_id" db:"ticket_id"`
	ProblemID int       `json:"problem_id" db:"problem_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// STATUS
// 1. Agendado
// 2. Aguarda Atendimento
// 3. Aguarda Faturamento
// 4. Aguarda Finalização
// 5. Compras
// 6. Concluído
// 7. Em Atendimento
// 8. Enviado
// 9. Equipamento Entregue
// 10. Estoque
// 11. Logística Reversa
// 12. Novo
// 13. Prestação de Contas
// 14. Contas Não Aprovadas
// 15. Emitir Nota
