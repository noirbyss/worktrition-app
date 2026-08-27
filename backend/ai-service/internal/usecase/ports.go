package usecase

import (
	"ai-service/internal/domain"
	"context"

)

type AIProvider interface {
	GeneratePlan(ctx context.Context, systemPrompt, userPromt string) (string, error)
}
type UserClient interface {
	GetProfile(ctx context.Context, UserId string) (*domain.UserProfile, error) 
}

type NutritionClient interface {
	SaveGeneratedPlan(ctx context.Context, userId string, generationId string, plan []domain.NutritionPlanDTO, waterMl int) error
}

type WorkoutClient interface {
	SaveGeneratedPlan(ctx context.Context, user_id string, generationId string, plan []domain.WorkoutPlanDTO) error
}
