// internal/service/mocks_test.go
package service

import (
	// "context" // Not used by current mocks, but might be for future ones
	"fmt"
	"strings"
	"sync"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5" // Added for direct access to golang-jwt types
	"github.com/russian-steam/auth-service/internal/domain"
	"github.com/russian-steam/auth-service/internal/pkg/jwt" // For our local jwt.Claims struct
	// "github.com/russian-steam/auth-service/internal/pkg/password" // Not directly used by mocks if methods are simple
)

// --- Mock PasswordHasher ---
type MockPasswordHasher struct {
	HashPasswordFunc   func(password string) (string, error)
	VerifyPasswordFunc func(password, hashedPassword string) (bool, error)
}

func (m *MockPasswordHasher) HashPassword(password string) (string, error) {
	if m.HashPasswordFunc != nil {
		return m.HashPasswordFunc(password)
	}
	return "hashed_" + password, nil
}
func (m *MockPasswordHasher) VerifyPassword(password, hashedPassword string) (bool, error) {
	if m.VerifyPasswordFunc != nil {
		return m.VerifyPasswordFunc(password, hashedPassword)
	}
	return "hashed_"+password == hashedPassword, nil
}

// --- Mock UserRepository ---
type MockUserRepository struct {
	mu             sync.Mutex
	users          map[string]*domain.User
	CreateErr      error
	GetIDErr       error
	GetEmailErr    error
	GetUsernameErr error
	UpdateErr      error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: make(map[string]*domain.User)}
}
func (m *MockUserRepository) Create(user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.CreateErr != nil {
		return m.CreateErr
	}
	if _, exists := m.users[user.ID]; exists {
		return fmt.Errorf("mock: user with ID %s already exists", user.ID)
	}
	for _, u := range m.users {
		if u.Email == user.Email || u.Username == user.Username {
			return ErrUserAlreadyExists // service level error for testing service logic
		}
	}
	m.users[user.ID] = user
	return nil
}
func (m *MockUserRepository) GetByID(id string) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.GetIDErr != nil {
		return nil, m.GetIDErr
	}
	user, exists := m.users[id]
	if !exists || user.Status == domain.StatusDeleted {
		return nil, ErrUserNotFound // service level error
	}
	// Return a copy to prevent test side effects on the stored mock object
    userCopy := *user
	return &userCopy, nil
}
func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.GetEmailErr != nil {
		return nil, m.GetEmailErr
	}
	for _, user := range m.users {
		if user.Email == email { // Allow finding user even if status is 'deleted'
            userCopy := *user
			return &userCopy, nil
		}
	}
	return nil, ErrUserNotFound // service level error
}
func (m *MockUserRepository) GetByUsername(username string) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.GetUsernameErr != nil {
		return nil, m.GetUsernameErr
	}
	for _, user := range m.users {
		if user.Username == username { // Allow finding user even if status is 'deleted'
            userCopy := *user
			return &userCopy, nil
		}
	}
	return nil, ErrUserNotFound // service level error
}
func (m *MockUserRepository) Update(user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	if _, exists := m.users[user.ID]; !exists {
		return ErrUserNotFound
	}
	m.users[user.ID] = user
	return nil
}

// --- Mock RefreshTokenRepository ---
type MockRefreshTokenRepository struct {
	mu                sync.Mutex
	tokens            map[string]*domain.RefreshToken // Key: token ID (JTI)
	CreateErr         error
	GetByIDErr        error
	GetByTokenHashErr error
	SetRevokedErr     error
	DeleteErr         error
    DeleteByUserIDErr error
}

