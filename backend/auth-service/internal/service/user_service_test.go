package service

import (
	"strings"
	"testing"
	"time"
    "errors"

	"github.com/russian-steam/auth-service/internal/domain"
	"github.com/google/uuid"
)

func TestUserService_RegisterUser_Success(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockHasher := &MockPasswordHasher{}
	userService := NewUserService(mockUserRepo, mockHasher)

	username := "testuser"
	email := "test@example.com"
	password := "password123"

	user, err := userService.RegisterUser(username, email, password)
	if err != nil {
		t.Fatalf("RegisterUser() error = %v, wantErr %v", err, false)
	}
	if user == nil {
		t.Fatalf("RegisterUser() user is nil")
	}
	if user.Username != username || user.Email != email {
		t.Errorf("RegisterUser() user data mismatch")
	}
	if user.Status != domain.StatusPendingVerification {
		t.Errorf("RegisterUser() user status got = %s, want %s", user.Status, domain.StatusPendingVerification)
	}
	if !strings.HasPrefix(user.PasswordHash, "hashed_") {
		t.Errorf("RegisterUser() password not hashed correctly")
	}
}

func TestUserService_RegisterUser_AlreadyExistsByEmail(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	existingEmail := "existing@example.com"
	mockUserRepo.Create(&domain.User{ID: "id1", Username: "existinguser", Email: existingEmail, Status: domain.StatusActive})

	mockHasher := &MockPasswordHasher{}
	userService := NewUserService(mockUserRepo, mockHasher)

	_, err := userService.RegisterUser("newuser", existingEmail, "password123")
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Errorf("RegisterUser() with existing email, error = %v, want %v", err, ErrUserAlreadyExists)
	}
}

func TestUserService_RegisterUser_AlreadyExistsByUsername(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	existingUsername := "existing_username"
	mockUserRepo.Create(&domain.User{ID: "id2", Username: existingUsername, Email: "unique@example.com", Status: domain.StatusActive})

	mockHasher := &MockPasswordHasher{}
	userService := NewUserService(mockUserRepo, mockHasher)

	_, err := userService.RegisterUser(existingUsername, "newemail@example.com", "password123")
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Errorf("RegisterUser() with existing username, error = %v, want %v", err, ErrUserAlreadyExists)
	}
}


func TestUserService_AuthenticateUser_Success(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockHasher := &MockPasswordHasher{}

	userID := uuid.NewString()
	username := "authuser"
	email := "auth@example.com"
	plainPassword := "password123"
	hashedPassword, _ := mockHasher.HashPassword(plainPassword)

	mockUserRepo.Create(&domain.User{
		ID:           userID,
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Status:       domain.StatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	userService := NewUserService(mockUserRepo, mockHasher)


	user, err := userService.AuthenticateUser(email, plainPassword)
	if err != nil {
		t.Fatalf("AuthenticateUser() error = %v", err)
	}
	if user == nil || user.Email != email {
		t.Errorf("AuthenticateUser() failed to authenticate correct user")
	}
    if user.LastLoginAt == nil || user.LastLoginAt.IsZero() {
        t.Errorf("AuthenticateUser() LastLoginAt not updated")
    }
}

func TestUserService_AuthenticateUser_InvalidPassword(t *testing.T) {
	mockUserRepo := NewMockUserRepository()
	mockHasher := &MockPasswordHasher{}
    hashedPassword, _ := mockHasher.HashPassword("password123")
	mockUserRepo.Create(&domain.User{ID: "id1", Username: "user", Email: "user@example.com", PasswordHash: hashedPassword, Status: domain.StatusActive})
	userService := NewUserService(mockUserRepo, mockHasher)

	_, err := userService.AuthenticateUser("user@example.com", "wrongpassword")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("AuthenticateUser() error = %v, want %v", err, ErrInvalidCredentials)
	}
}

func TestUserService_AuthenticateUser_UserNotFound(t *testing.T) {
    mockUserRepo := NewMockUserRepository()
    mockHasher := &MockPasswordHasher{}
    userService := NewUserService(mockUserRepo, mockHasher)

    _, err := userService.AuthenticateUser("nonexistent@example.com", "password")
    if !errors.Is(err, ErrInvalidCredentials) { // Service returns InvalidCredentials for user not found
        t.Errorf("AuthenticateUser() for non-existent user, error = %v, want %v", err, ErrInvalidCredentials)
    }
}

func TestUserService_AuthenticateUser_UserNotActiveOrPending(t *testing.T) {
    // mockUserRepo := NewMockUserRepository() // This was unused
    mockHasher := &MockPasswordHasher{}
    hashedPassword, _ := mockHasher.HashPassword("password123")

    testCases := []struct {
        name          string
        status        domain.UserStatus
        expectAuthErr bool // True if AuthenticateUser itself should error for this status
        errContains   string // Substring to check in error message if expectAuthErr is true
    }{
        {name: "Blocked", status: domain.StatusBlocked, expectAuthErr: true, errContains: "user account is blocked"},
        {name: "Deleted", status: domain.StatusDeleted, expectAuthErr: true, errContains: "user account is deleted"},
        // StatusPendingVerification should pass at service level, but handler will block actual login
        {name: "PendingVerification", status: domain.StatusPendingVerification, expectAuthErr: false, errContains: ""},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            email := strings.ToLower(tc.name) + "@example.com"
            username := strings.ToLower(tc.name) + "_user"
            currentRepo := NewMockUserRepository()
            currentRepo.Create(&domain.User{
                ID:           uuid.NewString(),
                Username:     username,
                Email:        email,
                PasswordHash: hashedPassword,
                Status:       tc.status,
            })
            userService := NewUserService(currentRepo, mockHasher)
            user, err := userService.AuthenticateUser(email, "password123")

            if tc.expectAuthErr {
                if err == nil {
                    t.Errorf("AuthenticateUser() with status %s expected error, got nil", tc.status)
                } else {
                    if tc.errContains != "" && !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tc.errContains)) {
                        t.Errorf("AuthenticateUser() with status %s, error message '%s' did not contain expected substring '%s'", tc.status, err.Error(), tc.errContains)
                    }
                    t.Logf("AuthenticateUser() with status %s, got error: %v (expected error containing '%s')", tc.status, err, tc.errContains)
                }
            } else {
                if err != nil {
                    t.Errorf("AuthenticateUser() with status %s expected no error, got %v", tc.status, err)
                }
                if user == nil {
                     t.Errorf("AuthenticateUser() with status %s expected user, got nil", tc.status)
                }
            }
        })
    }
}

func TestUserService_GetUserByID_Success(t *testing.T) {
    mockUserRepo := NewMockUserRepository()
    userService := NewUserService(mockUserRepo, nil) // Hasher not needed for GetUserByID

    expectedUser := &domain.User{ID: "user1", Username: "test", Email: "test@example.com", Status: domain.StatusActive}
    mockUserRepo.Create(expectedUser)

    user, err := userService.GetUserByID("user1")
    if err != nil {
        t.Fatalf("GetUserByID() error = %v", err)
    }
    if user == nil || user.ID != "user1" {
        t.Errorf("GetUserByID() did not return the correct user")
    }
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
    mockUserRepo := NewMockUserRepository()
    userService := NewUserService(mockUserRepo, nil)

    _, err := userService.GetUserByID("nonexistent")
    if !errors.Is(err, ErrUserNotFound) {
        t.Errorf("GetUserByID() error = %v, want %v", err, ErrUserNotFound)
    }
}
