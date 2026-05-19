// Package repository provides data access for user entities.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/termhub/backend/internal/domain"
)

// ErrNotFound is returned when a requested entity does not exist.
var ErrNotFound = errors.New("repository: entity not found")

// ErrDuplicate is returned when a unique constraint is violated.
var ErrDuplicate = errors.New("repository: duplicate entity")

// UserRepository handles persistence operations for users.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, username, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(ctx, query,
		user.ID, user.Username, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		// Check for unique constraint violation (pgx error code 23505)
		if isDuplicateKeyError(err) {
			return fmt.Errorf("repository: user with username %q already exists: %w", user.Username, ErrDuplicate)
		}
		return fmt.Errorf("repository: failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by their UUID.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, username, password_hash, role, created_at, updated_at
		FROM users WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repository: user with id %s not found: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("repository: failed to get user by id: %w", err)
	}
	return user, nil
}

// GetByUsername retrieves a user by their username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, password_hash, role, created_at, updated_at
		FROM users WHERE username = $1`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repository: user with username %q not found: %w", username, ErrNotFound)
		}
		return nil, fmt.Errorf("repository: failed to get user by username: %w", err)
	}
	return user, nil
}

// Update modifies an existing user record.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET username = $1, password_hash = $2, role = $3, updated_at = NOW()
		WHERE id = $4`

	result, err := r.db.Exec(ctx, query, user.Username, user.PasswordHash, user.Role, user.ID)
	if err != nil {
		return fmt.Errorf("repository: failed to update user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("repository: user with id %s not found: %w", user.ID, ErrNotFound)
	}
	return nil
}

// UpdateRole changes a user's role by their UUID.
func (r *UserRepository) UpdateRole(ctx context.Context, id uuid.UUID, role domain.Role) error {
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.Exec(ctx, query, role, id)
	if err != nil {
		return fmt.Errorf("repository: failed to update user role: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("repository: user with id %s not found: %w", id, ErrNotFound)
	}
	return nil
}

// Delete removes a user by their UUID.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository: failed to delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("repository: user with id %s not found: %w", id, ErrNotFound)
	}
	return nil
}

// List returns all users with pagination.
func (r *UserRepository) List(ctx context.Context, limit int, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, username, password_hash, role, created_at, updated_at
		FROM users ORDER BY created_at ASC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(
			&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("repository: failed to scan user row: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: error iterating user rows: %w", err)
	}
	return users, nil
}

// Count returns the total number of users.
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to count users: %w", err)
	}
	return count, nil
}

// isDuplicateKeyError checks if the error is a PostgreSQL unique constraint violation.
func isDuplicateKeyError(err error) bool {
	// pgx wraps the PostgreSQL error; check the error message for the unique_violation code
	return err != nil && (fmt.Sprintf("%v", err) != "" &&
		(contains(err.Error(), "23505") || contains(err.Error(), "duplicate key")))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
