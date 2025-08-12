package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	"github.com/ericolvr/maintenance-v2/internal/dto"
	_ "github.com/lib/pq"
)

var ErrSlaNotFound = errors.New("sla not found")

type SlaRepository interface {
	Create(ctx context.Context, sla *domain.Sla) (int, error)
	List(ctx context.Context) ([]domain.Sla, error)
	ListWithClientNames(ctx context.Context) ([]dto.SlaWithClient, error)
	FindByID(ctx context.Context, id int) (*domain.Sla, error)
	FindByParams(ctx context.Context, client_id int, priority int) (*dto.SlaJustHours, error)
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
	query := `SELECT id, client_id, priority, hours FROM slas ORDER BY id ASC`

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

func (r *slaRepository) ListWithClientNames(ctx context.Context) ([]dto.SlaWithClient, error) {
	query := `
		SELECT 
			s.id, 
			s.client_id, 
			c.name as client_name,
			s.priority, 
			s.hours 
		FROM slas s
		LEFT JOIN clients c ON s.client_id = c.id
		ORDER BY s.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing slas with client names: %w", err)
	}
	defer rows.Close()

	var slas []dto.SlaWithClient
	for rows.Next() {
		var sla dto.SlaWithClient
		if err := rows.Scan(
			&sla.ID,
			&sla.ClientID,
			&sla.ClientName,
			&sla.Priority,
			&sla.Hours); err != nil {
			return nil, fmt.Errorf("error scanning sla with client: %w", err)
		}
		slas = append(slas, sla)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating slas with client names: %w", err)
	}

	return slas, nil
}

func (r *slaRepository) FindByID(ctx context.Context, id int) (*domain.Sla, error) {
	query := `SELECT id, client_id, priority, hours FROM sla WHERE id = $1`
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
		return nil, fmt.Errorf("error finding sla by id: %w", err)
	}

	return &sla, nil
}

func (r *slaRepository) FindByParams(ctx context.Context, client_id int, priority int) (*dto.SlaJustHours, error) {
	query := `SELECT hours FROM slas WHERE client_id = $1 AND priority = $2`
	var slaHours dto.SlaJustHours

	err := r.db.QueryRowContext(ctx, query, client_id, priority).Scan(
		&slaHours.Hours,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSlaNotFound
		}
		return nil, fmt.Errorf("error finding sla by params: %w", err)
	}

	return &slaHours, nil
}

func (r *slaRepository) Update(ctx context.Context, sla *domain.Sla) error {
	query := `UPDATE slas SET 
			client_id = $1, 
			priority = $2, 
			hours = $3 
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
		return ErrSlaNotFound
	}

	return nil
}

func (r *slaRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM sla WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting sla: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrSlaNotFound
	}

	return nil
}
