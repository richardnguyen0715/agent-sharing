package domain

import (
	"time"

	"github.com/google/uuid"
)

// TerminalStatus represents the current state of a terminal.
type TerminalStatus string

const (
	TerminalStatusActive     TerminalStatus = "active"
	TerminalStatusIdle       TerminalStatus = "idle"
	TerminalStatusTerminated TerminalStatus = "terminated"
)

// Terminal represents a managed terminal instance.
type Terminal struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	OwnerID   uuid.UUID      `json:"owner_id" db:"owner_id"`
	Title     string         `json:"title" db:"title"`
	CWD       string         `json:"cwd" db:"cwd"`
	Shell     string         `json:"shell" db:"shell"`
	Status    TerminalStatus `json:"status" db:"status"`
	ColorTag  string         `json:"color_tag" db:"color_tag"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// TerminalSession represents a WebSocket connection to a terminal.
type TerminalSession struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	TerminalID     uuid.UUID  `json:"terminal_id" db:"terminal_id"`
	WebSocketID    string     `json:"websocket_id" db:"websocket_id"`
	ConnectedAt    time.Time  `json:"connected_at" db:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at,omitempty" db:"disconnected_at"`
}
