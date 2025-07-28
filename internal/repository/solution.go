package repository

import (
	"context"
	"database/sql"
	"github.com/ericolvr/maintenance-v2/internal/domain"
)

type SolutionRepository interface {
	Create(ctx context.Context, solution *domain.Solution) (int, error)
	FindByID(ctx context.Context, id int) (*domain.Solution, error)
	FindAll(ctx context.Context) ([]domain.Solution, error)
	FindByProblem(ctx context.Context, problemID int) ([]domain.Solution, error)
	Update(ctx context.Context, id int, solution *domain.Solution) error
	Delete(ctx context.Context, id int) error
}

type solutionRepository struct {
	db *sql.DB
}

func NewSolutionRepository(db *sql.DB) SolutionRepository {
	return &solutionRepository{db: db}
}

func (r *solutionRepository) Create(ctx context.Context, solution *domain.Solution) (int, error) {
	query := `
		INSERT INTO solutions (name, description, unit_price, problem_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	
	var id int
	err := r.db.QueryRowContext(ctx, query, solution.Name, solution.Description, solution.UnitPrice, solution.ProblemID).Scan(&id)
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

func (r *solutionRepository) FindByID(ctx context.Context, id int) (*domain.Solution, error) {
	query := `
		SELECT id, name, description, unit_price, problem_id, created_at, updated_at
		FROM solutions
		WHERE id = $1
	`
	
	solution := &domain.Solution{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&solution.ID,
		&solution.Name,
		&solution.Description,
		&solution.UnitPrice,
		&solution.ProblemID,
		&solution.CreatedAt,
		&solution.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	
	return solution, nil
}

func (r *solutionRepository) FindAll(ctx context.Context) ([]domain.Solution, error) {
	query := `
		SELECT id, name, description, unit_price, problem_id, created_at, updated_at
		FROM solutions
		ORDER BY name
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var solutions []domain.Solution
	for rows.Next() {
		var solution domain.Solution
		err := rows.Scan(
			&solution.ID,
			&solution.Name,
			&solution.Description,
			&solution.UnitPrice,
			&solution.ProblemID,
			&solution.CreatedAt,
			&solution.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		solutions = append(solutions, solution)
	}
	
	return solutions, nil
}

func (r *solutionRepository) FindByProblem(ctx context.Context, problemID int) ([]domain.Solution, error) {
	query := `
		SELECT id, name, description, unit_price, problem_id, created_at, updated_at
		FROM solutions
		WHERE problem_id = $1
		ORDER BY name
	`
	
	rows, err := r.db.QueryContext(ctx, query, problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var solutions []domain.Solution
	for rows.Next() {
		var solution domain.Solution
		err := rows.Scan(
			&solution.ID,
			&solution.Name,
			&solution.Description,
			&solution.UnitPrice,
			&solution.ProblemID,
			&solution.CreatedAt,
			&solution.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		solutions = append(solutions, solution)
	}
	
	return solutions, nil
}

func (r *solutionRepository) Update(ctx context.Context, id int, solution *domain.Solution) error {
	query := `
		UPDATE solutions
		SET name = $1, description = $2, unit_price = $3, problem_id = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`
	
	_, err := r.db.ExecContext(ctx, query, solution.Name, solution.Description, solution.UnitPrice, solution.ProblemID, id)
	return err
}

func (r *solutionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM solutions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
