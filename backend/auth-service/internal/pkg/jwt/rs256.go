// internal/pkg/jwt/rs256.go
package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid" // For JTI
    "github.com/russian-steam/auth-service/internal/config"
)

var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenMalformed   = errors.New("token is malformed")
	ErrUnexpectedMethod = errors.New("unexpected signing method")
    ErrNoPrivateKey     = errors.New("RSA private key not loaded")
    ErrNoPublicKey      = errors.New("RSA public key not loaded")
)

type Service struct {
    cfg         *config.JWTConfig
    rsaKeysCfg  *config.RSAKeysConfig
    privateKey  *rsa.PrivateKey
    publicKey   *rsa.PublicKey
}

func NewJWTService(cfg *config.JWTConfig, rsaKeysCfg *config.RSAKeysConfig) (*Service, error) {
    s := &Service{cfg: cfg, rsaKeysCfg: rsaKeysCfg}
    err := s.loadKeys()
    if err != nil {
        return nil, fmt.Errorf("failed to load RSA keys: %w", err)
    }
    return s, nil
}

func (s *Service) loadKeys() error {
    if s.rsaKeysCfg.PrivateKeyPath == "" || s.rsaKeysCfg.PublicKeyPath == "" {
        return fmt.Errorf("private or public key path not specified in config")
    }

	privKeyBytes, err := os.ReadFile(s.rsaKeysCfg.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("could not read private key: %w", err)
	}
	s.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		return fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyBytes, err := os.ReadFile(s.rsaKeysCfg.PublicKeyPath)
	if err != nil {
		return fmt.Errorf("could not read public key: %w", err)
	}
	s.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		return fmt.Errorf("could not parse public key: %w", err)
	}
    if s.privateKey == nil {
        return ErrNoPrivateKey
    }
    if s.publicKey == nil {
        return ErrNoPublicKey
    }
	return nil
}

type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

func (s *Service) GenerateAccessToken(userID, username, email string, roles []string) (string, string, time.Time, error) {
    if s.privateKey == nil {
		return "", "", time.Time{}, ErrNoPrivateKey
	}
	expirationTime := time.Now().Add(s.cfg.AccessTokenExpiryMin)
	jti := uuid.NewString()

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.cfg.Issuer,
			Subject:   userID,
			ID:        jti,
			Audience:  []string{s.cfg.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to sign access token: %w", err)
	}
	return tokenString, jti, expirationTime, nil
}

func (s *Service) GenerateRefreshToken(userID string) (string, string, time.Time, error) {
    if s.privateKey == nil {
		return "", "", time.Time{}, ErrNoPrivateKey
	}
	expirationTime := time.Now().Add(s.cfg.RefreshTokenExpiryDays)
	jti := uuid.NewString() // Unique ID for the refresh token

	claims := &Claims{ // Refresh token might have fewer claims
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.cfg.Issuer,
			Subject:   userID,
			ID:        jti, // JTI for refresh token itself
			Audience:  []string{s.cfg.Audience}, // Or a specific audience for refresh tokens
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return tokenString, jti, expirationTime, nil
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
    if s.publicKey == nil {
		return nil, ErrNoPublicKey
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedMethod, token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}
