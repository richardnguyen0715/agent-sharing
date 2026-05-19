// Package pty manages pseudo-terminal instances for terminal emulation.
package pty

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Session represents a running PTY session.
type Session struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	cmd       *exec.Cmd
	ptmx      *os.File
	CreatedAt time.Time
	mu        sync.Mutex
	closed    bool
}

// Manager manages the lifecycle of PTY sessions.
type Manager struct {
	mu       sync.RWMutex
	sessions map[uuid.UUID]*Session
	logger   zerolog.Logger
}

// NewManager creates a new PTY Manager.
func NewManager() *Manager {
	return &Manager{
		sessions: make(map[uuid.UUID]*Session),
		logger:   log.With().Str("component", "pty-manager").Logger(),
	}
}

// Create starts a new PTY session with the given shell command.
func (m *Manager) Create(id, ownerID uuid.UUID, shell string, cwd string) (*Session, error) {
	if shell == "" {
		shell = "/bin/bash"
	}
	if cwd == "" {
		cwd = "/tmp"
	}

	cmd := exec.Command(shell)
	cmd.Dir = cwd
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	// Set process group so we can signal the entire group on kill.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	ptmx, err := pty.Start(cmd)
	if err != nil {
		m.logger.Error().Err(err).
			Str("session_id", id.String()).
			Str("shell", shell).
			Str("cwd", cwd).
			Msg("failed to start PTY session")
		return nil, fmt.Errorf("pty: failed to start shell %q: %w", shell, err)
	}

	// Set a default window size.
	_ = pty.Setsize(ptmx, &pty.Winsize{Rows: 24, Cols: 80})

	session := &Session{
		ID:        id,
		OwnerID:   ownerID,
		cmd:       cmd,
		ptmx:      ptmx,
		CreatedAt: time.Now().UTC(),
	}

	m.mu.Lock()
	m.sessions[id] = session
	m.mu.Unlock()

	m.logger.Info().
		Str("session_id", id.String()).
		Str("owner_id", ownerID.String()).
		Str("shell", shell).
		Str("cwd", cwd).
		Int("pid", cmd.Process.Pid).
		Msg("PTY session created")

	return session, nil
}

// Get returns a PTY session by ID.
func (m *Manager) Get(id uuid.UUID) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	return s, ok
}

// Resize changes the window size for a PTY session.
func (m *Manager) Resize(id uuid.UUID, cols, rows uint16) error {
	m.mu.RLock()
	session, ok := m.sessions[id]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("pty: session %s not found", id)
	}
	if err := session.Resize(rows, cols); err != nil {
		return fmt.Errorf("pty: failed to resize session %s: %w", id, err)
	}
	m.logger.Debug().
		Str("session_id", id.String()).
		Uint16("cols", cols).
		Uint16("rows", rows).
		Msg("PTY session resized")
	return nil
}

// Kill terminates a PTY session with graceful shutdown.
// Sends SIGTERM first, waits up to 3 seconds, then SIGKILL.
func (m *Manager) Kill(id uuid.UUID) error {
	m.mu.Lock()
	session, ok := m.sessions[id]
	if !ok {
		m.mu.Unlock()
		return nil
	}
	delete(m.sessions, id)
	m.mu.Unlock()

	m.logger.Info().Str("session_id", id.String()).Msg("killing PTY session")
	return session.kill()
}

// ListByOwner returns all session IDs for a given owner.
func (m *Manager) ListByOwner(ownerID uuid.UUID) []uuid.UUID {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var ids []uuid.UUID
	for _, s := range m.sessions {
		if s.OwnerID == ownerID {
			ids = append(ids, s.ID)
		}
	}
	return ids
}

// ActiveCount returns the number of active PTY sessions.
func (m *Manager) ActiveCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.sessions)
}

// Cleanup kills all active PTY sessions. Called during graceful shutdown.
func (m *Manager) Cleanup() {
	m.mu.Lock()
	sessions := make([]*Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		sessions = append(sessions, s)
	}
	m.sessions = make(map[uuid.UUID]*Session)
	m.mu.Unlock()

	m.logger.Info().Int("count", len(sessions)).Msg("cleaning up all PTY sessions")

	for _, s := range sessions {
		if err := s.kill(); err != nil {
			m.logger.Error().Err(err).Str("session_id", s.ID.String()).Msg("error killing session during cleanup")
		}
	}

	m.logger.Info().Msg("PTY cleanup complete")
}

// Reader returns an io.Reader for the PTY output.
func (s *Session) Reader() io.Reader {
	return s.ptmx
}

// Writer returns an io.Writer for the PTY input.
func (s *Session) Writer() io.Writer {
	return s.ptmx
}

// Resize changes the PTY window size.
func (s *Session) Resize(rows, cols uint16) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("pty: session is closed")
	}
	return pty.Setsize(s.ptmx, &pty.Winsize{Rows: rows, Cols: cols})
}

// IsClosed reports whether this session has been closed.
func (s *Session) IsClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

// kill terminates the session with SIGTERM, then SIGKILL after timeout.
func (s *Session) kill() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true

	if s.cmd.Process != nil {
		// Try SIGTERM to process group first.
		pgid, err := syscall.Getpgid(s.cmd.Process.Pid)
		if err == nil {
			_ = syscall.Kill(-pgid, syscall.SIGTERM)
		} else {
			_ = s.cmd.Process.Signal(syscall.SIGTERM)
		}

		// Wait with timeout for graceful exit.
		done := make(chan error, 1)
		go func() {
			done <- s.cmd.Wait()
		}()

		select {
		case <-done:
			// Process exited gracefully.
		case <-time.After(3 * time.Second):
			// Force kill.
			if pgid, err := syscall.Getpgid(s.cmd.Process.Pid); err == nil {
				_ = syscall.Kill(-pgid, syscall.SIGKILL)
			} else {
				_ = s.cmd.Process.Kill()
			}
			<-done
		}
	}

	return s.ptmx.Close()
}
