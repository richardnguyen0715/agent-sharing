package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/termhub/backend/internal/domain"
)

// AuditService handles creation and retrieval of audit log entries.
type AuditService struct {
	db     *pgxpool.Pool
	logger zerolog.Logger
}

// NewAuditService creates a new AuditService.
func NewAuditService(db *pgxpool.Pool) *AuditService {
	return &AuditService{
		db:     db,
		logger: log.With().Str("component", "audit-service").Logger(),
	}
}

// Log records an audit event.
func (s *AuditService) Log(ctx context.Context, userID *uuid.UUID, terminalID *uuid.UUID, action domain.AuditAction, metadata map[string]interface{}) error {
	id := uuid.New()
	now := time.Now().UTC()

	var metaJSON json.RawMessage
	if metadata != nil {
		data, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("audit: failed to marshal metadata: %w", err)
		}
		metaJSON = data
	} else {
		metaJSON = json.RawMessage("{}")
	}

	query := `
		INSERT INTO audit_logs (id, user_id, terminal_id, action, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(ctx, query, id, userID, terminalID, action, metaJSON, now)
	if err != nil {
		s.logger.Error().Err(err).
			Str("action", string(action)).
			Msg("failed to write audit log")
		return fmt.Errorf("audit: failed to insert log: %w", err)
	}

	s.logger.Debug().
		Str("id", id.String()).
		Str("action", string(action)).
		Msg("audit log recorded")

	return nil
}

// List retrieves audit log entries with limit and offset for pagination, ordered by created_at DESC.
func (s *AuditService) List(ctx context.Context, limit int, offset int) ([]*domain.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, user_id, terminal_id, action, metadata, created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("audit: failed to list logs: %w", err)
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		entry := &domain.AuditLog{}
		if err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.TerminalID, &entry.Action, &entry.Metadata, &entry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("audit: failed to scan log row: %w", err)
		}
		logs = append(logs, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("audit: error iterating log rows: %w", err)
	}
	return logs, nil
}

// ListByUser retrieves audit log entries for a specific user.
func (s *AuditService) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `
		SELECT id, user_id, terminal_id, action, metadata, created_at
		FROM audit_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := s.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("audit: failed to list logs by user: %w", err)
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		entry := &domain.AuditLog{}
		if err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.TerminalID, &entry.Action, &entry.Metadata, &entry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("audit: failed to scan log row: %w", err)
		}
		logs = append(logs, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("audit: error iterating log rows: %w", err)
	}
	return logs, nil
}
