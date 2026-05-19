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
	OwnerID   uuid.UUID      `json:"ownerId" db:"owner_id"`
	Title     string         `json:"title" db:"title"`
	CWD       string         `json:"cwd" db:"cwd"`
	Shell     string         `json:"shell" db:"shell"`
	Status    TerminalStatus `json:"status" db:"status"`
	ColorTag  string         `json:"colorTag" db:"color_tag"`
	CreatedAt time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time      `json:"updatedAt" db:"updated_at"`
}

// TerminalSession represents a WebSocket connection to a terminal.
type TerminalSession struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	TerminalID     uuid.UUID  `json:"terminalId" db:"terminal_id"`
	WebSocketID    string     `json:"websocketId" db:"websocket_id"`
	ConnectedAt    time.Time  `json:"connectedAt" db:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnectedAt,omitempty" db:"disconnected_at"`
}
