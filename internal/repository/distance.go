package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	_ "github.com/lib/pq"
)

var ErrDistanceNotFound = errors.New("distance not found")

type DistanceRepository interface {
	Create(ctx context.Context, distance *domain.Distance) (int, error)
	List(ctx context.Context) ([]domain.Distance, error)
	FindByID(ctx context.Context, id int) (*domain.Distance, error)
	FindByNumber(ctx context.Context, number string) (*domain.Distance, error)
	Update(ctx context.Context, distance *domain.Distance) error
	Delete(ctx context.Context, id int) error
}

type distanceRepository struct {
	db *sql.DB
}

func NewDistanceRepository(db *sql.DB) DistanceRepository {
	return &distanceRepository{
		db: db,
	}
}

func (r *distanceRepository) Create(ctx context.Context, distance *domain.Distance) (int, error) {
	query := `INSERT INTO distances (distance, ticket_number, provider_id, provider_name) 
			VALUES ($1, $2, $3, $4) RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		distance.Distance,
		distance.TicketNumber,
		distance.ProviderId,
		distance.ProviderName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating distance: %w", err)
	}

	return id, nil
}

func (r *distanceRepository) List(ctx context.Context) ([]domain.Distance, error) {
	query := `SELECT id, distance, ticket_number, provider_id, provider_name FROM distances`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing distances: %w", err)
	}
	defer rows.Close()

	var distances []domain.Distance
	for rows.Next() {
		var distance domain.Distance
		if err := rows.Scan(
			&distance.ID,
			&distance.Distance,
			&distance.TicketNumber,
			&distance.ProviderId,
			&distance.ProviderName); err != nil {
			return nil, fmt.Errorf("error scanning distance: %w", err)
		}
		distances = append(distances, distance)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating costs: %w", err)
	}

	return distances, nil
}

func (r *distanceRepository) FindByID(ctx context.Context, id int) (*domain.Distance, error) {
	query := `SELECT id, distance, ticket_number, provider_id, provider_name FROM distances WHERE id = $1`
	var distance domain.Distance

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&distance.ID,
		&distance.Distance,
		&distance.TicketNumber,
		&distance.ProviderId,
		&distance.ProviderName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDistanceNotFound
		}
		return nil, fmt.Errorf("error finding distance by id: %w", err)
	}

	return &distance, nil
}

func (r *distanceRepository) FindByNumber(ctx context.Context, number string) (*domain.Distance, error) {
	query := `SELECT id, distance, ticket_number, provider_id, provider_name FROM distances WHERE ticket_number = $1`
	var distance domain.Distance

	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&distance.ID,
		&distance.Distance,
		&distance.TicketNumber,
		&distance.ProviderId,
		&distance.ProviderName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDistanceNotFound
		}
		return nil, fmt.Errorf("error finding distance by ticket number: %w", err)
	}

	return &distance, nil
}

func (r *distanceRepository) Update(ctx context.Context, distance *domain.Distance) error {
	query := `UPDATE distances SET 
			distance = $1, 
			ticket_number = $2, 
			provider_id = $3, 
			provider_name = $4 
			WHERE id = $5`

	result, err := r.db.ExecContext(ctx, query,
		distance.Distance,
		distance.TicketNumber,
		distance.ProviderId,
		distance.ProviderName,
		distance.ID)
	if err != nil {
		return fmt.Errorf("error updating distance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrDistanceNotFound
	}

	return nil
}

func (r *distanceRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM distances WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting distance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrDistanceNotFound
	}

	return nil
}
