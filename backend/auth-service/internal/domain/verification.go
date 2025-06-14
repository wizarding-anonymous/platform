package domain

import (
	"time"
)

// VerificationType defines the type of verification code.
type VerificationType string

const (
	VerificationTypeEmail VerificationType = "email_verification"
	VerificationTypePasswordReset VerificationType = "password_reset"
)

// VerificationCode represents a code used for verification purposes (e.g., email, password reset).
// Corresponds to the 'VerificationCode' table in backend/auth-service/docs/README.md (section 4.1)
type VerificationCode struct {
	ID        string           `json:"id"`
	UserID    string           `json:"user_id"`
	Type      VerificationType `json:"type"`
	CodeHash  string           `json:"-"` // Store a hash of the code
	Target    string           `json:"target"` // e.g., email address for email verification
	ExpiresAt time.Time        `json:"expires_at"`
	CreatedAt time.Time        `json:"created_at"`
}

// VerificationCodeRepository defines the interface for interacting with verification code storage.
type VerificationCodeRepository interface {
	Create(vc *VerificationCode) error
	FindByCodeHash(codeHash string, vcType VerificationType) (*VerificationCode, error)
	Delete(id string) error
	DeleteByUserIDAndType(userID string, vcType VerificationType) error // Clean up old codes
}
