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

// TerminalRepository handles persistence operations for terminals.
type TerminalRepository struct {
	db *pgxpool.Pool
}

// NewTerminalRepository creates a new TerminalRepository.
func NewTerminalRepository(db *pgxpool.Pool) *TerminalRepository {
	return &TerminalRepository{db: db}
}

// Create inserts a new terminal record.
func (r *TerminalRepository) Create(ctx context.Context, t *domain.Terminal) error {
	query := `
		INSERT INTO terminals (id, owner_id, title, cwd, shell, status, color_tag, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(ctx, query,
		t.ID, t.OwnerID, t.Title, t.CWD, t.Shell, t.Status, t.ColorTag, t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("repository: failed to create terminal: %w", err)
	}
	return nil
}

// GetByID retrieves a terminal by its UUID.
func (r *TerminalRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Terminal, error) {
	query := `
		SELECT id, owner_id, title, cwd, shell, status, color_tag, created_at, updated_at
		FROM terminals WHERE id = $1`

	t := &domain.Terminal{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.OwnerID, &t.Title, &t.CWD, &t.Shell, &t.Status, &t.ColorTag, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repository: terminal %s not found: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("repository: failed to get terminal: %w", err)
	}
	return t, nil
}

// ListByOwner returns all terminals for a given owner, ordered by created_at DESC.
func (r *TerminalRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]*domain.Terminal, error) {
	query := `
		SELECT id, owner_id, title, cwd, shell, status, color_tag, created_at, updated_at
		FROM terminals WHERE owner_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list terminals by owner: %w", err)
	}
	defer rows.Close()

	return scanTerminals(rows)
}

// ListAll returns all terminals ordered by created_at DESC.
func (r *TerminalRepository) ListAll(ctx context.Context) ([]*domain.Terminal, error) {
	query := `
		SELECT id, owner_id, title, cwd, shell, status, color_tag, created_at, updated_at
		FROM terminals
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list all terminals: %w", err)
	}
	defer rows.Close()

	return scanTerminals(rows)
}

// UpdateStatus changes the status of a terminal.
func (r *TerminalRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TerminalStatus) error {
	query := `UPDATE terminals SET status = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("repository: failed to update terminal status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("repository: terminal %s not found: %w", id, ErrNotFound)
	}
	return nil
}

// Delete removes a terminal by its UUID.
func (r *TerminalRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM terminals WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository: failed to delete terminal: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("repository: terminal %s not found: %w", id, ErrNotFound)
	}
	return nil
}

// Count returns the total number of terminals.
func (r *TerminalRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM terminals`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to count terminals: %w", err)
	}
	return count, nil
}

// CountByStatus returns the number of terminals with a given status.
func (r *TerminalRepository) CountByStatus(ctx context.Context, status domain.TerminalStatus) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM terminals WHERE status = $1`, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to count terminals by status: %w", err)
	}
	return count, nil
}

// scanTerminals scans multiple terminal rows into a slice.
func scanTerminals(rows pgx.Rows) ([]*domain.Terminal, error) {
	var terminals []*domain.Terminal
	for rows.Next() {
		t := &domain.Terminal{}
		if err := rows.Scan(
			&t.ID, &t.OwnerID, &t.Title, &t.CWD, &t.Shell, &t.Status, &t.ColorTag, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("repository: failed to scan terminal row: %w", err)
		}
		terminals = append(terminals, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: error iterating terminal rows: %w", err)
	}
	return terminals, nil
}
