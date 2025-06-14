package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/google/uuid" // Added import
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/russian-steam/auth-service/internal/domain"
    "github.com/russian-steam/auth-service/internal/service" // For ErrVerificationCodeNotFound
)

type VerificationCodeRepository struct {
	db *pgxpool.Pool
}

func NewVerificationCodeRepository(db *pgxpool.Pool) *VerificationCodeRepository {
	return &VerificationCodeRepository{db: db}
}

func (r *VerificationCodeRepository) Create(vc *domain.VerificationCode) error {
	query := `INSERT INTO verification_codes (id, user_id, type, code_hash, target, expires_at, created_at)
                  VALUES ($1, $2, $3, $4, $5, $6, $7)`
    if vc.ID == "" {
        vc.ID = uuid.NewString()
    }
	_, err := r.db.Exec(context.Background(), query,
		vc.ID, vc.UserID, vc.Type, vc.CodeHash, vc.Target, vc.ExpiresAt, vc.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create verification code in db: %w", err)
	}
	return nil
}

func (r *VerificationCodeRepository) FindByCodeHash(codeHash string, vcType domain.VerificationType) (*domain.VerificationCode, error) {
	query := `SELECT id, user_id, type, code_hash, target, expires_at, created_at
                  FROM verification_codes WHERE code_hash = $1 AND type = $2`
	row := r.db.QueryRow(context.Background(), query, codeHash, vcType)
	vc := &domain.VerificationCode{}
	err := row.Scan(&vc.ID, &vc.UserID, &vc.Type, &vc.CodeHash, &vc.Target, &vc.ExpiresAt, &vc.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, service.ErrVerificationCodeNotFound
		}
		return nil, fmt.Errorf("failed to scan verification code row: %w", err)
	}
	return vc, nil
}

func (r *VerificationCodeRepository) Delete(id string) error {
	query := `DELETE FROM verification_codes WHERE id = $1`
	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete verification code: %w", err)
	}
    if result.RowsAffected() == 0 {
        // Optionally, return an error if no rows were affected, though for delete it might not be critical
        // For example: return service.ErrVerificationCodeNotFound
    }
	return nil
}

func (r *VerificationCodeRepository) DeleteByUserIDAndType(userID string, vcType domain.VerificationType) error {
    query := `DELETE FROM verification_codes WHERE user_id = $1 AND type = $2`
    _, err := r.db.Exec(context.Background(), query, userID, vcType)
    if err != nil {
        return fmt.Errorf("failed to delete verification codes by user ID and type: %w", err)
    }
    return nil
}
