package repository

import (
	"context"
	"nutrition-service/internal/service"

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
	INSER INTO plan_templates (user_id, generation_id, calories, protein, fat, carb, water_goal, is_active)
	VALUES($1, $2, $3, $4, $5, $6, $7, true)
	RETURNING id;
	`, plan.UserID,
		plan.GenerationID,
		plan.Calories,
		plan.Protein,
		plan.Fat,
		plan.Carb,
		plan.WaterGoalMl).Scan(&planTemplateID); err != nil {
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
			`, mealTemplateID, item.Name, item.Calories, item.Protein, item.Fat, item.Carb)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
