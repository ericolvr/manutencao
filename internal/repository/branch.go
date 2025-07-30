package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ericolvr/maintenance-v2/internal/domain"
	_ "github.com/lib/pq"
)

var ErrBranchNotFound = errors.New("branch not found")

type BranchRepository interface {
	Create(ctx context.Context, branch *domain.Branch) (int, error)
	List(ctx context.Context) ([]domain.Branch, error)
	FindByID(ctx context.Context, id int) (*domain.Branch, error)
	FindByUniorg(ctx context.Context, uniorg string) (*domain.Branch, error)
	GetByClient(ctx context.Context, client string) ([]domain.Branch, error)
	Update(ctx context.Context, branch *domain.Branch) error
	Delete(ctx context.Context, id int) error
}

type branchRepository struct {
	db *sql.DB
}

func NewBranchRepository(db *sql.DB) BranchRepository {
	return &branchRepository{
		db: db,
	}
}

func (r *branchRepository) Create(ctx context.Context, branch *domain.Branch) (int, error) {
	query := `INSERT INTO branchs (client, name, uniorg, zipcode, state, city, neighborhood, address, complement) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		branch.Client,
		branch.Name,
		branch.Uniorg,
		branch.Zipcode,
		branch.State,
		branch.City,
		branch.Neighborhood,
		branch.Address,
		branch.Complement).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating branch: %w", err)
	}

	return id, nil
}

func (r *branchRepository) List(ctx context.Context) ([]domain.Branch, error) {
	query := `SELECT id, client, name, uniorg, zipcode, state, city, neighborhood, address, complement FROM branchs ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing branchs: %w", err)
	}
	defer rows.Close()

	var branchs []domain.Branch
	for rows.Next() {
		var branch domain.Branch
		if err := rows.Scan(
			&branch.ID,
			&branch.Client,
			&branch.Name,
			&branch.Uniorg,
			&branch.Zipcode,
			&branch.State,
			&branch.City,
			&branch.Neighborhood,
			&branch.Address,
			&branch.Complement); err != nil {
			return nil, fmt.Errorf("error scanning branch: %w", err)
		}
		branchs = append(branchs, branch)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating branchs: %w", err)
	}

	return branchs, nil
}

func (r *branchRepository) FindByID(ctx context.Context, id int) (*domain.Branch, error) {
	query := `SELECT id, name, client, uniorg, zipcode, state, city, neighborhood, address, complement FROM branchs WHERE id = $1`
	var branch domain.Branch

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBranchNotFound
		}
		return nil, fmt.Errorf("error finding branch by id: %w", err)
	}

	return &branch, nil
}

func (r *branchRepository) FindByUniorg(ctx context.Context, uniorg string) (*domain.Branch, error) {
	query := `SELECT * FROM branchs WHERE uniorg = $1 LIMIT 1`

	var branch domain.Branch
	err := r.db.QueryRowContext(ctx, query, uniorg).Scan(
		&branch.ID,
		&branch.Client,
		&branch.Name,
		&branch.Uniorg,
		&branch.Zipcode,
		&branch.State,
		&branch.City,
		&branch.Neighborhood,
		&branch.Address,
		&branch.Complement)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error finding branch by uniorg: %w", err)
	}

	return &branch, nil
}

func (r *branchRepository) GetByClient(ctx context.Context, client string) ([]domain.Branch, error) {
	query := `SELECT id, name, uniorg, zipcode, state, city, neighborhood, address FROM branchs WHERE client = $1 ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query, client)
	if err != nil {
		return nil, fmt.Errorf("error getting branches by client: %w", err)
	}
	defer rows.Close()

	var branches []domain.Branch
	for rows.Next() {
		var branch domain.Branch
		if err := rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.Uniorg,
			&branch.Zipcode,
			&branch.State,
			&branch.City,
			&branch.Neighborhood,
			&branch.Address); err != nil {
			return nil, fmt.Errorf("error scanning branch: %w", err)
		}
		branches = append(branches, branch)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating branches: %w", err)
	}

	return branches, nil
}

func (r *branchRepository) Update(ctx context.Context, branch *domain.Branch) error {
	query := `UPDATE branchs SET 
			name = $1, 
			client = $2, 
			uniorg = $3, 
			zipcode = $4, 
			state = $5, 
			city = $6, 
			neighborhood = $7, 
			address = $8, 
			complement = $9
			WHERE id = $10`

	result, err := r.db.ExecContext(ctx, query,
		branch.Name,
		branch.Client,
		branch.Uniorg,
		branch.Zipcode,
		branch.State,
		branch.City,
		branch.Neighborhood,
		branch.Address,
		branch.Complement,

		branch.ID)
	if err != nil {
		return fmt.Errorf("error updating branch: %w", err)
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

func (r *branchRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM branchs WHERE id = $1`

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
