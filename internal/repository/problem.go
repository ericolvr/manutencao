package repository

import (
	"context"
	"database/sql"
	"github.com/ericolvr/maintenance-v2/internal/domain"
)

type ProblemRepository interface {
	Create(ctx context.Context, problem *domain.Problem) (int, error)
	FindByID(ctx context.Context, id int) (*domain.Problem, error)
	FindAll(ctx context.Context) ([]domain.Problem, error)
	Update(ctx context.Context, id int, problem *domain.Problem) error
	Delete(ctx context.Context, id int) error
}

type problemRepository struct {
	db *sql.DB
}

func NewProblemRepository(db *sql.DB) ProblemRepository {
	return &problemRepository{db: db}
}

func (r *problemRepository) Create(ctx context.Context, problem *domain.Problem) (int, error) {
	query := `
		INSERT INTO problems (name, description)
		VALUES ($1, $2)
		RETURNING id
	`
	
	var id int
	err := r.db.QueryRowContext(ctx, query, problem.Name, problem.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

func (r *problemRepository) FindByID(ctx context.Context, id int) (*domain.Problem, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM problems
		WHERE id = $1
	`
	
	problem := &domain.Problem{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&problem.ID,
		&problem.Name,
		&problem.Description,
		&problem.CreatedAt,
		&problem.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	
	return problem, nil
}

func (r *problemRepository) FindAll(ctx context.Context) ([]domain.Problem, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM problems
		ORDER BY name
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var problems []domain.Problem
	for rows.Next() {
		var problem domain.Problem
		err := rows.Scan(
			&problem.ID,
			&problem.Name,
			&problem.Description,
			&problem.CreatedAt,
			&problem.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}
	
	return problems, nil
}

func (r *problemRepository) Update(ctx context.Context, id int, problem *domain.Problem) error {
	query := `
		UPDATE problems
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	
	_, err := r.db.ExecContext(ctx, query, problem.Name, problem.Description, id)
	return err
}

func (r *problemRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM problems WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
