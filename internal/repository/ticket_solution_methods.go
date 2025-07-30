package repository

import (
	"context"
	"fmt"

	"github.com/ericolvr/maintenance-v2/internal/domain"
)

// AddSolutionToTicket associa uma solution do catálogo a um ticket
func (r *ticketRepository) AddSolutionToTicket(ctx context.Context, ticketID int, solutionID int, quantity int) error {
	// Buscar dados da solution e problema para calcular custo
	var unitPrice float64
	var solutionName string
	var problemID int
	var problemName string
	err := r.db.QueryRowContext(ctx, `
		SELECT s.name, s.unit_price, s.problem_id, p.name 
		FROM solutions s 
		JOIN problems p ON s.problem_id = p.id 
		WHERE s.id = $1
	`, solutionID).Scan(&solutionName, &unitPrice, &problemID, &problemName)
	if err != nil {
		return fmt.Errorf("failed to get solution data: %w", err)
	}

	subtotal := float64(quantity) * unitPrice

	// Inserir na tabela ticket_costs
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO ticket_costs (ticket_id, problem_id, problem_name, solution_id, solution_name, quantity, unit_price, subtotal) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, ticketID, problemID, problemName, solutionID, solutionName, quantity, unitPrice, subtotal)
	if err != nil {
		return fmt.Errorf("failed to add solution to ticket: %w", err)
	}
	return nil
}

// GetTicketSolutions retorna todas as solutions associadas a um ticket
func (r *ticketRepository) GetTicketSolutions(ctx context.Context, ticketID int) ([]domain.TicketCost, error) {
	// Reutilizar método existente GetTicketCosts
	return r.GetTicketCosts(ctx, ticketID)
}

// RemoveSolutionFromTicket remove uma solution de um ticket
func (r *ticketRepository) RemoveSolutionFromTicket(ctx context.Context, ticketID int, solutionID int) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM ticket_costs 
		WHERE ticket_id = $1 AND solution_id = $2
	`, ticketID, solutionID)
	if err != nil {
		return fmt.Errorf("failed to remove solution from ticket: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("solution not found in ticket")
	}

	return nil
}
