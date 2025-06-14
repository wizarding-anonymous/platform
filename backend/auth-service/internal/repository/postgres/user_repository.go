package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/russian-steam/auth-service/internal/domain"
    "github.com/russian-steam/auth-service/internal/service" // For ErrUserNotFound
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (id, username, email, password_hash, status, created_at, updated_at, email_verified_at, last_login_at)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(context.Background(), query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.Status,
		user.CreatedAt, user.UpdatedAt, user.EmailVerifiedAt, user.LastLoginAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Unique violation
				return fmt.Errorf("%w: %s", service.ErrUserAlreadyExists, pgErr.Detail)
			}
		}
		return fmt.Errorf("failed to create user in db: %w", err)
	}
	return nil
}

func (r *UserRepository) scanUser(row pgx.Row) (*domain.User, error) {
    user := &domain.User{}
    err := row.Scan(
        &user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Status,
        &user.CreatedAt, &user.UpdatedAt, &user.EmailVerifiedAt, &user.LastLoginAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, service.ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to scan user row: %w", err)
    }
    return user, nil
}

func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, status, created_at, updated_at, email_verified_at, last_login_at
                  FROM users WHERE id = $1 AND status != 'deleted'`
	row := r.db.QueryRow(context.Background(), query, id)
	return r.scanUser(row)
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, status, created_at, updated_at, email_verified_at, last_login_at
                  FROM users WHERE email = $1 AND status != 'deleted'`
	row := r.db.QueryRow(context.Background(), query, email)
	return r.scanUser(row)
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, status, created_at, updated_at, email_verified_at, last_login_at
                  FROM users WHERE username = $1 AND status != 'deleted'`
	row := r.db.QueryRow(context.Background(), query, username)
    return r.scanUser(row)
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, status = $4,
                  updated_at = $5, email_verified_at = $6, last_login_at = $7
                  WHERE id = $8`
    user.UpdatedAt = time.Now().UTC()
	_, err := r.db.Exec(context.Background(), query,
		user.Username, user.Email, user.PasswordHash, user.Status,
		user.UpdatedAt, user.EmailVerifiedAt, user.LastLoginAt, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user in db: %w", err)
	}
	return nil
}
