package service

import (
	"context"
	"time"
)

type Repository interface {
	SavePlan(ctx context.Context, r SaveGeneratedPlanRequest) error
	GetDayPlan(ctx context.Context, r GetDayPlanRequest) (GetDayPlanResponse, error)
	CompleteTraining(ctx context.Context, r CompleteTrainingRequest) (trainingType string, err error)

	GetCompletionDates(ctx context.Context, userID string) ([]time.Time, error)
	GetTotalTrainingTimeSeconds(ctx context.Context, userID string) (int32, error)
	GetActivePlanFulfillment(ctx context.Context, userID string) (completed int, total int, err error)
}
