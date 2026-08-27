package repository

import (
	"ai-service/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, generation *domain.Generation) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO ai_generation_requests (id, user_id, plan_type, status, error_message)
		 VALUES ($1, $2, $3, $4, $5)`,
		generation.ID,
		generation.UserID,
		generation.PlanType,
		generation.Status,
		generation.ErrorMessage,
	)
	if err != nil {
		return fmt.Errorf("insert generation request: %w", err)
	}

	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO generated_plan_versions (generation_id, user_id, plan_type) VALUES ($1, $2, $3)`,
		generation.ID,
		generation.UserID,
		generation.PlanType,
	)
	if err != nil {
		return fmt.Errorf("insert generated plan version: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, generationID string) (*domain.Generation, error) {
	generation := &domain.Generation{}
	var errorMessage sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, plan_type, status, error_message
		 FROM ai_generation_requests WHERE id = $1`,
		generationID,
	).Scan(
		&generation.ID,
		&generation.UserID,
		&generation.PlanType,
		&generation.Status,
		&errorMessage,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrGenerationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("select generation request: %w", err)
	}

	if errorMessage.Valid {
		generation.ErrorMessage = errorMessage.String
	}

	return generation, nil
}

func (r *PostgresRepository) UpdateStatus(
	ctx context.Context,
	generationID string,
	status domain.GenerationStatus,
	errorMessage string,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE ai_generation_requests
		 SET status = $1, error_message = $2, updated_at = now()
		 WHERE id = $3`,
		status,
		errorMessage,
		generationID,
	)
	if err != nil {
		return fmt.Errorf("update generation status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrGenerationNotFound
	}

	return nil
}

func (r *PostgresRepository) SavePrompt(ctx context.Context, generationID, systemPrompt, userPrompt string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO ai_prompt_versions (generation_id, system_prompt, user_prompt) VALUES ($1, $2, $3)`,
		generationID,
		systemPrompt,
		userPrompt,
	)
	if err != nil {
		return fmt.Errorf("insert prompt version: %w", err)
	}
	return nil
}

func (r *PostgresRepository) SaveRawResult(ctx context.Context, generationID, rawResponse string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO ai_generation_results (generation_id, raw_response) VALUES ($1, $2)`,
		generationID,
		rawResponse,
	)
	if err != nil {
		return fmt.Errorf("insert generation result: %w", err)
	}
	return nil
}
