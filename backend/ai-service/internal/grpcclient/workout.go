package grpcclient

import (
	"ai-service/internal/domain"
	"context"
	"fmt"

	"github.com/noirbyss/worktrition-app/gen/workout-service"
	"go.uber.org/zap"
)

type WorkoutClient struct {
	client workout.WorkoutServiceClient
	logger *zap.SugaredLogger
}

func NewWorkoutClient(client workout.WorkoutServiceClient, logger *zap.SugaredLogger) *WorkoutClient {
	return &WorkoutClient{
		client: client,
		logger: logger,
	}
}

func (wc *WorkoutClient) SaveGeneratedPlan(
	ctx context.Context,
	userId string,
	generationId string,
	plan []domain.WorkoutPlanDTO,
) error {
	_, err := wc.client.SaveGeneratedPlan(ctx, &workout.SaveGeneratedPlanRequest{
		UserId:       userId,
		GenerationId: generationId,
		WorkoutDays:  mapWorkoutPlanToProto(plan),
	})
	if err != nil {
		wc.logger.Errorf(
			"failed to save generated workout plan for user_id=%s generation_id=%s: %v",
			userId,
			generationId,
			err,
		)
		return fmt.Errorf("workout client save generated plan: %w", err)
	}
	return nil
}

func mapWorkoutPlanToProto(plans []domain.WorkoutPlanDTO) []*workout.WorkoutDayRequest {
	days := make([]*workout.WorkoutDayRequest, 0, len(plans))
	for _, p := range plans {
		days = append(days, &workout.WorkoutDayRequest{
			DayOfWeek: workout.DaysOfWeek(p.Day),
			Type:      p.Type,
			Exercises: p.Exercises,
		})
	}
	return days
}
