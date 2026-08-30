package repository

import (
	"context"
	"errors"
	"time"

	"workout-service/internal/service"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) (*PostgresDB, error) {
	if pool == nil {
		return nil, ErrInvalidPool
	}

	return &PostgresDB{pool: pool}, nil
}

func (db *PostgresDB) SavePlan(ctx context.Context, plan service.SaveGeneratedPlanRequest) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
	UPDATE plan_templates
	SET is_active = false
	WHERE user_id = $1 AND is_active = true;
	`, plan.UserID); err != nil {
		return err
	}

	var planID int32
	if err := tx.QueryRow(ctx, `
	INSERT INTO plan_templates (user_id, generation_id, is_active)
	VALUES ($1, $2, true)
	RETURNING id;
	`, plan.UserID, plan.GenerationID).Scan(&planID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrPlanAlreadyExists
		}

		return err
	}

	for _, day := range plan.WorkoutDays {
		var trainingTemplateID int32
		if err := tx.QueryRow(ctx, `
		INSERT INTO training_templates (plan_id, day_of_week, type)
		VALUES ($1, $2, $3)
		RETURNING id;
		`, planID, int16(day.DayOfWeek), day.Type).Scan(&trainingTemplateID); err != nil {
			return err
		}

		for position, exercise := range day.Exercises {
			if _, err := tx.Exec(ctx, `
			INSERT INTO exercises (training_template_id, position, name)
			VALUES ($1, $2, $3);
			`, trainingTemplateID, int16(position), exercise); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (db *PostgresDB) GetDayPlan(ctx context.Context, r service.GetDayPlanRequest) (service.GetDayPlanResponse, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return service.GetDayPlanResponse{}, err
	}
	defer tx.Rollback(ctx)

	trainingTemplateID, trainingType, err := activeTrainingTemplate(ctx, tx, r.UserID, r.DayOfWeek)
	if err != nil {
		return service.GetDayPlanResponse{}, err
	}

	rows, err := tx.Query(ctx, `
	SELECT name
	FROM exercises
	WHERE training_template_id = $1
	ORDER BY position;
	`, trainingTemplateID)
	if err != nil {
		return service.GetDayPlanResponse{}, err
	}
	defer rows.Close()

	exercises := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return service.GetDayPlanResponse{}, err
		}

		exercises = append(exercises, name)
	}

	if err := rows.Err(); err != nil {
		return service.GetDayPlanResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return service.GetDayPlanResponse{}, err
	}

	return service.GetDayPlanResponse{
		DayOfWeek: r.DayOfWeek,
		Type:      trainingType,
		Exercises: exercises,
	}, nil
}

func (db *PostgresDB) CompleteTraining(ctx context.Context, r service.CompleteTrainingRequest) (string, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	trainingTemplateID, trainingType, err := activeTrainingTemplate(ctx, tx, r.UserID, r.DayOfWeek)
	if err != nil {
		return "", err
	}

	if _, err := tx.Exec(ctx, `
	INSERT INTO training_completions (training_template_id, duration_seconds)
	VALUES ($1, $2);
	`, trainingTemplateID, r.DurationSeconds); err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return trainingType, nil
}

func (db *PostgresDB) GetCompletionDates(ctx context.Context, userID string) ([]time.Time, error) {
	rows, err := db.pool.Query(ctx, `
	SELECT DISTINCT tc.completed_at::date
	FROM training_completions tc
	JOIN training_templates tt ON tt.id = tc.training_template_id
	JOIN plan_templates pt ON pt.id = tt.plan_id
	WHERE pt.user_id = $1
	ORDER BY 1 DESC;
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dates := make([]time.Time, 0)
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}

		dates = append(dates, date)
	}

	return dates, rows.Err()
}

func (db *PostgresDB) GetTotalTrainingTimeSeconds(ctx context.Context, userID string) (int32, error) {
	var total int32
	if err := db.pool.QueryRow(ctx, `
	SELECT COALESCE(SUM(tc.duration_seconds), 0)::int
	FROM training_completions tc
	JOIN training_templates tt ON tt.id = tc.training_template_id
	JOIN plan_templates pt ON pt.id = tt.plan_id
	WHERE pt.user_id = $1;
	`, userID).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (db *PostgresDB) GetActivePlanFulfillment(ctx context.Context, userID string) (int, int, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx)

	var (
		planID      int32
		activatedAt time.Time
	)

	if err := tx.QueryRow(ctx, `
	SELECT id, activated_at
	FROM plan_templates
	WHERE user_id = $1 AND is_active = true;
	`, userID).Scan(&planID, &activatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, nil
		}

		return 0, 0, err
	}

	var weeklyTrainings int
	if err := tx.QueryRow(ctx, `
	SELECT COUNT(*)
	FROM training_templates tt
	WHERE tt.plan_id = $1
	  AND EXISTS (SELECT 1 FROM exercises e WHERE e.training_template_id = tt.id);
	`, planID).Scan(&weeklyTrainings); err != nil {
		return 0, 0, err
	}

	var completed int
	if err := tx.QueryRow(ctx, `
	SELECT COUNT(*)
	FROM training_completions tc
	JOIN training_templates tt ON tt.id = tc.training_template_id
	WHERE tt.plan_id = $1;
	`, planID).Scan(&completed); err != nil {
		return 0, 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, 0, err
	}

	elapsedWeeks := int(time.Since(activatedAt).Hours()/(24*7)) + 1
	total := weeklyTrainings * elapsedWeeks

	return completed, total, nil
}

func activeTrainingTemplate(ctx context.Context, tx pgx.Tx, userID string, day service.DaysOfWeek) (int32, string, error) {
	var planID int32
	if err := tx.QueryRow(ctx, `
	SELECT id
	FROM plan_templates
	WHERE user_id = $1 AND is_active = true;
	`, userID).Scan(&planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", ErrPlanNotFound
		}

		return 0, "", err
	}

	var (
		trainingTemplateID int32
		trainingType       string
	)
	if err := tx.QueryRow(ctx, `
	SELECT id, type
	FROM training_templates
	WHERE plan_id = $1 AND day_of_week = $2;
	`, planID, int16(day)).Scan(&trainingTemplateID, &trainingType); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", ErrTrainingNotFound
		}

		return 0, "", err
	}

	return trainingTemplateID, trainingType, nil
}
