package jwt

import (
	"os"
	// "path/filepath" // Not strictly needed if paths are hardcoded relative to test
	"testing"
	"time"
    "log"
    "errors"


	"github.com/russian-steam/auth-service/internal/config"
	// "github.com/golang-jwt/jwt/v5" // Already imported by jwt package itself
)

// Note: The setupTestKeys helper from the prompt is removed as the subtask now directly handles key generation.
// The test will rely on the keys being present in "testdata/keys/" relative to this package.

func TestJWTService_GenerateAndValidateAccessToken(t *testing.T) {
	// Paths are relative to the package directory where the test is run.
	// The subtask execution previously created these files.
	testPrivKeyPath := "testdata/keys/jwtRS256.key"
	testPubKeyPath := "testdata/keys/jwtRS256.key.pub"

    // Verify keys exist before running test, to give a clear error if subtask failed setup
    if _, err := os.Stat(testPrivKeyPath); os.IsNotExist(err) {
        t.Fatalf("Test private key not found at %s. Ensure key generation step was successful.", testPrivKeyPath)
    }
    if _, err := os.Stat(testPubKeyPath); os.IsNotExist(err) {
        t.Fatalf("Test public key not found at %s. Ensure key generation step was successful.", testPubKeyPath)
    }


	cfg := &config.JWTConfig{
		AccessTokenExpiryMin:   1 * time.Minute,
		RefreshTokenExpiryDays: 1 * 24 * time.Hour, // Not used in this specific test function
		Issuer:                 "test-issuer",
		Audience:               "test-audience",
	}
    rsaCfg := &config.RSAKeysConfig{
        PrivateKeyPath: testPrivKeyPath,
        PublicKeyPath:  testPubKeyPath,
    }

	jwtSvc, err := NewJWTService(cfg, rsaCfg)
	if err != nil {
		t.Fatalf("NewJWTService() error = %v", err)
	}

	userID := "user123"
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user"}

	tokenString, jti, expTime, err := jwtSvc.GenerateAccessToken(userID, username, email, roles)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if tokenString == "" {
		t.Errorf("GenerateAccessToken() tokenString is empty")
	}
    if jti == "" {
        t.Errorf("GenerateAccessToken() jti is empty")
    }
    if expTime.IsZero() || expTime.Before(time.Now()) {
        t.Errorf("GenerateAccessToken() expTime is invalid: %v", expTime)
    }

	// Validate the token
	claims, err := jwtSvc.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("ValidateToken() UserID got = %v, want %v", claims.UserID, userID)
	}
	if claims.Username != username {
		t.Errorf("ValidateToken() Username got = %v, want %v", claims.Username, username)
	}
    if claims.ID != jti {
         t.Errorf("ValidateToken() JTI got = %v, want %v", claims.ID, jti)
    }
    if claims.Issuer != cfg.Issuer {
        t.Errorf("ValidateToken() Issuer got = %v, want %v", claims.Issuer, cfg.Issuer)
    }

    foundAudience := false
    for _, aud := range claims.Audience {
        if aud == cfg.Audience {
            foundAudience = true
            break
        }
    }
    if !foundAudience {
         t.Errorf("ValidateToken() Audience got = %v, want %v to be present", claims.Audience, cfg.Audience)
    }


	// Test expired token
	cfgExpired := &config.JWTConfig{
        AccessTokenExpiryMin:   -(1 * time.Minute), // Negative duration for expiry in the past
        Issuer:                 "test-issuer",
        Audience:               "test-audience",
    }

    // Need a new service instance if key loading is tied to it, or reconfigure.
    // NewJWTService reloads keys, so that's fine.
    jwtSvcExpired, err := NewJWTService(cfgExpired, rsaCfg)
    if err != nil {
        t.Fatalf("Failed to create JWTService for expired token test: %v", err)
    }
	expiredToken, _, _, _ := jwtSvcExpired.GenerateAccessToken(userID, username, email, roles)

    // Allow a small clock skew for validation, though JWT lib usually handles this.
    // Validate with the original service instance (jwtSvc) to ensure its ValidateToken works.
	time.Sleep(1 * time.Second) // Ensure token is definitely expired if clocks are perfectly synced.
	_, err = jwtSvc.ValidateToken(expiredToken)
	if err == nil || !errors.Is(err, ErrTokenExpired) {
		t.Errorf("ValidateToken() expected ErrTokenExpired for expired token, got %v", err)
	}

    // Test validation with a corrupted token string
    corruptedToken := tokenString + "corruption"
    _, err = jwtSvc.ValidateToken(corruptedToken)
    if err == nil {
        t.Errorf("ValidateToken() expected error for corrupted token, got nil")
    } else {
        log.Printf("Corrupted token validation error (expected): %v", err)
    }
}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
    testPrivKeyPath := "testdata/keys/jwtRS256.key"
	testPubKeyPath := "testdata/keys/jwtRS256.key.pub"

	cfg := &config.JWTConfig{
		RefreshTokenExpiryDays: 1 * 24 * time.Hour,
		Issuer:                 "test-issuer",
		Audience:               "test-audience",
	}
    rsaCfg := &config.RSAKeysConfig{
        PrivateKeyPath: testPrivKeyPath,
        PublicKeyPath:  testPubKeyPath,
    }
    jwtSvc, err := NewJWTService(cfg, rsaCfg)
    if err != nil {
        t.Fatalf("NewJWTService() error = %v", err)
    }

    userID := "user456"
    tokenString, jti, expTime, err := jwtSvc.GenerateRefreshToken(userID)
    if err != nil {
        t.Fatalf("GenerateRefreshToken() error = %v", err)
    }
    if tokenString == "" { t.Error("GenerateRefreshToken() tokenString is empty") }
    if jti == "" { t.Error("GenerateRefreshToken() jti is empty") }
    if expTime.IsZero() || expTime.Before(time.Now()) { t.Errorf("GenerateRefreshToken() expTime is invalid: %v", expTime) }

    claims, err := jwtSvc.ValidateToken(tokenString)
    if err != nil { t.Fatalf("ValidateToken() for refresh token error = %v", err) }
    if claims.UserID != userID { t.Errorf("ValidateToken() UserID got = %v, want %v", claims.UserID, userID) }
    if claims.ID != jti {t.Errorf("ValidateToken() JTI got = %v, want %v", claims.ID, jti) }
}

