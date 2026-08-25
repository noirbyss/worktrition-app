package service

import "context"

type Repository interface {
	SavePlan(ctx context.Context, plan CreatePlanRequest) error
	GetDayPlan(ctx context.Context, userID string, day DaysOfWeek) ([]GetMealResponse, error)
	CompleteMeal(ctx context.Context, userID string, mealID int32) error
	CompleteWater(ctx context.Context, userID string, amountMl int32) error
	//GetStats(ctx context.Context, userID string) (StatsRawData, error)
}
