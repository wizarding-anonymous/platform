// internal/service/user_service.go
package service

import (
	"errors"
	"fmt"
	"time"
    "log"

	"github.com/google/uuid"
	"github.com/russian-steam/auth-service/internal/domain"
	// "github.com/russian-steam/auth-service/internal/pkg/password" // Not directly used, PasswordHasher interface is
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user with this email or username already exists")
	ErrInvalidCredentials   = errors.New("invalid email or password")
    ErrPasswordVerification = errors.New("password verification failed")
)

type PasswordHasher interface {
	HashPassword(password string) (string, error) // No params here, configured in impl
	VerifyPassword(password, hashedPassword string) (bool, error)
}

type UserService struct {
	userRepo domain.UserRepository
	hasher   PasswordHasher
}

func NewUserService(userRepo domain.UserRepository, hasher PasswordHasher) *UserService {
	return &UserService{userRepo: userRepo, hasher: hasher}
}

func (s *UserService) RegisterUser(username, email, plainPassword string) (*domain.User, error) {
	// Check if user already exists
	existingUserByEmail, err := s.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, ErrUserNotFound) { // Assuming GetByEmail returns ErrUserNotFound from repo
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if existingUserByEmail != nil {
		return nil, ErrUserAlreadyExists
	}

    existingUserByUsername, err := s.userRepo.GetByUsername(username)
    if err != nil && !errors.Is(err, ErrUserNotFound) { // Assuming GetByUsername returns ErrUserNotFound from repo
        return nil, fmt.Errorf("failed to check username existence: %w", err)
    }
    if existingUserByUsername != nil {
        return nil, ErrUserAlreadyExists
    }

	hashedPassword, err := s.hasher.HashPassword(plainPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Status:       domain.StatusPendingVerification,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUserByID(userID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) { // Assuming GetByID returns ErrUserNotFound from repo
			return nil, ErrUserNotFound // Return the service-level error
		}
		return nil, fmt.Errorf("failed to get user by ID from repository: %w", err)
	}
	return user, nil
}

func (s *UserService) AuthenticateUser(email, plainPassword string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) { // Assuming GetByEmail returns ErrUserNotFound from repo
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user.Status == domain.StatusBlocked || user.Status == domain.StatusDeleted {
		return nil, fmt.Errorf("user account is %s", user.Status)
	}

    if user.Status == domain.StatusPendingVerification {
        log.Printf("User %s is pending verification.", user.Email)
    }


	match, err := s.hasher.VerifyPassword(plainPassword, user.PasswordHash)
	if err != nil {
		log.Printf("Error verifying password for user %s: %v", email, err)
		return nil, ErrPasswordVerification
	}
	if !match {
		return nil, ErrInvalidCredentials
	}

    now := time.Now().UTC()
    user.LastLoginAt = &now
    user.UpdatedAt = now
    if err := s.userRepo.Update(user); err != nil {
        log.Printf("Failed to update last login time for user %s: %v", user.ID, err)
    }

	return user, nil
}
