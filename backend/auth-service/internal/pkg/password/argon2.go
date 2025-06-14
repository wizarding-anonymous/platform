package password

import (
	"crypto/rand" // Use crypto/rand for a cryptographically secure random number generator
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ArgonParams holds the parameters for Argon2id.
// These are example parameters and should be tuned based on security requirements and server performance.
type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultArgonParams provides sensible default parameters for Argon2id.
var DefaultArgonParams = ArgonParams{
	Memory:      64 * 1024, // 64MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// Argon2PasswordHasher implements the service.PasswordHasher interface.
type Argon2PasswordHasher struct {
	params ArgonParams
}

// NewArgon2PasswordHasher creates a new Argon2PasswordHasher.
// If params is nil, DefaultArgonParams will be used.
func NewArgon2PasswordHasher(params *ArgonParams) *Argon2PasswordHasher {
	p := DefaultArgonParams
	if params != nil {
		p = *params
	}
	return &Argon2PasswordHasher{params: p}
}

// HashPassword generates a salted Argon2id hash of the password.
// The format of the output string is: $argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
func (h *Argon2PasswordHasher) HashPassword(password string) (string, error) {
	params := h.params
	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil { // Corrected import: crypto/rand
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	// Encode salt and hash to base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format into standard string representation
	// Example: $argon2id$v=19$m=65536,t=3,p=2$YWFhYWFhYWFhYWFhYWFhYQ$ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZY
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, params.Memory, params.Iterations, params.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword compares a plain password with a stored Argon2id hash that was created by HashPassword.
func (h *Argon2PasswordHasher) VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return false, fmt.Errorf("unsupported hash type: %s", parts[1])
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return false, fmt.Errorf("unsupported argon2 version: %v", err)
	}

	params := &ArgonParams{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return false, fmt.Errorf("failed to parse argon2 params: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}
	params.SaltLength = uint32(len(salt))


	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}
    params.KeyLength = uint32(len(decodedHash))


	comparisonHash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

    if len(decodedHash) != len(comparisonHash) {
        return false, nil
    }

    var diff byte
    for i := 0; i < len(decodedHash); i++ {
        diff |= decodedHash[i] ^ comparisonHash[i]
    }
    return diff == 0, nil
}
