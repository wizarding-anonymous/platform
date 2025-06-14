package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/russian-steam/auth-service/internal/domain"
    "github.com/russian-steam/auth-service/internal/service" // For ErrRefreshTokenNotFound
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(token *domain.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, session_id, is_revoked, expires_at, created_at)
                  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(context.Background(), query,
		token.ID, token.UserID, token.TokenHash, token.SessionID, token.IsRevoked, token.ExpiresAt, token.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create refresh token in db: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) scanToken(row pgx.Row) (*domain.RefreshToken, error) {
    token := &domain.RefreshToken{}
    err := row.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.SessionID, &token.IsRevoked, &token.ExpiresAt, &token.CreatedAt)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, service.ErrRefreshTokenNotFound
        }
        return nil, fmt.Errorf("failed to scan refresh token row: %w", err)
    }
    return token, nil
}

func (r *RefreshTokenRepository) GetByTokenHash(tokenHash string) (*domain.RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, session_id, is_revoked, expires_at, created_at
                  FROM refresh_tokens WHERE token_hash = $1`
	row := r.db.QueryRow(context.Background(), query, tokenHash)
	return r.scanToken(row)
}

func (r *RefreshTokenRepository) GetByID(id string) (*domain.RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, session_id, is_revoked, expires_at, created_at
                  FROM refresh_tokens WHERE id = $1`
	row := r.db.QueryRow(context.Background(), query, id)
    return r.scanToken(row)
}

func (r *RefreshTokenRepository) SetRevoked(id string, isRevoked bool) error {
	query := `UPDATE refresh_tokens SET is_revoked = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(context.Background(), query, isRevoked, time.Now().UTC(), id)
	if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return service.ErrRefreshTokenNotFound
        }
		return fmt.Errorf("failed to set refresh token revoked status: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) Delete(id string) error {
	query := `DELETE FROM refresh_tokens WHERE id = $1`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		// pgx.ErrNoRows is not typically returned by Exec for DELETE.
		// We rely on RowsAffected if we need to confirm something was deleted.
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	if result.RowsAffected() == 0 {
		return service.ErrRefreshTokenNotFound // Return specific error if no rows were deleted
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteByUserID(userID string) error {
    query := `DELETE FROM refresh_tokens WHERE user_id = $1`
    _, err := r.db.Exec(context.Background(), query, userID)
    if err != nil {
        return fmt.Errorf("failed to delete refresh tokens by user ID: %w", err)
    }
    return nil
}
