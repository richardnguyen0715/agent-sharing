// Package domain defines the core business entities and value types.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user's access level in the system.
type Role string

const (
	RoleAdmin      Role = "admin"
	RoleOperator   Role = "operator"
	RoleViewer     Role = "viewer"
	RoleRestricted Role = "restricted"
)

// ValidRoles contains all valid role values.
var ValidRoles = []Role{RoleAdmin, RoleOperator, RoleViewer, RoleRestricted}

// IsValid checks whether the role is a recognized value.
func (r Role) IsValid() bool {
	for _, v := range ValidRoles {
		if r == v {
			return true
		}
	}
	return false
}

// User represents a registered user of the system.
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         Role      `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
