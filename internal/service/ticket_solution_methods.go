package service

import (
	"context"
	"fmt"

	"github.com/ericolvr/maintenance-v2/internal/dto"
)

// AddSolutionToTicket associa uma solution do catálogo a um ticket
func (s *ticketService) AddSolutionToTicket(ctx context.Context, ticketID int, req *dto.TicketSolutionRequest) error {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Verificar se solution existe
	_, err = s.solutionRepo.FindByID(ctx, req.SolutionID)
	if err != nil {
		return fmt.Errorf("solution not found: %w", err)
	}

	// Associar solution ao ticket
	err = s.ticketRepo.AddSolutionToTicket(ctx, ticketID, req.SolutionID, req.Quantity)
	if err != nil {
		return fmt.Errorf("failed to add solution to ticket: %w", err)
	}

	return nil
}

// GetTicketSolutions retorna todas as solutions associadas a um ticket
func (s *ticketService) GetTicketSolutions(ctx context.Context, ticketID int) ([]dto.TicketSolutionResponse, error) {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}

	// Buscar solutions do ticket (via ticket_costs)
	ticketCosts, err := s.ticketRepo.GetTicketSolutions(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket solutions: %w", err)
	}

	// Converter para DTO com detalhes da solution
	var responses []dto.TicketSolutionResponse
	for _, cost := range ticketCosts {
		// Só incluir custos que têm solution_id (não são custos customizados)
		if cost.SolutionID == nil {
			continue
		}

		// Buscar detalhes da solution
		solution, err := s.solutionRepo.FindByID(ctx, *cost.SolutionID)
		if err != nil {
			// Se não encontrar a solution, continua sem os detalhes
			responses = append(responses, dto.TicketSolutionResponse{
				ID:         cost.ID,
				TicketID:   cost.TicketID,
				SolutionID: *cost.SolutionID,
				Quantity:   cost.Quantity,
				UnitPrice:  cost.UnitPrice,
				Subtotal:   cost.Subtotal,
				CreatedAt:  cost.CreatedAt,
			})
			continue
		}

		responses = append(responses, dto.TicketSolutionResponse{
			ID:         cost.ID,
			TicketID:   cost.TicketID,
			SolutionID: *cost.SolutionID,
			Solution: &dto.SolutionResponse{
				ID:          solution.ID,
				Name:        solution.Name,
				Description: solution.Description,
				ProblemID:   solution.ProblemID,
			},
			Quantity:  cost.Quantity,
			UnitPrice: cost.UnitPrice,
			Subtotal:  cost.Subtotal,
			CreatedAt: cost.CreatedAt,
		})
	}

	return responses, nil
}

// RemoveSolutionFromTicket remove a associação entre uma solution e um ticket
func (s *ticketService) RemoveSolutionFromTicket(ctx context.Context, ticketID int, solutionID int) error {
	// Verificar se ticket existe
	_, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Remover associação
	err = s.ticketRepo.RemoveSolutionFromTicket(ctx, ticketID, solutionID)
	if err != nil {
		return fmt.Errorf("failed to remove solution from ticket: %w", err)
	}

	return nil
}
