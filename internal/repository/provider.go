package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	_ "github.com/lib/pq"
)

var ErrProviderNotFound = errors.New("provider not found")

type ProviderRepository interface {
	Create(ctx context.Context, provider *domain.Provider) (int, error)
	List(ctx context.Context) ([]domain.Provider, error)
	FindByID(ctx context.Context, id int) (*domain.Provider, error)
	FindByName(ctx context.Context, name string) (*domain.Provider, error)
	Update(ctx context.Context, provider *domain.Provider) error
	Delete(ctx context.Context, id int) error
}

type providerRepository struct {
	db *sql.DB
}

func NewProviderRepository(db *sql.DB) ProviderRepository {
	return &providerRepository{
		db: db,
	}
}

func (r *providerRepository) Create(ctx context.Context, provider *domain.Provider) (int, error) {
	query := `INSERT INTO providers (name, mobile, zipcode, state, city, neighborhood, address, complement) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		provider.Name,
		provider.Mobile,
		provider.Zipcode,
		provider.State,
		provider.City,
		provider.Neighborhood,
		provider.Address,
		provider.Complement).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating provider: %w", err)
	}

	return id, nil
}

func (r *providerRepository) List(ctx context.Context) ([]domain.Provider, error) {
	query := `SELECT id, name, mobile, zipcode, state, city, neighborhood, address, complement FROM providers ORDER BY name DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing providers: %w", err)
	}
	defer rows.Close()

	var providers []domain.Provider
	for rows.Next() {
		var provider domain.Provider
		if err := rows.Scan(
			&provider.ID,
			&provider.Name,
			&provider.Mobile,
			&provider.Zipcode,
			&provider.State,
			&provider.City,
			&provider.Neighborhood,
			&provider.Address,
			&provider.Complement); err != nil {
			return nil, fmt.Errorf("error scanning provider: %w", err)
		}
		providers = append(providers, provider)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating providers: %w", err)
	}

	return providers, nil
}

func (r *providerRepository) FindByID(ctx context.Context, id int) (*domain.Provider, error) {
	query := `SELECT id, name, mobile, zipcode, state, city, neighborhood, address, complement FROM providers WHERE id = $1`
	var provider domain.Provider

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProviderNotFound
		}
		return nil, fmt.Errorf("error finding provider by id: %w", err)
	}

	return &provider, nil
}

func (r *providerRepository) FindByName(ctx context.Context, name string) (*domain.Provider, error) {
	query := `SELECT * FROM providers WHERE name = $1 ORDER BY name LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, name)

	var provider domain.Provider
	if err := row.Scan(
		&provider.ID,
		&provider.Name,
		&provider.Mobile,
		&provider.Zipcode,
		&provider.State,
		&provider.City,
		&provider.Neighborhood,
		&provider.Address,
		&provider.Complement); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No provider found
		}
		return nil, fmt.Errorf("error scanning provider: %w", err)
	}

	return &provider, nil
}

func (r *providerRepository) Update(ctx context.Context, provider *domain.Provider) error {
	query := `UPDATE providers SET 
			name = $1, 
			mobile = $2, 
			zipcode = $3, 
			state = $4, 
			city = $5, 
			neighborhood = $6, 
			address = $7, 
			complement = $8
			WHERE id = $9`

	result, err := r.db.ExecContext(ctx, query,
		provider.Name,
		provider.Mobile,
		provider.Zipcode,
		provider.State,
		provider.City,
		provider.Neighborhood,
		provider.Address,
		provider.Complement,
		provider.ID)
	if err != nil {
		return fmt.Errorf("error updating provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrProviderNotFound
	}

	return nil
}

func (r *providerRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM providers WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrProviderNotFound
	}

	return nil
}
