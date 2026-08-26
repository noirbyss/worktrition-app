package usecase

import (
	"ai-service/internal/domain"
	"ai-service/internal/provider"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UseCase struct {
	Provider 	AIProvider
	UserClient 	UserClient
	Nutrition	NutritionClient
	Workout 	WorkoutClient
	logger 		*zap.SugaredLogger
}

func NewUseCase (provider AIProvider, userClient UserClient, nutrition NutritionClient,workout WorkoutClient, logger *zap.SugaredLogger) *UseCase {
	return &UseCase{
		Provider: 	provider,
		UserClient:	userClient,
		Nutrition: 	nutrition,
		Workout: 	workout,
		logger: 	logger,
	}
} 

func (us *UseCase) StartGeneration(ctx context.Context, userId string, planType string) (string, error) {
	generationId := uuid.New().String()
	go us.generatePlan(userId, generationId, planType)
	return  generationId, nil 
}

func (us *UseCase) generatePlan(generationId, userId, planType string) {
		bgCtx := context.Background()
		profile,err := us.UserClient.GetProfile(bgCtx, userId)
		if err != nil {
			us.logger.Errorf("failed %v", err)
			return 
		}
		userPrompt := MapUserProfilePrompt(profile)
		
		resp, err := us.Provider.GeneratePlan(bgCtx, provider.SystemPrompt, userPrompt)
		if err != nil {
			us.logger.Errorf("failed to generate response %v", err)
			return 
		}
		plan := domain.GeneratedPlanDTO{}
		if err := json.Unmarshal([]byte(resp), &plan); err != nil {
			us.logger.Errorf("failed to decode json %v", err)
			return 
		}
		us.Nutrition.SaveGeneratedPlan(bgCtx, userId, generationId, plan.Nutrition)
		us.Workout.SaveGeneratedPlan(bgCtx, userId, generationId, plan.Workouts)
	}