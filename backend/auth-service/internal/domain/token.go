package domain

import (
	"time"
)

// RefreshToken represents a refresh token in the system.
// Corresponds to the 'RefreshToken' table in backend/auth-service/docs/README.md (section 4.1)
type RefreshToken struct {
	ID          string    `json:"id"` // Can be JTI (JWT ID)
	UserID      string    `json:"user_id"`
	TokenHash   string    `json:"-"` // Store a hash of the refresh token, not the token itself if it's long-lived and opaque
	SessionID   *string   `json:"session_id,omitempty"` // Optional: link to a session
	IsRevoked   bool      `json:"is_revoked"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// RefreshTokenRepository defines the interface for interacting with refresh token storage.
type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	GetByTokenHash(tokenHash string) (*RefreshToken, error)
	GetByID(id string) (*RefreshToken, error) // If ID is JTI
	SetRevoked(id string, isRevoked bool) error
	Delete(id string) error
	DeleteByUserID(userID string) error // For "logout all sessions"
}

// JtiBlacklistRepository defines the interface for managing a JTI blacklist (for access tokens).
// This will likely be implemented with Redis.
type JtiBlacklistRepository interface {
	AddToBlacklist(jti string, expiresAt time.Time) error
	IsBlacklisted(jti string) (bool, error)
}