func NewMockRefreshTokenRepository() *MockRefreshTokenRepository {
	return &MockRefreshTokenRepository{tokens: make(map[string]*domain.RefreshToken)}
}
func (m *MockRefreshTokenRepository) Create(token *domain.RefreshToken) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.CreateErr != nil {
		return m.CreateErr
	}
	m.tokens[token.ID] = token
	return nil
}
func (m *MockRefreshTokenRepository) GetByID(id string) (*domain.RefreshToken, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	token, ok := m.tokens[id]
	if !ok {
		return nil, ErrRefreshTokenNotFound
	}
    tokenCopy := *token
	return &tokenCopy, nil
}
func (m *MockRefreshTokenRepository) GetByTokenHash(tokenHash string) (*domain.RefreshToken, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.GetByTokenHashErr != nil {
		return nil, m.GetByTokenHashErr
	}
	for _, t := range m.tokens {
		if t.TokenHash == tokenHash {
            tokenCopy := *t
			return &tokenCopy, nil
		}
	}
	return nil, ErrRefreshTokenNotFound
}
func (m *MockRefreshTokenRepository) SetRevoked(id string, isRevoked bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
    if m.SetRevokedErr != nil { return m.SetRevokedErr }
	token, ok := m.tokens[id]
	if !ok {
		return ErrRefreshTokenNotFound
	}
	token.IsRevoked = isRevoked
	// If SetRevoked implies deletion in the real repo (as TokenService expects for one-time use)
    // then mimic that here for consistent testing of TokenService logic.
    // Based on TokenService.RevokeRefreshToken, it calls repo.Delete.
    // So, this SetRevoked in mock might not be directly called by TokenService if it prefers Delete.
    // Let's assume this is for direct SetRevoked calls if any, or for testing the repo interface itself.
	return nil
}
func (m *MockRefreshTokenRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	if _, ok := m.tokens[id]; !ok {
		return ErrRefreshTokenNotFound // Or return nil if not finding is okay for delete
	}
	delete(m.tokens, id)
	return nil
}
func (m *MockRefreshTokenRepository) DeleteByUserID(userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
    if m.DeleteByUserIDErr != nil { return m.DeleteByUserIDErr }
	for id, token := range m.tokens {
		if token.UserID == userID {
			delete(m.tokens, id)
		}
	}
	return nil
}

// --- Mock JtiBlacklistRepository ---
type MockJtiBlacklistRepository struct {
	mu        sync.Mutex
	blacklist map[string]time.Time
    AddErr    error
    IsErr     error
}

func NewMockJtiBlacklistRepository() *MockJtiBlacklistRepository {
	return &MockJtiBlacklistRepository{blacklist: make(map[string]time.Time)}
}
func (m *MockJtiBlacklistRepository) AddToBlacklist(jti string, expiresAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
    if m.AddErr != nil { return m.AddErr }
	m.blacklist[jti] = expiresAt
	return nil
}
func (m *MockJtiBlacklistRepository) IsBlacklisted(jti string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
    if m.IsErr != nil { return false, m.IsErr }
	exp, ok := m.blacklist[jti]
	if !ok {
		return false, nil
	}
	return time.Now().UTC().Before(exp), nil
}

// --- Mock JWTTokenGeneratorValidator (JWT Service) ---
type MockJWTService struct {
	GenerateAccessTokenFunc  func(userID, username, email string, roles []string) (string, string, time.Time, error)
	GenerateRefreshTokenFunc func(userID string) (string, string, time.Time, error)
	ValidateTokenFunc        func(tokenString string) (*jwt.Claims, error)
}

