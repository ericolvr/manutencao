package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/dto"
)

type TicketRepository interface {
	Create(ctx context.Context, ticket *domain.Ticket) (int, error)
	List(ctx context.Context, limit, offset int) ([]domain.Ticket, int, error)
	ListWithDetails(ctx context.Context, limit, offset int) ([]dto.TicketWithDetails, int, error)
	FindByID(ctx context.Context, id int) (*domain.Ticket, error)
	Update(ctx context.Context, ticket *domain.Ticket) error
	Delete(ctx context.Context, ticketID int) error
	GetTicketNumber(ctx context.Context) (int, error)
	AddProvider(ctx context.Context, ticketID int, providerID int) error
	RemoveProvider(ctx context.Context, ticketID int) error
	GetProviderOnTicket(ctx context.Context, ticketID int) (*domain.Provider, error)

	// Ticket Costs methods
	CreateTicketCosts(ctx context.Context, ticketID int, costs []domain.TicketCost) error
	GetTicketCosts(ctx context.Context, ticketID int) ([]domain.TicketCost, error)
	UpdateTicketCosts(ctx context.Context, ticketID int, costs []domain.TicketCost) error
	DeleteTicketCosts(ctx context.Context, ticketID int) error
	// Ticket Problems methods

	AddProblemToTicket(ctx context.Context, ticketID int, problemID int) error
	GetTicketProblems(ctx context.Context, ticketID int) ([]domain.TicketProblem, error)
	RemoveProblemFromTicket(ctx context.Context, ticketID int, problemID int) error

	// Ticket Solutions methods
	AddSolutionToTicket(ctx context.Context, ticketID int, solutionID int, quantity int) error
	GetTicketSolutions(ctx context.Context, ticketID int) ([]domain.TicketCost, error)
	RemoveSolutionFromTicket(ctx context.Context, ticketID int, solutionID int) error
}

type ticketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *domain.Ticket) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Log rollback error
			}
		}
	}()

	var ticketID int
	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO tickets (
			number, status, priority, description, open_date, close_date, branch_id, provider_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		ticket.Number,
		ticket.Status,
		ticket.Priority,
		ticket.Description,
		ticket.OpenDate,
		ticket.CloseDate,
		ticket.BranchID,
		ticket.ProviderID,
	).Scan(&ticketID)

	if err != nil {
		return 0, fmt.Errorf("failed to create ticket: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ticketID, nil
}

func (r *ticketRepository) List(ctx context.Context, limit, offset int) ([]domain.Ticket, int, error) {
	var records int

	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tickets").Scan(&records)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	query := `
		SELECT id, number, status, priority, description, open_date, close_date, branch_id, provider_id
		FROM tickets
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tickets: %w", err)
	}
	defer rows.Close()

	var tickets []domain.Ticket
	for rows.Next() {
		var ticket domain.Ticket
		var providerID sql.NullInt64
		if err := rows.Scan(
			&ticket.ID,
			&ticket.Number,
			&ticket.Status,
			&ticket.Priority,
			&ticket.Description,
			&ticket.OpenDate,
			&ticket.CloseDate,
			&ticket.BranchID,
			&providerID,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning ticket: %w", err)
		}

		if providerID.Valid {
			providerIDValue := int(providerID.Int64)
			ticket.ProviderID = &providerIDValue
		}

		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating tickets: %w", err)
	}

	return tickets, records, nil
}

func (r *ticketRepository) ListWithDetails(ctx context.Context, limit, offset int) ([]dto.TicketWithDetails, int, error) {
	var records int

	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tickets").Scan(&records)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	query := `
		SELECT 
			t.id, t.number, t.status, t.priority, t.description, 
			t.open_date, t.close_date, t.branch_id, t.provider_id,
			t.created_at, t.updated_at,
			b.name as branch_name, 
			b.uniorg as branch_uniorg,
			p.name as provider_name,
			d.distance
		FROM tickets t
		LEFT JOIN branches b ON t.branch_id = b.id
		LEFT JOIN providers p ON t.provider_id = p.id  
		LEFT JOIN distances d ON t.number = d.number
		ORDER BY t.id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tickets with details: %w", err)
	}
	defer rows.Close()

	var tickets []dto.TicketWithDetails
	for rows.Next() {
		var ticket dto.TicketWithDetails
		var providerID sql.NullInt64
		var providerName sql.NullString
		var distance sql.NullFloat64
		
		if err := rows.Scan(
			&ticket.ID,
			&ticket.Number,
			&ticket.Status,
			&ticket.Priority,
			&ticket.Description,
			&ticket.OpenDate,
			&ticket.CloseDate,
			&ticket.BranchID,
			&providerID,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.BranchName,
			&ticket.BranchUniorg,
			&providerName,
			&distance,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning ticket with details: %w", err)
		}

		// Handle nullable fields
		if providerID.Valid {
			providerIDValue := int(providerID.Int64)
			ticket.ProviderID = &providerIDValue
		}
		
		if providerName.Valid {
			ticket.ProviderName = &providerName.String
		}
		
		if distance.Valid {
			ticket.Distance = &distance.Float64
		}

		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating tickets with details: %w", err)
	}

	return tickets, records, nil
}

