// internal/service/token_service.go
package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
    "log"
    "errors"


	"github.com/russian-steam/auth-service/internal/domain"
	"github.com/russian-steam/auth-service/internal/pkg/jwt"
)

var (
    ErrRefreshTokenNotFound = errors.New("refresh token not found")
    ErrRefreshTokenRevoked  = errors.New("refresh token has been revoked")
    ErrRefreshTokenExpired  = errors.New("refresh token has expired")
)

type JWTTokenGeneratorValidator interface {
    GenerateAccessToken(userID, username, email string, roles []string) (string, string, time.Time, error)
    GenerateRefreshToken(userID string) (string, string, time.Time, error)
    ValidateToken(tokenString string) (*jwt.Claims, error)
}

type TokenService struct {
	refreshTokenRepo domain.RefreshTokenRepository
	jtiBlacklistRepo domain.JtiBlacklistRepository
	jwtService       JWTTokenGeneratorValidator
}

func NewTokenService(
	refreshTokenRepo domain.RefreshTokenRepository,
	jtiBlacklistRepo domain.JtiBlacklistRepository,
	jwtService JWTTokenGeneratorValidator,
) *TokenService {
	return &TokenService{
		refreshTokenRepo: refreshTokenRepo,
		jtiBlacklistRepo: jtiBlacklistRepo,
		jwtService:       jwtService,
	}
}

func (s *TokenService) hashToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *TokenService) GenerateTokens(user *domain.User) (accessToken, rawRefreshToken, accessTokenJTI string, accessExp, refreshExp time.Time, err error) {
	roles := []string{"user"}
    if user.Status != domain.StatusActive && user.Status != domain.StatusPendingVerification {
        return "", "", "", time.Time{}, time.Time{}, fmt.Errorf("user account is not active or pending verification, status: %s", user.Status)
    }


	accessToken, accessTokenJTI, accessExp, err = s.jwtService.GenerateAccessToken(user.ID, user.Username, user.Email, roles)
	if err != nil {
		return "", "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	rawRefreshToken, refreshTokenJTI, refreshExp, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshToken := &domain.RefreshToken{
		ID:        refreshTokenJTI,
		UserID:    user.ID,
		TokenHash: s.hashToken(rawRefreshToken),
		ExpiresAt: refreshExp,
		CreatedAt: time.Now().UTC(),
        IsRevoked: false,
	}
	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		log.Printf("Warning: failed to store refresh token for user %s: %v", user.ID, err)
		// Depending on policy, might want to return an error here
	}

	return accessToken, rawRefreshToken, accessTokenJTI, accessExp, refreshExp, nil
}

