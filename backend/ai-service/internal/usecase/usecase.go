package usecase

import (
	"ai-service/internal/provider"
	"context"
	"encoding/json"
	"ai-service/internal/domain"
	"go.uber.org/zap"
)

type UseCase struct {
	Provider AIProvider
	logger *zap.SugaredLogger
}

func NewUseCase (provider AIProvider, logger *zap.SugaredLogger) *UseCase {
	return &UseCase{
		Provider: provider,
		logger: logger,
	}
} 

func (us *UseCase) StartGeneration(ctx context.Context, userID string, planType string) (*domain.GeneratedPlanDTO, error) {
	userPrompt :=  provider.UserPrompt //GetInfo(userID)
	resp, err := us.Provider.GeneratePlan(ctx, provider.SystemPrompt, userPrompt)
	if err != nil {
		//....
		us.logger.Errorf("failed to generate response %v", err)
		return nil, err
	}
	dto := domain.GeneratedPlanDTO{}
	if err := json.Unmarshal([]byte(resp), &dto); err != nil {
		us.logger.Errorf("failed to decode json %v", err)
		return nil, err
	}
	// Get
	return  &dto, nil 
}