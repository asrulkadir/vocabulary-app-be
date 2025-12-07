package auth

import (
	"context"
	"database/sql"
)

// Repository handles data access for auth
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// FindByEmail finds a user by email
func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password, name, created_at, updated_at FROM users WHERE email = $1`

	var user User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Create creates a new user
func (r *repository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (email, password, name, created_at, updated_at) 
			  VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`

	return r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Name).Scan(&user.ID)
}

// FindByID finds a user by ID
func (r *repository) FindByID(ctx context.Context, id string) (*User, error) {
	query := `SELECT id, email, password, name, created_at, updated_at FROM users WHERE id = $1`

	var user User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
