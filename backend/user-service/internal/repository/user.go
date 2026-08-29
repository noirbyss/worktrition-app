package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/domain"
)

const (
	uniqueViolation     = "23505"
	foreignKeyViolation = "23503"
	invalidTextFormat   = "22P02"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

var _ domain.UserRepository = (*PostgresUserRepository)(nil)

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) (string, error) {
	if user == nil {
		return "", domain.NewValidationError("user", "is required")
	}

	const query = `
		INSERT INTO users (name, email, password_hash, birth_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, profile_completed, created_at, updated_at
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.BirthDate,
	).Scan(&user.ID, &user.ProfileCompleted, &user.CreatedAt, &user.UpdatedAt)
	if mappedErr := mapUserWriteError(err); mappedErr != nil {
		if mappedErr == domain.ErrUserAlreadyExists {
			return "", mappedErr
		}

		return "", fmt.Errorf("create user: %w", mappedErr)
	}

	return user.ID, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT id::text, name, email, password_hash, birth_date, profile_completed, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`

	user, err := r.get(ctx, query, email)
	if isNoRows(err) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const query = `
		SELECT id::text, name, email, password_hash, birth_date, profile_completed, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user, err := r.get(ctx, query, id)
	if isNoRows(err) {
		return nil, domain.ErrUserNotFound
	}
	if mappedErr := mapUUIDValidationError(err, "user_id"); mappedErr != nil {
		return nil, mappedErr
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepository) SetProfileCompleted(ctx context.Context, id string, completed bool) error {
	const query = `
		UPDATE users
		SET profile_completed = $2
		WHERE id = $1
	`

	commandTag, err := r.pool.Exec(ctx, query, id, completed)
	if mappedErr := mapUUIDValidationError(err, "user_id"); mappedErr != nil {
		return mappedErr
	}
	if err != nil {
		return fmt.Errorf("set profile completed: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *PostgresUserRepository) get(ctx context.Context, query string, args ...any) (*domain.User, error) {
	var user domain.User

	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.BirthDate,
		&user.ProfileCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func isPgErrorCode(err error, code string) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == code
}

func mapUserWriteError(err error) error {
	if err == nil {
		return nil
	}
	if isPgErrorCode(err, uniqueViolation) {
		return domain.ErrUserAlreadyExists
	}

	return err
}

func mapProfileWriteError(err error) error {
	if err == nil {
		return nil
	}
	if isPgErrorCode(err, foreignKeyViolation) {
		return domain.ErrUserNotFound
	}

	return err
}

func mapUUIDValidationError(err error, field string) error {
	if err == nil {
		return nil
	}
	if isPgErrorCode(err, invalidTextFormat) {
		return domain.NewValidationError(field, "must be a valid UUID")
	}

	return nil
}
