package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	_ "github.com/lib/pq"
)

var ErrCosthNotFound = errors.New("cost not found")

type CostRepository interface {
	List(ctx context.Context) ([]domain.Cost, error)
	FindByID(ctx context.Context, id int) (*domain.Cost, error)
	Update(ctx context.Context, cost *domain.Cost) error
	Delete(ctx context.Context, id int) error
}

type costRepository struct {
	db *sql.DB
}

func NewCostRepository(db *sql.DB) CostRepository {
	return &costRepository{
		db: db,
	}
}

func (r *costRepository) List(ctx context.Context) ([]domain.Cost, error) {
	query := `SELECT id, value_per_km, initial_value FROM costs`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing branchs: %w", err)
	}
	defer rows.Close()

	var costs []domain.Cost
	for rows.Next() {
		var cost domain.Cost
		if err := rows.Scan(
			&cost.ID,
			&cost.ValuePerKm,
			&cost.InitialValue); err != nil {
			return nil, fmt.Errorf("error scanning cost: %w", err)
		}
		costs = append(costs, cost)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating costs: %w", err)
	}

	return costs, nil
}

func (r *costRepository) FindByID(ctx context.Context, id int) (*domain.Cost, error) {
	query := `SELECT id, value_per_km, initial_value FROM costs WHERE id = $1`
	var cost domain.Cost

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cost.ID,
		&cost.ValuePerKm,
		&cost.InitialValue,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBranchNotFound
		}
		return nil, fmt.Errorf("error finding cost by id: %w", err)
	}

	return &cost, nil
}

func (r *costRepository) Update(ctx context.Context, cost *domain.Cost) error {
	query := `UPDATE costs SET 
			value_per_km = $1, 
			initial_value = $2, 
			WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query,
		cost.ValuePerKm,
		cost.InitialValue,
		cost.ID)
	if err != nil {
		return fmt.Errorf("error updating cost: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrBranchNotFound
	}

	return nil
}

func (r *costRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM costs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting branch: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrBranchNotFound
	}

	return nil
}
