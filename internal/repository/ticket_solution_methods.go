package repository

import (
	"context"
	"fmt"

	"github.com/ericolvr/maintenance-v2/internal/domain"
)

// AddSolutionToTicket associa uma solution do catálogo a um ticket
func (r *ticketRepository) AddSolutionToTicket(ctx context.Context, ticketID int, solutionID int, quantity int) error {
	// Buscar dados da solution para calcular custo
	var unitPrice float64
	err := r.db.QueryRowContext(ctx, "SELECT unit_price FROM solutions WHERE id = $1", solutionID).Scan(&unitPrice)
	if err != nil {
		return fmt.Errorf("failed to get solution price: %w", err)
	}

	subtotal := float64(quantity) * unitPrice

	// Inserir na tabela ticket_costs
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO ticket_costs (ticket_id, solution_id, quantity, unit_price, subtotal) 
		VALUES ($1, $2, $3, $4, $5)
	`, ticketID, solutionID, quantity, unitPrice, subtotal)
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