func TestJWTService_LoadKeys_Error(t *testing.T) {
    cfg := &config.JWTConfig{}
    // Test with non-existent key paths
    rsaCfgMissing := &config.RSAKeysConfig{
        PrivateKeyPath: "testdata/keys/nonexistent.key",
        PublicKeyPath:  "testdata/keys/nonexistent.pub",
    }
    _, err := NewJWTService(cfg, rsaCfgMissing)
    if err == nil {
        t.Errorf("NewJWTService() with missing keys expected error, got nil")
    } else {
        t.Logf("NewJWTService() with missing keys, got error: %v (expected)", err)
    }

    // Test with invalid key content (e.g., create empty files)
    _ = os.MkdirAll("testdata/keys_invalid", 0755)
    defer os.RemoveAll("testdata/keys_invalid")

    emptyPrivKey := "testdata/keys_invalid/empty.key"
    emptyPubKey := "testdata/keys_invalid/empty.pub"
    os.WriteFile(emptyPrivKey, []byte(""), 0644)
    os.WriteFile(emptyPubKey, []byte(""), 0644)

    rsaCfgInvalid := &config.RSAKeysConfig{ PrivateKeyPath: emptyPrivKey, PublicKeyPath: emptyPubKey }
    _, err = NewJWTService(cfg, rsaCfgInvalid)
    if err == nil {
        t.Errorf("NewJWTService() with invalid key content expected error, got nil")
    } else {
        t.Logf("NewJWTService() with invalid key content, got error: %v (expected)", err)
    }
}
