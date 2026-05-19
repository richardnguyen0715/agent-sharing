package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditAction represents the type of auditable event.
type AuditAction string

const (
	AuditActionLogin              AuditAction = "login"
	AuditActionLogout             AuditAction = "logout"
	AuditActionTerminalCreate     AuditAction = "terminal.create"
	AuditActionTerminalDelete     AuditAction = "terminal.delete"
	AuditActionTerminalConnect    AuditAction = "terminal.connect"
	AuditActionTerminalDisconnect AuditAction = "terminal.disconnect"
	AuditActionUserCreate         AuditAction = "user.create"
	AuditActionUserUpdate         AuditAction = "user.update"
	AuditActionUserLogin          AuditAction = "user.login"
	AuditActionUserRegister       AuditAction = "user.register"
	AuditActionPolicyViolation    AuditAction = "policy.violation"
)

// AuditLog records a security-relevant event in the system.
type AuditLog struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	UserID     *uuid.UUID      `json:"user_id,omitempty" db:"user_id"`
	TerminalID *uuid.UUID      `json:"terminal_id,omitempty" db:"terminal_id"`
	Action     AuditAction     `json:"action" db:"action"`
	Metadata   json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}
