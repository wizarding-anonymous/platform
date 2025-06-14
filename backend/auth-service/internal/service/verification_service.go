// internal/service/verification_service.go
package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"
    "log"
    "errors"


	"github.com/russian-steam/auth-service/internal/domain"
)

const verificationCodeLength = 6
const verificationCodeExpiry = 15 * time.Minute

var (
    ErrVerificationCodeNotFound = errors.New("verification code not found or already used")
    ErrVerificationCodeExpired  = errors.New("verification code has expired")
    ErrVerificationCodeInvalid  = errors.New("invalid verification code")
)


type VerificationService struct {
	vcRepo   domain.VerificationCodeRepository
    userRepo domain.UserRepository
}

func NewVerificationService(vcRepo domain.VerificationCodeRepository, userRepo domain.UserRepository) *VerificationService {
	return &VerificationService{vcRepo: vcRepo, userRepo: userRepo}
}

func (s *VerificationService) generateSecureRandomString(length int) (string, error) {
	// For a 6-digit code, we can use numbers, but hex is simpler for now
    // For truly numeric codes, a different approach would be needed.
	numBytes := length / 2 // Each byte becomes 2 hex characters
    if length%2 != 0 {
        numBytes++
    }
	bytes := make([]byte, numBytes)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}
	hexStr := hex.EncodeToString(bytes)
	return hexStr[:length], nil // Trim to desired length
}

func (s *VerificationService) hashVerificationCode(code string) string {
	hasher := sha256.New()
	hasher.Write([]byte(code))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *VerificationService) GenerateEmailVerificationCode(userID, email string) (*domain.VerificationCode, string, error) {
    // Delete any existing email verification codes for this user to prevent multiple valid codes
    if err := s.vcRepo.DeleteByUserIDAndType(userID, domain.VerificationTypeEmail); err != nil {
        log.Printf("Warning: Failed to delete old email verification codes for user %s: %v", userID, err)
        // Continue, as this is not necessarily a fatal error for generating a new code
    }

    rawCode, err := s.generateSecureRandomString(verificationCodeLength)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate raw code: %w", err)
	}

	hashedCode := s.hashVerificationCode(rawCode)
	expiresAt := time.Now().UTC().Add(verificationCodeExpiry)

	vc := &domain.VerificationCode{
        // ID will be set by the repository upon creation if it's auto-generated
		UserID:    userID,
		Type:      domain.VerificationTypeEmail,
		CodeHash:  hashedCode,
		Target:    email, // Store the email this code is for
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.vcRepo.Create(vc); err != nil {
		return nil, "", fmt.Errorf("failed to create verification code: %w", err)
	}
	return vc, rawCode, nil
}

// VerifyEmailCode checks if the provided raw code is valid for the given email.
// If valid, it marks the user's email as verified and deletes the code.
func (s *VerificationService) VerifyEmailCode(rawCode string, targetEmail string) (*domain.User, error) {
	hashedCode := s.hashVerificationCode(rawCode)
	vc, err := s.vcRepo.FindByCodeHash(hashedCode, domain.VerificationTypeEmail)
	if err != nil {
        log.Printf("Verification code not found for hash derived from '%s' (type: %s): %v", rawCode, domain.VerificationTypeEmail, err)
		return nil, ErrVerificationCodeNotFound
	}

	if time.Now().UTC().After(vc.ExpiresAt) {
        log.Printf("Verification code ID %s for email %s has expired.", vc.ID, targetEmail)
        if delErr := s.vcRepo.Delete(vc.ID); delErr != nil { // Delete expired code
            log.Printf("Failed to delete expired verification code %s: %v", vc.ID, delErr)
        }
		return nil, ErrVerificationCodeExpired
	}

    if vc.Target != targetEmail {
        log.Printf("Verification code target mismatch. Expected %s, got %s for code ID %s", targetEmail, vc.Target, vc.ID)
        // Do NOT delete the code here, as it might be a typo by the user for a different email.
        return nil, ErrVerificationCodeInvalid
    }

    // Code is valid, not expired, and target matches.
    user, err := s.markUserEmailAsVerified(vc.UserID, vc.Target)
    if err != nil {
        // Log the error but don't necessarily delete the code if user update failed,
        // as it might be a transient DB issue. However, for security, if the user update fails,
        // it might be better to require a new code. This depends on the desired behavior.
        log.Printf("Failed to mark email as verified for user %s after code %s validation: %v", vc.UserID, vc.ID, err)
        return nil, fmt.Errorf("failed to mark email as verified: %w", err)
    }

    // Successfully verified and user updated, now delete the code.
    if err := s.vcRepo.Delete(vc.ID); err != nil {
        log.Printf("Warning: Failed to delete verification code %s after successful use: %v", vc.ID, err)
        // Not returning an error here as the core operation (verification) succeeded.
    }

	return user, nil
}

// markUserEmailAsVerified updates the user's status to active and sets email_verified_at.
func (s *VerificationService) markUserEmailAsVerified(userID string, email string) (*domain.User, error) {
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return nil, fmt.Errorf("could not find user %s: %w", userID, err)
    }

    // Double check that the email being verified actually belongs to this user.
    // This is important if the `targetEmail` in VerifyEmailCode could somehow be manipulated.
    if user.Email != email {
        log.Printf("Security check: User %s email (%s) does not match verification target email (%s).", userID, user.Email, email)
        return nil, fmt.Errorf("verification target email '%s' does not match user's primary email '%s'", email, user.Email)
    }

    changed := false
    if user.Status == domain.StatusPendingVerification {
        user.Status = domain.StatusActive
        changed = true
        log.Printf("User %s status updated to Active.", userID)
    }

    if user.EmailVerifiedAt == nil {
        now := time.Now().UTC()
        user.EmailVerifiedAt = &now
        changed = true
        log.Printf("User %s email %s marked as verified at %v.", userID, email, now)
    }

    if changed {
        user.UpdatedAt = time.Now().UTC()
        if err := s.userRepo.Update(user); err != nil {
            return nil, fmt.Errorf("failed to update user %s after email verification: %w", userID, err)
        }
    } else {
         log.Printf("User %s email %s was already verified and active.", userID, email)
    }
    return user, nil
}
