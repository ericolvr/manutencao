package dto

import (
	"time"

	"github.com/ericolvr/maintenance-v2/internal/domain"
)

type TicketRequest struct {
	Number      string `json:"number" binding:"required"`
	Status      int64  `json:"status" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
	Description string `json:"description" binding:"required"`
	OpenDate    string `json:"open_date" binding:"required"` // Formato "2006-01-02T15:04:05Z"
	BranchID    int    `json:"branch_id" binding:"required"`
}

type UpdateTicketRequest struct {
	Number          string                   `json:"number" binding:"required"`
	Status          int64                    `json:"status" binding:"required"`
	Priority        string                   `json:"priority" binding:"required"`
	Description     string                   `json:"description" binding:"required"`
	OpenDate        string                   `json:"open_date" binding:"required"` // Formato "2006-01-02T15:04:05Z"
	CloseDate       *string                  `json:"close_date,omitempty"`         // Opcional, formato "2006-01-02T15:04:05Z"
	BranchID        int                      `json:"branch_id" binding:"required"`
	ProviderID      int                      `json:"provider_id"`
	SolutionItems   []SolutionItemRequest    `json:"solution_items,omitempty"`
}

type TicketResponse struct {
	ID              int                      `json:"id"`
	Number          string                   `json:"number"`
	Status          int                      `json:"status"`
	Priority        string                   `json:"priority"`
	Description     string                   `json:"description"`
	OpenDate        time.Time                `json:"open_date"`
	CloseDate       *time.Time               `json:"close_date,omitempty"`
	BranchID        int                      `json:"branch_id"`
	ProviderID      *int                     `json:"provider_id,omitempty"`
	Costs           []SolutionItemResponse   `json:"costs,omitempty"`
	TotalCost       float64                  `json:"total_cost"`
}

// SolutionItemRequest representa um item de solução na requisição
type SolutionItemRequest struct {
	Description string  `json:"description" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
}

// SolutionItemResponse representa um item de solução na resposta
type SolutionItemResponse struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
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
		Costs:       []SolutionItemResponse{}, // Será preenchido pelo service
		TotalCost:   0.0, // Será calculado pelo service
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
			ID:          cost.ID,
			Description: cost.SolutionName,
			UnitPrice:   cost.UnitPrice,
			Quantity:    cost.Quantity,
			Subtotal:    cost.Subtotal,
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
		Costs:       costItems,
		TotalCost:   totalCost,
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
