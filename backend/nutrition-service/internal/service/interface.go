package service

import "context"

type Repository interface {
	SavePlan(ctx context.Context, r SaveGeneratedPlanRequest) error
	GetDayPlan(ctx context.Context, r GetDayPlanRequest) (GetDayPlanResponse, error)
	CompleteMeal(ctx context.Context, r CompleteMealRequest) error
	CompleteWater(ctx context.Context, r CompleteWaterRequest) error

	GetNutritionHistory(ctx context.Context, userID string) ([]NutritionDayRecord, error)
	GetWaterHistory(ctx context.Context, userID string) ([]WaterDayRecord, error)
	GetActivePlanFulfillment(ctx context.Context, userID string) (completed int, total int, err error)
}
