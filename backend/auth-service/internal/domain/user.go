package domain

import (
	"time"
)

// UserStatus defines the possible statuses of a user.
type UserStatus string

const (
	StatusPendingVerification UserStatus = "pending_verification"
	StatusActive              UserStatus = "active"
	StatusBlocked             UserStatus = "blocked"
	StatusDeleted             UserStatus = "deleted"
)

// User represents a user in the system.
// Corresponds to the 'users' table in backend/auth-service/docs/README.md (section 4.1)
type User struct {
	ID                string     `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"` // Never expose password hash
	Status            UserStatus `json:"status"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// UserRepository defines the interface for interacting with user storage.
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	// Delete(id string) error // Maybe not for MVP, use status instead
}