func (m *MockJWTService) GenerateAccessToken(userID, username, email string, roles []string) (string, string, time.Time, error) {
	if m.GenerateAccessTokenFunc != nil {
		return m.GenerateAccessTokenFunc(userID, username, email, roles)
	}
	return "mock_access_token_" + userID, "jti_access_" + userID, time.Now().Add(15 * time.Minute), nil
}
func (m *MockJWTService) GenerateRefreshToken(userID string) (string, string, time.Time, error) {
	if m.GenerateRefreshTokenFunc != nil {
		return m.GenerateRefreshTokenFunc(userID)
	}
	return "mock_refresh_token_" + userID, "jti_refresh_" + userID, time.Now().Add(7 * 24 * time.Hour), nil
}
    func (m *MockJWTService) ValidateToken(tokenString string) (*jwt.Claims, error) { // Returns local jwt.Claims
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(tokenString)
	}
	if strings.HasPrefix(tokenString, "mock_access_token_") || strings.HasPrefix(tokenString, "mock_refresh_token_") {
		parts := strings.Split(tokenString, "_") // Simple split, assumes userID doesn't have "_"
        if len(parts) < 3 { return nil, fmt.Errorf("invalid mock token format for parts extraction") }
		userID := parts[len(parts)-1]

        jtiPrefix := "jti_access_"
        if strings.HasPrefix(tokenString, "mock_refresh_token_") {
            jtiPrefix = "jti_refresh_"
        }
        jti := jtiPrefix + userID

		exp := time.Now().Add(15 * time.Minute)
            isExpired := strings.Contains(tokenString, "expired_token")
		if isExpired {
			exp = time.Now().Add(-15 * time.Minute)
		}

            // Construct the RegisteredClaims part using gojwt types
            registeredClaims := gojwt.RegisteredClaims{
                ID:        jti,
                ExpiresAt: gojwt.NewNumericDate(exp),
                Issuer:    "mock-issuer",
                Audience:  gojwt.ClaimStrings{"mock-audience"},
            }

            // Our local jwt.Claims embeds jwt.RegisteredClaims
            // So we need to ensure the mock returns this structure.
            // The fields UserID, Username, Email are specific to our local jwt.Claims.
            localClaims := &jwt.Claims{
                UserID:   userID,
                Username: "mockuser",
                Email:    "mock@example.com",
                RegisteredClaims: registeredClaims,
            }

            if isExpired {
                return localClaims, jwt.ErrTokenExpired // Return our local error type
            }
		return localClaims, nil
	}
        if strings.Contains(tokenString, "expired_token") {
             return nil, jwt.ErrTokenExpired // Return our local error type
    }
	return nil, fmt.Errorf("invalid mock token: %s", tokenString)
}

// --- Mock VerificationCodeRepository ---
type MockVerificationCodeRepository struct {
	mu                sync.Mutex
	codes             map[string]*domain.VerificationCode // Key: codeHash
	CreateErr         error
	FindByCodeHashErr error
	DeleteErr         error
    DeleteByUserIDErr error
}

func NewMockVerificationCodeRepository() *MockVerificationCodeRepository {
	return &MockVerificationCodeRepository{codes: make(map[string]*domain.VerificationCode)}
}
func (m *MockVerificationCodeRepository) Create(vc *domain.VerificationCode) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.CreateErr != nil {
		return m.CreateErr
	}
	m.codes[vc.CodeHash] = vc
	return nil
}
func (m *MockVerificationCodeRepository) FindByCodeHash(codeHash string, vcType domain.VerificationType) (*domain.VerificationCode, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.FindByCodeHashErr != nil {
		return nil, m.FindByCodeHashErr
	}
	vc, ok := m.codes[codeHash]
	if !ok || vc.Type != vcType {
		return nil, ErrVerificationCodeNotFound // service level error
	}
    vcCopy := *vc
	return &vcCopy, nil
}
func (m *MockVerificationCodeRepository) Delete(id string) error { // ID here is vc.ID (uuid), not hash
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	found := false
	for hash, vc := range m.codes {
		if vc.ID == id {
			delete(m.codes, hash)
            found = true
			break
		}
	}
    if !found {
        // return ErrVerificationCodeNotFound // Or nil if not finding is okay for delete
    }
	return nil
}
func (m *MockVerificationCodeRepository) DeleteByUserIDAndType(userID string, vcType domain.VerificationType) error {
	m.mu.Lock()
	defer m.mu.Unlock()
    if m.DeleteByUserIDErr != nil { return m.DeleteByUserIDErr }
	for hash, vc := range m.codes {
		if vc.UserID == userID && vc.Type == vcType {
			delete(m.codes, hash)
		}
	}
	return nil
}
