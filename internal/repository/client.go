package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	_ "github.com/lib/pq"
)

type ClientRepository interface {
	Create(ctx context.Context, client *domain.Client) (int, error)
	List(ctx context.Context) ([]domain.Client, error)
	FindByID(ctx context.Context, id int) (*domain.Client, error)
	Update(ctx context.Context, client *domain.Client) error
	Delete(ctx context.Context, id int) error
}

type clientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) ClientRepository {
	return &clientRepository{
		db: db,
	}
}

func (r *clientRepository) Create(ctx context.Context, client *domain.Client) (int, error) {
	query := `INSERT INTO clients (name) VALUES ($1) RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query, client.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating client: %w", err)
	}

	return id, nil
}

func (r *clientRepository) List(ctx context.Context) ([]domain.Client, error) {
	query := `SELECT id, name FROM clients ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing clients: %w", err)
	}
	defer rows.Close()

	var clients []domain.Client
	for rows.Next() {
		var client domain.Client
		if err := rows.Scan(&client.ID, &client.Name); err != nil {
			return nil, fmt.Errorf("error scanning client: %w", err)
		}
		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clients: %w", err)
	}

	return clients, nil
}

func (r *clientRepository) FindByID(ctx context.Context, id int) (*domain.Client, error) {
	query := `SELECT id, name FROM clients WHERE id = $1`
	var client domain.Client

	err := r.db.QueryRowContext(ctx, query, id).Scan(&client.ID, &client.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error finding client by id: %w", err)
	}

	return &client, nil
}

func (r *clientRepository) Update(ctx context.Context, client *domain.Client) error {
	query := `UPDATE clients SET name = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, client.Name, client.ID)
	if err != nil {
		return fmt.Errorf("error updating client: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *clientRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM clients WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting client: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
