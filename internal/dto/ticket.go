package dto

import (
	"time"

	"github.com/ericolvr/maintenance-v2/internal/domain"
)

type TicketRequest struct {
	Number      string `json:"number" binding:"required"`
	Status      int    `json:"status" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
	Description string `json:"description" binding:"required"`
	OpenDate    string `json:"open_date" binding:"required"` // Formato "2006-01-02T15:04:05Z"
	BranchID    int    `json:"branch_id" binding:"required"`
	// ClientNumber  string     `json:"client_number" db:"client_number"`         // numero de chamado do cliente
	// SupportPeople string     `json:"support_people" db:"support_people"`       // pessoa do suporte responsavel pelo atendimento
	// SlaEndDate    *time.Time `json:"sla_end_date,omitempty" db:"sla_end_date"` // Data de término do SLA
}

type UpdateTicketRequest struct {
	Number        string                `json:"number" binding:"required"`
	Status        int                   `json:"status" binding:"required"`
	Priority      string                `json:"priority" binding:"required"`
	Description   string                `json:"description" binding:"required"`
	OpenDate      string                `json:"open_date" binding:"required"` // Formato "2006-01-02T15:04:05Z"
	CloseDate     *string               `json:"close_date,omitempty"`         // Opcional, formato "2006-01-02T15:04:05Z"
	BranchID      int                   `json:"branch_id" binding:"required"`
	ProviderID    int                   `json:"provider_id"`
	SolutionItems []SolutionItemRequest `json:"solution_items,omitempty"`
}

type TicketResponse struct {
	ID           int                    `json:"id"`
	Number       string                 `json:"number"`
	Status       int                    `json:"status"`
	Priority     string                 `json:"priority"`
	Description  string                 `json:"description"`
	OpenDate     time.Time              `json:"open_date"`
	CloseDate    *time.Time             `json:"close_date,omitempty"`
	BranchID     int                    `json:"branch_id"`
	BranchName   string                 `json:"branch_name"`
	BranchUniorg string                 `json:"branch_uniorg"`
	ProviderID   *int                   `json:"provider_id,omitempty"`
	ProviderName *string                `json:"provider_name,omitempty"`
	Distance     *float64               `json:"distance,omitempty"`
	Costs        []SolutionItemResponse `json:"costs,omitempty"`
	TotalCost    float64                `json:"total_cost"`
}

// SolutionItemRequest representa um item de solução na requisição
type SolutionItemRequest struct {
	Description string  `json:"description" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
}

// SolutionItemResponse representa um item de solução na resposta
type SolutionItemResponse struct {
	ProblemName  string  `json:"problem_name"`
	SolutionName string  `json:"solution_name"`
	UnitPrice    float64 `json:"unit_price"`
	Subtotal     float64 `json:"subtotal"`
}

// MapToTicketResponse mapeia um domínio Ticket para sua representação DTO TicketResponse
func ToTicketResponse(ticket *domain.Ticket) *TicketResponse {
	if ticket == nil {
		return nil
	}

	return &TicketResponse{
		ID:          ticket.ID,
		Number:      ticket.Number,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		Description: ticket.Description,
		OpenDate:    ticket.OpenDate,
		CloseDate:   ticket.CloseDate,
		BranchID:    ticket.BranchID,
		ProviderID:  ticket.ProviderID,
		Distance:    nil,                      // Será preenchido pelo service
		Costs:       []SolutionItemResponse{}, // Será preenchido pelo service
		TotalCost:   0.0,                      // Será calculado pelo service
	}
}

// ToTicketResponseWithCosts mapeia um domínio Ticket com seus custos para TicketResponse
func ToTicketResponseWithCosts(ticket *domain.Ticket, costs []domain.TicketCost) *TicketResponse {
	if ticket == nil {
		return nil
	}

	var costItems []SolutionItemResponse
	var totalCost float64

	for _, cost := range costs {
		costItem := SolutionItemResponse{
			ProblemName:  cost.ProblemName,
			SolutionName: cost.SolutionName,
			UnitPrice:    cost.UnitPrice,
			Subtotal:     cost.Subtotal,
		}
		costItems = append(costItems, costItem)
		totalCost += cost.Subtotal
	}

	return &TicketResponse{
		ID:          ticket.ID,
		Number:      ticket.Number,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		Description: ticket.Description,
		OpenDate:    ticket.OpenDate,
		CloseDate:   ticket.CloseDate,
		BranchID:    ticket.BranchID,
		ProviderID:  ticket.ProviderID,
		Distance:    nil, // Será preenchido pelo service
		Costs:       costItems,
		TotalCost:   totalCost,
	}
}