func (r *ticketRepository) FindByID(ctx context.Context, id int) (*domain.Ticket, error) {
	ticket := domain.Ticket{}

	var closeDate sql.NullTime
	var providerID sql.NullInt64

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, number, status, priority, description, open_date, close_date, branch_id, provider_id
		FROM tickets
		WHERE id = $1`,
		id,
	).Scan(
		&ticket.ID,
		&ticket.Number,
		&ticket.Status,
		&ticket.Priority,
		&ticket.Description,
		&ticket.OpenDate,
		&closeDate,
		&ticket.BranchID,
		&providerID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ticket with id %d not found", id)
		}
		return nil, fmt.Errorf("error scanning ticket: %w", err)
	}

	if closeDate.Valid {
		ticket.CloseDate = &closeDate.Time
	}

	if providerID.Valid {
		providerIDValue := int(providerID.Int64)
		ticket.ProviderID = &providerIDValue
	}

	// Provider ID já foi atribuído acima
	// Branch ID já está no ticket, não precisamos carregar o objeto completo

	return &ticket, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE tickets SET 
			number = $1, status = $2, priority = $3, description = $4, 
			open_date = $5, close_date = $6, branch_id = $7, provider_id = $8
		WHERE id = $9`,
		ticket.Number,
		ticket.Status,
		ticket.Priority,
		ticket.Description,
		ticket.OpenDate,
		ticket.CloseDate,
		ticket.BranchID,
		ticket.ProviderID,
		ticket.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update ticket: %w", err)
	}
	return nil
}

func (r *ticketRepository) Delete(ctx context.Context, ticketID int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM tickets WHERE id = $1", ticketID)
	if err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ticket with id %d not found", ticketID)
	}
	return nil
}

func (r *ticketRepository) GetTicketNumber(ctx context.Context) (int, error) {
	var maxNumber sql.NullString
	err := r.db.QueryRowContext(ctx, "SELECT MAX(number) FROM tickets").Scan(&maxNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to get max ticket number: %w", err)
	}

	// Se não há tickets ainda, começar com 1
	if !maxNumber.Valid || maxNumber.String == "" {
		return 1, nil
	}

	// Converter string para int
	maxInt, err := strconv.Atoi(maxNumber.String)
	if err != nil {
		return 0, fmt.Errorf("failed to convert max number to int: %w", err)
	}

	return maxInt + 1, nil
}

func (r *ticketRepository) AddProvider(ctx context.Context, ticketID int, providerID int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE tickets SET provider_id = $1 WHERE id = $2`,
		providerID,
		ticketID,
	)
	if err != nil {
		return fmt.Errorf("failed to add provider to ticket: %w", err)
	}
	return nil
}

func (r *ticketRepository) RemoveProvider(ctx context.Context, ticketID int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE tickets SET provider_id = NULL WHERE id = $1`,
		ticketID,
	)
	if err != nil {
		return fmt.Errorf("failed to remove provider from ticket: %w", err)
	}
	return nil
}

func (r *ticketRepository) GetProviderOnTicket(ctx context.Context, ticketID int) (*domain.Provider, error) {
	var provider domain.Provider
	err := r.db.QueryRowContext(
		ctx,
		`SELECT p.id, p.name, p.mobile, p.zipcode, p.state, p.city, p.neighborhood, p.address, p.complement
		FROM providers p
		JOIN tickets t ON p.id = t.provider_id
		WHERE t.id = $1`,
		ticketID,
	).Scan(
		&provider.ID,
		&provider.Name,
		&provider.Mobile,
		&provider.Zipcode,
		&provider.State,
		&provider.City,
		&provider.Neighborhood,
		&provider.Address,
		&provider.Complement,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No provider associated
		}
		return nil, fmt.Errorf("error getting provider: %w", err)
	}
	return &provider, nil
}

