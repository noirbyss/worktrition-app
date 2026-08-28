package repository

import (
	"context"
	"errors"
	"nutrition-service/internal/service"

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

	var planTemplateID int32
	if err := tx.QueryRow(ctx, `
	INSERT INTO plan_templates (user_id, generation_id, calories, protein, fat, carb, water_goal, is_active)
	VALUES($1, $2, $3, $4, $5, $6, $7, true)
	RETURNING id;
	`, plan.UserID,
		plan.GenerationID,
		plan.Calories,
		plan.Protein,
		plan.Fat,
		plan.Carb,
		plan.WaterGoalMl).Scan(&planTemplateID); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrPlanAlreadyExists
		}

		return err
	}

	for _, pm := range plan.PlannedMeals {
		var mealTemplateID int32
		if err := tx.QueryRow(ctx, `
		INSERT INTO meal_templates (plan_id, day_of_week, calories, protein, fat, carb)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id;
		`, planTemplateID,
			pm.DayOfWeek,
			pm.Calories,
			pm.Protein,
			pm.Fat,
			pm.Carb).Scan(&mealTemplateID); err != nil {
			return err
		}

		for _, item := range pm.MealItems {
			_, err := tx.Exec(ctx, `
			INSERT INTO meal_items (meal_template_id, name, recipe, calories, protein, fat, carb)
			VALUES($1, $2, $3, $4, $5, $6, $7)
			`, mealTemplateID, item.Name, item.Recipe, item.Calories, item.Protein, item.Fat, item.Carb)
			if err != nil {
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

	var (
		planTemplateID int32
		waterGoal      int32
	)

	if err := tx.QueryRow(ctx, `
	SELECT id, water_goal
	FROM plan_templates
	WHERE user_id = $1 AND is_active = true;
	`, r.UserID).Scan(&planTemplateID, &waterGoal); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return service.GetDayPlanResponse{}, ErrPlanNotFound
		}
		return service.GetDayPlanResponse{}, err
	}

	var (
		mealTemplateID int32
		nutritionFacts service.NutritionFacts
	)

	if err := tx.QueryRow(ctx, `
	SELECT id, calories, protein, fat, carb
	FROM meal_templates
	WHERE plan_id = $1 AND day_of_week = $2;
	`, planTemplateID, r.DayOfWeek).Scan(&mealTemplateID,
		&nutritionFacts.Calories,
		&nutritionFacts.Protein,
		&nutritionFacts.Fat,
		&nutritionFacts.Carb); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return service.GetDayPlanResponse{}, ErrPlanNotFound
		}
		return service.GetDayPlanResponse{}, err
	}

	meals := make([]service.MealItemsResponse, 0)

	rows, err := tx.Query(ctx, `
	SELECT id, name, recipe
	FROM meal_items
	WHERE meal_template_id = $1;
	`, mealTemplateID)
	if err != nil {
		return service.GetDayPlanResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var meal service.MealItemsResponse

		if err := rows.Scan(&meal.ID, &meal.Name, &meal.Recipe); err != nil {
			return service.GetDayPlanResponse{}, err
		}

		meals = append(meals, meal)
	}

	if err := rows.Err(); err != nil {
		return service.GetDayPlanResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return service.GetDayPlanResponse{}, err
	}

	return service.GetDayPlanResponse{
		MealItems:      meals,
		NutritionFacts: nutritionFacts,
		WaterGoalMl:    waterGoal,
	}, nil
}

func (db *PostgresDB) CompleteMeal(ctx context.Context, r service.CompleteMealRequest) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool
	if err := tx.QueryRow(ctx, `
	SELECT EXISTS(
		SELECT 1
		FROM meal_items mi
		JOIN meal_templates mt ON mt.id = mi.meal_template_id
		JOIN plan_templates pt ON pt.id = mt.plan_id
		WHERE mi.id = $1 AND pt.user_id = $2 AND pt.is_active = true
	);
	`, r.MealItemID, r.UserID).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		return ErrMealItemNotFound
	}

	if _, err := tx.Exec(ctx, `
	INSERT INTO meal_completions (meal_item_id)
	VALUES($1);
	`, r.MealItemID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (db *PostgresDB) CompleteWater(ctx context.Context, r service.CompleteWaterRequest) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var planID int32

	if err := tx.QueryRow(ctx, `
	SELECT id
	FROM plan_templates
	WHERE user_id = $1 AND is_active = true
	`, r.UserID).Scan(&planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPlanNotFound
		}

		return err
	}

	if _, err := tx.Exec(ctx, `
	INSERT INTO water_completions (plan_id, amount_ml)
	VALUES($1, $2);
	`, planID, r.AmountMl); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (db *PostgresDB) GetNutritionHistory(ctx context.Context, userID string) ([]service.NutritionDayRecord, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
	SELECT mt.calories, mt.protein, mt.fat, mt.carb,
	SUM(mi.calories) AS consumed_calories,
	SUM(mi.protein) AS consumed_protein,
	SUM(mi.fat) AS consumed_fat,
	SUM(mi.carb) AS consumed_carb

	FROM meal_completions mc

	JOIN meal_items mi ON mi.id = mc.meal_item_id
	JOIN meal_templates mt ON mt.id = mi.meal_template_id
	JOIN plan_templates pt ON pt.id = mt.plan_id

	WHERE pt.user_id = $1

	GROUP BY mc.completed_at::date, mt.id, mt.calories, mt.protein, mt.fat, mt.carb;
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]service.NutritionDayRecord, 0)
	for rows.Next() {
		var record service.NutritionDayRecord

		if err := rows.Scan(
			&record.Target.Calories, &record.Target.Protein, &record.Target.Fat, &record.Target.Carb,
			&record.Consumed.Calories, &record.Consumed.Protein, &record.Consumed.Fat, &record.Consumed.Carb,
		); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return records, nil
}

func (db *PostgresDB) GetWaterHistory(ctx context.Context, userID string) ([]service.WaterDayRecord, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
	SELECT pt.water_goal ,
	SUM(wc.amount_ml)::int AS consumed_amount

	FROM water_completions wc

	JOIN plan_templates pt ON pt.id = wc.plan_id

	WHERE pt.user_id = $1

	GROUP BY wc.completed_at::date, pt.id, pt.water_goal
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]service.WaterDayRecord, 0)

	for rows.Next() {
		var record service.WaterDayRecord

		if err := rows.Scan(&record.GoalMl, &record.ConsumedMl); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return records, nil
}
