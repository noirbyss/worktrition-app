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
	Provider 			AIProvider
	UserClient		 	UserClient
	NutritionClient 	NutritionClient
	WorkoutClient  		WorkoutClient
	logger 				*zap.SugaredLogger
}

func NewUseCase (provider AIProvider, userClient UserClient, nutritionClient NutritionClient,workoutClient WorkoutClient, logger *zap.SugaredLogger) *UseCase {
	return &UseCase{
		Provider: 			provider,
		UserClient:			userClient,
		NutritionClient:	nutritionClient,
		WorkoutClient: 		workoutClient,
		logger: 			logger,
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
			us.logger.Errorw(
				"failed to get user profile",
				"userId", 		userId, 
				"generationId", generationId,
			 	"error",		err,
			)
			return 
		}
		userPrompt := MapUserProfilePrompt(profile)
		
		resp, err := us.Provider.GeneratePlan(bgCtx, provider.SystemPrompt, userPrompt)
		if err != nil {
			us.logger.Errorw(
				"failed to generate response", 
				"userId", 		userId, 
				"generationId", generationId,
			 	"error",		err,
			)
			return 
		}
		plan := domain.GeneratedPlanDTO{}
		if err := json.Unmarshal([]byte(resp), &plan); err != nil {
			us.logger.Errorw(
				"failed to decode ai response",
				"userId", 		userId, 
				"generationId", generationId,
			 	"error",		err,
			)
			return 
		}
		us.NutritionClient.SaveGeneratedPlan(bgCtx, userId, generationId, plan.Nutrition, plan.WaterMl)
		us.WorkoutClient.SaveGeneratedPlan(bgCtx, userId, generationId, plan.Workouts)
	}