func (r *ticketRepository) getBranchByID(ctx context.Context, branchID int) (*domain.Branch, error) {
	var branch domain.Branch
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, client, name, uniorg, zipcode, state, city, neighborhood, address, complement
		FROM branches
		WHERE id = $1`,
		branchID,
	).Scan(
		&branch.ID,
		&branch.Client,
		&branch.Name,
		&branch.Uniorg,
		&branch.Zipcode,
		&branch.State,
		&branch.City,
		&branch.Neighborhood,
		&branch.Address,
		&branch.Complement,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting branch: %w", err)
	}
	return &branch, nil
}

// CreateTicketCosts insere custos para um ticket
func (r *ticketRepository) CreateTicketCosts(ctx context.Context, ticketID int, costs []domain.TicketCost) error {
	if len(costs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, cost := range costs {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO ticket_costs (ticket_id, problem_id, problem_name, solution_id, solution_name, quantity, unit_price, subtotal) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			ticketID, cost.ProblemID, cost.ProblemName, cost.SolutionID, cost.SolutionName, cost.Quantity, cost.UnitPrice, cost.Subtotal)
		if err != nil {
			return fmt.Errorf("failed to insert ticket cost: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTicketCosts retorna todos os custos de um ticket
func (r *ticketRepository) GetTicketCosts(ctx context.Context, ticketID int) ([]domain.TicketCost, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, ticket_id, problem_id, problem_name, solution_id, solution_name, quantity, unit_price, subtotal, created_at 
		 FROM ticket_costs WHERE ticket_id = $1 ORDER BY created_at`,
		ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to query ticket costs: %w", err)
	}
	defer rows.Close()

	var costs []domain.TicketCost
	for rows.Next() {
		var cost domain.TicketCost
		err := rows.Scan(
			&cost.ID,
			&cost.TicketID,
			&cost.ProblemID,
			&cost.ProblemName,
			&cost.SolutionID,
			&cost.SolutionName,
			&cost.Quantity,
			&cost.UnitPrice,
			&cost.Subtotal,
			&cost.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket cost: %w", err)
		}
		costs = append(costs, cost)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ticket costs: %w", err)
	}

	return costs, nil
}

// UpdateTicketCosts atualiza os custos de um ticket (remove antigos e insere novos)
func (r *ticketRepository) UpdateTicketCosts(ctx context.Context, ticketID int, costs []domain.TicketCost) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Remove custos antigos
	_, err = tx.ExecContext(ctx, "DELETE FROM ticket_costs WHERE ticket_id = $1", ticketID)
	if err != nil {
		return fmt.Errorf("failed to delete old ticket costs: %w", err)
	}

	// Insere novos custos
	for _, cost := range costs {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO ticket_costs (ticket_id, problem_id, problem_name, solution_id, solution_name, quantity, unit_price, subtotal) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			ticketID, cost.ProblemID, cost.ProblemName, cost.SolutionID, cost.SolutionName, cost.Quantity, cost.UnitPrice, cost.Subtotal)
		if err != nil {
			return fmt.Errorf("failed to insert ticket cost: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteTicketCosts remove todos os custos de um ticket
func (r *ticketRepository) DeleteTicketCosts(ctx context.Context, ticketID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM ticket_costs WHERE ticket_id = $1", ticketID)
	if err != nil {
		return fmt.Errorf("failed to delete ticket costs: %w", err)
	}
	return nil
}

// AddProblemToTicket associa um problema a um ticket
func (r *ticketRepository) AddProblemToTicket(ctx context.Context, ticketID int, problemID int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO ticket_problems (ticket_id, problem_id) 
		VALUES ($1, $2)
		ON CONFLICT (ticket_id, problem_id) DO NOTHING
	`, ticketID, problemID)
	if err != nil {
		return fmt.Errorf("failed to add problem to ticket: %w", err)
	}
	return nil
}

// GetTicketProblems retorna todos os problemas associados a um ticket
func (r *ticketRepository) GetTicketProblems(ctx context.Context, ticketID int) ([]domain.TicketProblem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT tp.id, tp.ticket_id, tp.problem_id, tp.created_at
		FROM ticket_problems tp
		WHERE tp.ticket_id = $1
		ORDER BY tp.created_at DESC
	`, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket problems: %w", err)
	}
	defer rows.Close()

	var problems []domain.TicketProblem
	for rows.Next() {
		var problem domain.TicketProblem
		err := rows.Scan(&problem.ID, &problem.TicketID, &problem.ProblemID, &problem.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket problem: %w", err)
		}
		problems = append(problems, problem)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate ticket problems: %w", err)
	}

	return problems, nil
}

// RemoveProblemFromTicket remove a associação entre um problema e um ticket
func (r *ticketRepository) RemoveProblemFromTicket(ctx context.Context, ticketID int, problemID int) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM ticket_problems 
		WHERE ticket_id = $1 AND problem_id = $2
	`, ticketID, problemID)
	if err != nil {
		return fmt.Errorf("failed to remove problem from ticket: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("problem not found in ticket")
	}

	return nil
}
