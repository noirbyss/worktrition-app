package grpcclient

import (
	"ai-service/internal/domain"
	"context"

	"github.com/noirbyss/worktrition-app/gen/nutrition-service"
	"go.uber.org/zap"
)

type NutritionClient struct {
	client  nutrition.NutritionServiceClient
	logger 	*zap.SugaredLogger
}

func NewNutritionClient(client nutrition.NutritionServiceClient, logger *zap.SugaredLogger) *NutritionClient {
	return &NutritionClient{
		client: client,
		logger: logger,
	}
}

func (nc *NutritionClient) SaveGeneratedPlan(ctx context.Context, userId string, generationId string, plan *domain.NutriationPlanDTO) error {
	nc.client.SaveGeneratedPlan(ctx, &nutrition.SaveGeneratedPlanRequest{
		UserId: userId,
		GenerationId: generationId,
		Meals: n,
	})//userId, generationId, plan)
	return nil
}