// ToTicketResponseWithDistanceAndCosts mapeia um domínio Ticket com distance e custos para TicketResponse
func ToTicketResponseWithDistanceAndCosts(ticket *domain.Ticket, distance *float64, costs []domain.TicketCost) *TicketResponse {
	if ticket == nil {
		return nil
	}

	var costItems []SolutionItemResponse
	var totalCost float64

	for _, cost := range costs {
		costItem := SolutionItemResponse{
			ProblemName:  cost.ProblemName,
			SolutionName: cost.SolutionName,
			UnitPrice:    cost.UnitPrice,
			Subtotal:     cost.Subtotal,
		}
		costItems = append(costItems, costItem)
		totalCost += cost.Subtotal
	}

	return &TicketResponse{
		ID:           ticket.ID,
		Number:       ticket.Number,
		Status:       ticket.Status,
		Priority:     ticket.Priority,
		Description:  ticket.Description,
		OpenDate:     ticket.OpenDate,
		CloseDate:    ticket.CloseDate,
		BranchID:     ticket.BranchID,
		BranchName:   "", // Será preenchido pela nova função
		BranchUniorg: "", // Será preenchido pela nova função
		ProviderID:   ticket.ProviderID,
		ProviderName: nil, // Será preenchido pela nova função
		Distance:     distance,
		Costs:        costItems,
		TotalCost:    totalCost,
	}
}

// AddProviderRequest contém o ID do provider a ser adicionado ao ticket
type AddProviderRequest struct {
	ProviderID int `json:"provider_id" binding:"required"`
}

// AddSolutionItemRequest contém dados para adicionar item de solução ao ticket
type AddSolutionItemRequest struct {
	Description string  `json:"description" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
}

// UpdateSolutionItemRequest contém dados para atualizar item de solução
type UpdateSolutionItemRequest struct {
	Description string  `json:"description" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
}

// ToTicketResponseWithBranchProviderDistanceAndCosts mapeia um domínio Ticket com branch, provider, distance e custos para TicketResponse
func ToTicketResponseWithBranchProviderDistanceAndCosts(ticket *domain.Ticket, branch *domain.Branch, provider *domain.Provider, distance *float64, costs []domain.TicketCost) *TicketResponse {
	if ticket == nil {
		return nil
	}

	var costItems []SolutionItemResponse
	var totalCost float64

	for _, cost := range costs {
		costItem := SolutionItemResponse{
			ProblemName:  cost.ProblemName,
			SolutionName: cost.SolutionName,
			UnitPrice:    cost.UnitPrice,
			Subtotal:     cost.Subtotal,
		}
		costItems = append(costItems, costItem)
		totalCost += cost.Subtotal
	}

	// Preparar nome do provider (se existir)
	var providerName *string
	if provider != nil {
		providerName = &provider.Name
	}

	// Preparar informações do branch
	branchName := ""
	branchUniorg := ""
	if branch != nil {
		branchName = branch.Name
		branchUniorg = branch.Uniorg
	}

	return &TicketResponse{
		ID:           ticket.ID,
		Number:       ticket.Number,
		Status:       ticket.Status,
		Priority:     ticket.Priority,
		Description:  ticket.Description,
		OpenDate:     ticket.OpenDate,
		CloseDate:    ticket.CloseDate,
		BranchID:     ticket.BranchID,
		BranchName:   branchName,
		BranchUniorg: branchUniorg,
		ProviderID:   ticket.ProviderID,
		ProviderName: providerName,
		Distance:     distance,
		Costs:        costItems,
		TotalCost:    totalCost,
	}
}

// TicketWithDetails representa um ticket com todos os dados relacionados em uma única query
type TicketWithDetails struct {
	ID          int        `json:"id" db:"id"`
	Number      string     `json:"number" db:"number"`
	Status      int        `json:"status" db:"status"`
	Priority    string     `json:"priority" db:"priority"`
	Description string     `json:"description" db:"description"`
	OpenDate    time.Time  `json:"open_date" db:"open_date"`
	CloseDate   *time.Time `json:"close_date,omitempty" db:"close_date"`
	BranchID    int        `json:"branch_id" db:"branch_id"`
	ProviderID  *int       `json:"provider_id,omitempty" db:"provider_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// Dados relacionados via JOIN
	BranchName   string   `json:"branch_name" db:"branch_name"`
	BranchUniorg string   `json:"branch_uniorg" db:"branch_uniorg"`
	ProviderName *string  `json:"provider_name,omitempty" db:"provider_name"`
	Distance     *float64 `json:"distance,omitempty" db:"distance"`
}

// ToTicketResponseFromDetails converte TicketWithDetails para TicketResponse com custos
func ToTicketResponseFromDetails(ticketDetail *TicketWithDetails, costs []domain.TicketCost) *TicketResponse {
	if ticketDetail == nil {
		return nil
	}

	var costItems []SolutionItemResponse
	var totalCost float64

	for _, cost := range costs {
		costItem := SolutionItemResponse{
			ProblemName:  cost.ProblemName,
			SolutionName: cost.SolutionName,
			UnitPrice:    cost.UnitPrice,
			Subtotal:     cost.Subtotal,
		}
		costItems = append(costItems, costItem)
		totalCost += cost.Subtotal
	}

	return &TicketResponse{
		ID:           ticketDetail.ID,
		Number:       ticketDetail.Number,
		Status:       ticketDetail.Status,
		Priority:     ticketDetail.Priority,
		Description:  ticketDetail.Description,
		OpenDate:     ticketDetail.OpenDate,
		CloseDate:    ticketDetail.CloseDate,
		BranchID:     ticketDetail.BranchID,
		BranchName:   ticketDetail.BranchName,
		BranchUniorg: ticketDetail.BranchUniorg,
		ProviderID:   ticketDetail.ProviderID,
		ProviderName: ticketDetail.ProviderName,
		Distance:     ticketDetail.Distance,
		Costs:        costItems,
		TotalCost:    totalCost,
	}
}
