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

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (int, error)
	List(ctx context.Context) ([]domain.User, error)
	FindByID(ctx context.Context, id int) (*domain.User, error)
	FindByName(ctx context.Context, name string) ([]domain.User, error)
	FindByMobile(ctx context.Context, mobile string) ([]domain.User, error)
	FindUsersToTicket(ctx context.Context) ([]dto.UsersToTicket, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (int, error) {
	query := `INSERT INTO users (name, mobile, password, role, status) 
            VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int
	var err error
	err = r.db.QueryRowContext(ctx, query,
		user.Name,
		user.Mobile,
		user.Password,
		user.Role,
		user.Status).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}

	return id, nil
}

func (r *userRepository) List(ctx context.Context) ([]domain.User, error) {
	query := `SELECT id, name, mobile, password, role, status FROM users ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Mobile,
			&user.Password,
			&user.Role,
			&user.Status); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, name, mobile, password, role, status FROM users WHERE id = $1`
	var user domain.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Mobile,
		&user.Password,
		&user.Role,
		&user.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error finding user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByName(ctx context.Context, name string) ([]domain.User, error) {
	query := `SELECT id, name, mobile, password, role, status FROM users WHERE name = $1`

	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("error finding users by name: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Mobile,
			&user.Password,
			&user.Role,
			&user.Status); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *userRepository) FindByMobile(ctx context.Context, mobile string) ([]domain.User, error) {
	query := `SELECT id, name, mobile, password, role, status FROM users WHERE mobile = $1`

	rows, err := r.db.QueryContext(ctx, query, mobile)
	if err != nil {
		return nil, fmt.Errorf("error finding users by mobile: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Mobile,
			&user.Password,
			&user.Role,
			&user.Status); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *userRepository) FindUsersToTicket(ctx context.Context) ([]dto.UsersToTicket, error) {
	query := `SELECT id, name FROM users WHERE (role = 1 OR role = 2) AND status = $1 ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query, 1)
	if err != nil {
		return nil, fmt.Errorf("error finding users by role: %w", err)
	}
	defer rows.Close()

	var users []dto.UsersToTicket
	for rows.Next() {
		var user dto.UsersToTicket
		if err := rows.Scan(
			&user.ID,
			&user.Name); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET 
			name = $1, 
			mobile = $2, 
			password = $3, 
			role = $4, 
			status = $5 
		WHERE id = $6`

	result, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Mobile,
		user.Password,
		user.Role,
		user.Status,
		user.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
