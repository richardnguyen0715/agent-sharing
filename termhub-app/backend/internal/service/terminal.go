// Package service contains business logic for terminal management.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/termhub/backend/internal/domain"
	ptymanager "github.com/termhub/backend/internal/pty"
	"github.com/termhub/backend/internal/repository"
)

// Terminal service sentinel errors.
var (
	ErrTerminalNotFound = errors.New("terminal: not found")
	ErrNotOwner         = errors.New("terminal: access denied, not owner")
)

// TerminalService handles terminal lifecycle operations.
type TerminalService struct {
	repo       *repository.TerminalRepository
	ptyManager *ptymanager.Manager
	logger     zerolog.Logger
}

// NewTerminalService creates a new TerminalService.
func NewTerminalService(repo *repository.TerminalRepository, ptyManager *ptymanager.Manager) *TerminalService {
	return &TerminalService{
		repo:       repo,
		ptyManager: ptyManager,
		logger:     log.With().Str("component", "terminal-service").Logger(),
	}
}

// Create provisions a new terminal for the given user.
func (s *TerminalService) Create(ctx context.Context, ownerID uuid.UUID, title, cwd, shell, colorTag string) (*domain.Terminal, error) {
	// Apply defaults.
	if title == "" {
		title = "Terminal"
	}
	if cwd == "" {
		cwd = "/tmp"
	}
	if shell == "" {
		shell = "/bin/bash"
	}
	if colorTag == "" {
		colorTag = "#00FF00"
	}

	id := uuid.New()
	now := time.Now().UTC()

	terminal := &domain.Terminal{
		ID:        id,
		OwnerID:   ownerID,
		Title:     title,
		CWD:       cwd,
		Shell:     shell,
		Status:    domain.TerminalStatusActive,
		ColorTag:  colorTag,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save to database first.
	if err := s.repo.Create(ctx, terminal); err != nil {
		return nil, fmt.Errorf("terminal: failed to create record: %w", err)
	}

	// Spawn PTY session.
	_, err := s.ptyManager.Create(id, ownerID, shell, cwd)
	if err != nil {
		// Clean up the database record if PTY creation fails.
		_ = s.repo.Delete(ctx, id)
		return nil, fmt.Errorf("terminal: failed to spawn PTY: %w", err)
	}

	s.logger.Info().
		Str("terminal_id", id.String()).
		Str("owner_id", ownerID.String()).
		Str("title", title).
		Str("shell", shell).
		Msg("terminal created")

	return terminal, nil
}

// Get retrieves a terminal by ID and verifies ownership.
func (s *TerminalService) Get(ctx context.Context, id, ownerID uuid.UUID) (*domain.Terminal, error) {
	terminal, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTerminalNotFound
		}
		return nil, fmt.Errorf("terminal: failed to get: %w", err)
	}

	if terminal.OwnerID != ownerID {
		return nil, ErrNotOwner
	}

	return terminal, nil
}

// List returns all terminals owned by a user.
func (s *TerminalService) List(ctx context.Context, ownerID uuid.UUID) ([]*domain.Terminal, error) {
	terminals, err := s.repo.ListByOwner(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("terminal: failed to list: %w", err)
	}
	if terminals == nil {
		terminals = []*domain.Terminal{}
	}
	return terminals, nil
}

// ListAll returns all terminals across all users (admin use).
func (s *TerminalService) ListAll(ctx context.Context) ([]*domain.Terminal, error) {
	terminals, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("terminal: failed to list all: %w", err)
	}
	if terminals == nil {
		terminals = []*domain.Terminal{}
	}
	return terminals, nil
}

// Delete terminates and removes a terminal, verifying ownership.
func (s *TerminalService) Delete(ctx context.Context, id, ownerID uuid.UUID) error {
	terminal, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTerminalNotFound
		}
		return fmt.Errorf("terminal: failed to get for delete: %w", err)
	}

	if terminal.OwnerID != ownerID {
		return ErrNotOwner
	}

	// Kill the PTY session.
	if err := s.ptyManager.Kill(id); err != nil {
		s.logger.Error().Err(err).Str("terminal_id", id.String()).Msg("failed to kill PTY during delete")
	}

	// Delete from database.
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("terminal: failed to delete record: %w", err)
	}

	s.logger.Info().
		Str("terminal_id", id.String()).
		Str("owner_id", ownerID.String()).
		Msg("terminal deleted")

	return nil
}

// Resize sends a resize signal to the terminal's PTY.
func (s *TerminalService) Resize(_ context.Context, id uuid.UUID, cols, rows uint16) error {
	return s.ptyManager.Resize(id, cols, rows)
}

// GetRepo returns the terminal repository for direct access when needed.
func (s *TerminalService) GetRepo() *repository.TerminalRepository {
	return s.repo
}

// GetPTYManager returns the PTY manager.
func (s *TerminalService) GetPTYManager() *ptymanager.Manager {
	return s.ptyManager
}
