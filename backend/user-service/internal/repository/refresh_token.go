package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/domain"
)

type PostgresRefreshTokenRepository struct {
	pool *pgxpool.Pool
}

var _ domain.RefreshTokenRepository = (*PostgresRefreshTokenRepository)(nil)

func NewPostgresRefreshTokenRepository(pool *pgxpool.Pool) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{pool: pool}
}

func (r *PostgresRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	if token == nil {
		return domain.NewValidationError("refresh_token", "is required")
	}

	const query = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id::text, created_at
	`

	err := r.pool.QueryRow(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt).
		Scan(&token.ID, &token.CreatedAt)
	if mappedErr := mapRefreshTokenWriteError(err); mappedErr != nil {
		if mappedErr == domain.ErrUserNotFound || domain.IsValidationError(mappedErr) {
			return mappedErr
		}

		return fmt.Errorf("create refresh token: %w", mappedErr)
	}

	return nil
}

func (r *PostgresRefreshTokenRepository) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	const query = `
		SELECT id::text, user_id::text, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var token domain.RefreshToken
	var revokedAt pgtype.Timestamptz

	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&revokedAt,
		&token.CreatedAt,
	)
	if isNoRows(err) {
		return nil, domain.ErrInvalidToken
	}
	if err != nil {
		return nil, fmt.Errorf("get refresh token by hash: %w", err)
	}

	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
	}

	return &token, nil
}

func (r *PostgresRefreshTokenRepository) RevokeByHash(ctx context.Context, tokenHash string) error {
	const query = `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1 AND revoked_at IS NULL
	`

	commandTag, err := r.pool.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("revoke refresh token by hash: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrInvalidToken
	}

	return nil
}

func (r *PostgresRefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	const query = `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL
	`

	_, err := r.pool.Exec(ctx, query, userID)
	if mappedErr := mapUUIDValidationError(err, "user_id"); mappedErr != nil {
		return mappedErr
	}
	if err != nil {
		return fmt.Errorf("revoke refresh tokens by user id: %w", err)
	}

	return nil
}

func mapRefreshTokenWriteError(err error) error {
	if err == nil {
		return nil
	}
	if isPgErrorCode(err, foreignKeyViolation) {
		return domain.ErrUserNotFound
	}
	if isPgErrorCode(err, uniqueViolation) {
		return domain.ErrInvalidToken
	}
	if isPgErrorCode(err, invalidTextFormat) {
		return domain.NewValidationError("user_id", "must be a valid UUID")
	}

	return err
}
