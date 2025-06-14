package password

import (
	"strings"
	"testing"
)

func TestArgon2PasswordHasher_HashAndVerify(t *testing.T) {
	hasher := NewArgon2PasswordHasher(nil) // Use default params
	password := "strongP@sswOrd123!"

	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if !strings.HasPrefix(hashedPassword, "$argon2id$") {
		t.Errorf("Hashed password does not have argon2id prefix: %s", hashedPassword)
	}

	// Test valid password
	match, err := hasher.VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Fatalf("VerifyPassword() error for correct password = %v", err)
	}
	if !match {
		t.Errorf("VerifyPassword() got = %v, want %v for correct password", match, true)
	}

	// Test invalid password
	match, err = hasher.VerifyPassword("wrongPassword", hashedPassword)
	if err != nil {
		t.Fatalf("VerifyPassword() error for incorrect password = %v", err)
	}
	if match {
		t.Errorf("VerifyPassword() got = %v, want %v for incorrect password", match, false)
	}
}

func TestArgon2PasswordHasher_VerifyInvalidFormat(t *testing.T) {
	hasher := NewArgon2PasswordHasher(nil)
	_, err := hasher.VerifyPassword("anypassword", "invalidhashformat")
	if err == nil {
		t.Errorf("VerifyPassword() expected error for invalid hash format, got nil")
	}
}
