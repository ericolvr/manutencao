package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	_ "github.com/lib/pq"
)

var ErrSlaNotFound = errors.New("sla not found")

type SlaRepository interface {
	Create(ctx context.Context, sla *domain.Sla) (int, error)
	List(ctx context.Context) ([]domain.Sla, error)
	FindByID(ctx context.Context, id int) (*domain.Sla, error)
	GetByClient(ctx context.Context, client int) ([]domain.Sla, error)
	Update(ctx context.Context, sla *domain.Sla) error
	Delete(ctx context.Context, id int) error
}

type slaRepository struct {
	db *sql.DB
}

func NewSlaRepository(db *sql.DB) SlaRepository {
	return &slaRepository{
		db: db,
	}
}

func (r *slaRepository) Create(ctx context.Context, sla *domain.Sla) (int, error) {
	query := `INSERT INTO slas (client_id, priority, hours) 
			VALUES ($1, $2, $3) RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		sla.ClientID,
		sla.Priority,
		sla.Hours).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating sla: %w", err)
	}

	return id, nil
}

func (r *slaRepository) List(ctx context.Context) ([]domain.Sla, error) {
	query := `SELECT id, client_id, priority, hours FROM slas ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing slas: %w", err)
	}
	defer rows.Close()

	var slas []domain.Sla
	for rows.Next() {
		var sla domain.Sla
		if err := rows.Scan(
			&sla.ID,
			&sla.ClientID,
			&sla.Priority,
			&sla.Hours); err != nil {
			return nil, fmt.Errorf("error scanning sla: %w", err)
		}
		slas = append(slas, sla)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating slas: %w", err)
	}

	return slas, nil
}

func (r *slaRepository) FindByID(ctx context.Context, id int) (*domain.Sla, error) {
	query := `SELECT id, client_id, priority, hours FROM slas WHERE id = $1`
	var sla domain.Sla

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sla.ID,
		&sla.ClientID,
		&sla.Priority,
		&sla.Hours,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSlaNotFound
		}
		return nil, fmt.Errorf("error finding branch by id: %w", err)
	}

	return &sla, nil
}

func (r *slaRepository) GetByClient(ctx context.Context, client int) ([]domain.Sla, error) {
	query := `SELECT id, client_id, priority, hours FROM slas WHERE client_id = $1 ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query, client)
	if err != nil {
		return nil, fmt.Errorf("error getting slas by client: %w", err)
	}
	defer rows.Close()

	var slas []domain.Sla
	for rows.Next() {
		var sla domain.Sla
		if err := rows.Scan(
			&sla.ID,
			&sla.ClientID,
			&sla.Priority,
			&sla.Hours); err != nil {
			return nil, fmt.Errorf("error scanning sla: %w", err)
		}
		slas = append(slas, sla)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating slas: %w", err)
	}

	return slas, nil
}

func (r *slaRepository) Update(ctx context.Context, sla *domain.Sla) error {
	query := `UPDATE slas SET 
			client_id = $1, 
			priority = $2, 
			hours = $3, 
			WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query,
		sla.ClientID,
		sla.Priority,
		sla.Hours,
		sla.ID)
	if err != nil {
		return fmt.Errorf("error updating sla: %w", err)
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

func (r *slaRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM slas WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting sla: %w", err)
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