func (s *TokenService) ValidateAndRefreshTokens(rawRefreshTokenValue string, user *domain.User) (newAccessToken, newRawRefreshToken, newAccessTokenJTI string, newAccessExp, newRefreshExp time.Time, err error) {
    claims, err := s.jwtService.ValidateToken(rawRefreshTokenValue)
    if err != nil {
        log.Printf("Incoming refresh token failed basic validation: %v", err)
        if errors.Is(err, jwt.ErrTokenExpired) {
            // Attempt to find and delete the token from DB if it's just expired
            storedRT, findErr := s.refreshTokenRepo.GetByTokenHash(s.hashToken(rawRefreshTokenValue))
            if findErr == nil && storedRT != nil {
                if delErr := s.refreshTokenRepo.Delete(storedRT.ID); delErr != nil {
                    log.Printf("Failed to delete expired refresh token %s from DB: %v", storedRT.ID, delErr)
                }
            }
            return "", "", "", time.Time{}, time.Time{}, ErrRefreshTokenExpired
        }
        return "", "", "", time.Time{}, time.Time{}, fmt.Errorf("invalid refresh token structure: %w", err)
    }

    storedRefreshTokenID := claims.ID // JTI of the refresh token

    // Fetch the refresh token from the database using its JTI
    storedRT, err := s.refreshTokenRepo.GetByID(storedRefreshTokenID)
    if err != nil {
        if errors.Is(err, ErrRefreshTokenNotFound) { // Assuming GetByID returns a specific error for not found
             log.Printf("Refresh token with JTI %s not found in DB", storedRefreshTokenID)
            return "", "", "", time.Time{}, time.Time{}, ErrRefreshTokenNotFound
        }
        log.Printf("Error fetching stored refresh token by JTI %s: %v", storedRefreshTokenID, err)
        return "", "", "", time.Time{}, time.Time{}, fmt.Errorf("could not retrieve refresh token: %w", err)
    }

    // Verify the token hash to ensure the provided raw token matches the stored one
    if storedRT.TokenHash != s.hashToken(rawRefreshTokenValue) {
        log.Printf("Refresh token hash mismatch for JTI %s. Potential token theft or misuse.", storedRefreshTokenID)
        // Optionally, revoke all tokens for the user here as a security measure.
        // s.RevokeAllUserRefreshTokens(storedRT.UserID)
        return "", "", "", time.Time{}, time.Time{}, ErrRefreshTokenNotFound // Treat as not found or invalid
    }

    if storedRT.IsRevoked {
        log.Printf("Attempt to use revoked refresh token with JTI %s", storedRefreshTokenID)
        return "", "", "", time.Time{}, time.Time{}, ErrRefreshTokenRevoked
    }
    if time.Now().UTC().After(storedRT.ExpiresAt) {
        log.Printf("Attempt to use expired refresh token with JTI %s (expired at %v)", storedRefreshTokenID, storedRT.ExpiresAt)
        if err := s.refreshTokenRepo.Delete(storedRT.ID); err != nil {
            log.Printf("Failed to delete expired refresh token %s from DB: %v", storedRT.ID, err)
        }
        return "", "", "", time.Time{}, time.Time{}, ErrRefreshTokenExpired
    }

    // If we reached here, the refresh token is valid, not revoked, and not expired.
    // It's good practice to revoke the old refresh token (one-time use) and issue new ones.
    if err := s.RevokeRefreshToken(storedRT.ID); err != nil {
         log.Printf("Failed to revoke old refresh token %s: %v. Proceeding with new token generation.", storedRT.ID, err)
         // Decide if this is a critical error. For now, we proceed.
    }

    // Generate new pair of tokens
    return s.GenerateTokens(user)
}


func (s *TokenService) RevokeRefreshToken(tokenID string) error {
    // rt, err := s.refreshTokenRepo.GetByID(tokenID) // This GetByID call is not strictly needed if we just delete.
    // if err != nil {
    // 	if errors.Is(err, ErrRefreshTokenNotFound) {
    // 		return nil
    // 	}
    // 	return fmt.Errorf("failed to get refresh token by ID before deleting: %w", err)
    // }
	return s.refreshTokenRepo.Delete(tokenID)
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a given user.
// This is useful for "logout all sessions" functionality.
func (s *TokenService) RevokeAllUserRefreshTokens(userID string) error {
    err := s.refreshTokenRepo.DeleteByUserID(userID)
    if err != nil {
        return fmt.Errorf("failed to delete refresh tokens for user %s: %w", userID, err)
    }
    return nil
}

func (s *TokenService) BlacklistAccessToken(jti string, expiresAt time.Time) error {
	return s.jtiBlacklistRepo.AddToBlacklist(jti, expiresAt)
}

func (s *TokenService) IsAccessTokenBlacklisted(jti string) (bool, error) {
    return s.jtiBlacklistRepo.IsBlacklisted(jti)
}

// JwtService returns the underlying JWT token generator/validator.
// This is used by handlers that need direct access to token validation logic,
// e.g. gRPC ValidateToken or HTTP middleware.
func (s *TokenService) JwtService() JWTTokenGeneratorValidator {
	return s.jwtService
}
