package usecase

import (
	"ai-service/internal/domain"
	"ai-service/internal/provider"
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UseCase struct {
	Provider         AIProvider
	UserClient       UserClient
	NutritionClient  NutritionClient
	WorkoutClient    WorkoutClient
	GenerationRepo   GenerationRepository
	logger           *zap.SugaredLogger
}

func NewUseCase(
	provider AIProvider,
	userClient UserClient,
	nutritionClient NutritionClient,
	workoutClient WorkoutClient,
	generationRepo GenerationRepository,
	logger *zap.SugaredLogger,
) *UseCase {
	return &UseCase{
		Provider:        provider,
		UserClient:      userClient,
		NutritionClient: nutritionClient,
		WorkoutClient:   workoutClient,
		GenerationRepo:  generationRepo,
		logger:          logger,
	}
}

func (us *UseCase) StartGeneration(ctx context.Context, userId string, planType string) (string, error) {
	if userId == "" {
		return "", errors.New("user_id is required")
	}

	generationId := uuid.New().String()
	generation := &domain.Generation{
		ID:       generationId,
		UserID:   userId,
		PlanType: planType,
		Status:   domain.GenerationStatusPending,
	}
	if err := us.GenerationRepo.Create(ctx, generation); err != nil {
		return "", err
	}

	go us.generatePlan(generationId, userId, planType)
	return generationId, nil
}

func (us *UseCase) GetGenerationStatus(ctx context.Context, generationId string) (*domain.Generation, error) {
	return us.GenerationRepo.GetByID(ctx, generationId)
}

func (us *UseCase) generatePlan(generationId, userId, planType string) {
	bgCtx := context.Background()

	profile, err := us.UserClient.GetProfile(bgCtx, userId)
	if err != nil {
		us.logger.Errorw(
			"failed to get user profile",
			"userId", userId,
			"generationId", generationId,
			"error", err,
		)
		us.markFailed(bgCtx, generationId, "failed to get user profile")
		return
	}

	userPrompt := MapUserProfilePrompt(profile)
	if err := us.GenerationRepo.SavePrompt(bgCtx, generationId, provider.SystemPrompt, userPrompt); err != nil {
		us.logger.Errorw(
			"failed to save prompt version",
			"userId", userId,
			"generationId", generationId,
			"error", err,
		)
		us.markFailed(bgCtx, generationId, "failed to save prompt version")
		return
	}

	resp, err := us.Provider.GeneratePlan(bgCtx, provider.SystemPrompt, userPrompt)
	if err != nil {
		us.logger.Errorw(
			"failed to generate response",
			"userId", userId,
			"generationId", generationId,
			"error", err,
		)
		us.markFailed(bgCtx, generationId, "failed to generate response")
		return
	}

	if err := us.GenerationRepo.SaveRawResult(bgCtx, generationId, resp); err != nil {
		us.logger.Errorw(
			"failed to save generation result",
			"userId", userId,
			"generationId", generationId,
			"error", err,
		)
		us.markFailed(bgCtx, generationId, "failed to save generation result")
		return
	}

	plan := domain.GeneratedPlanDTO{}
	if err := json.Unmarshal([]byte(resp), &plan); err != nil {
		us.logger.Errorw(
			"failed to decode ai response",
			"userId", userId,
			"generationId", generationId,
			"error", err,
		)
		us.markFailed(bgCtx, generationId, "failed to decode ai response")
		return
	}

	if domain.ShouldSaveNutrition(planType) {
		waterMl := plan.WaterMl
		if profile.WaterGoalMl > 0 {
			waterMl = int(profile.WaterGoalMl)
		}

		if err := us.NutritionClient.SaveGeneratedPlan(bgCtx, userId, generationId, plan.Nutrition, waterMl); err != nil {
			us.logger.Errorw(
				"failed to save nutrition plan",
				"userId", userId,
				"generationId", generationId,
				"error", err,
			)
			us.markFailed(bgCtx, generationId, "failed to save nutrition plan")
			return
		}
	}

	if domain.ShouldSaveWorkout(planType) {
		if us.WorkoutClient == nil {
			us.logger.Errorw(
				"workout client is not configured",
				"userId", userId,
				"generationId", generationId,
			)
			us.markFailed(bgCtx, generationId, "workout client is not configured")
			return
		}

		if err := us.WorkoutClient.SaveGeneratedPlan(bgCtx, userId, generationId, plan.Workouts); err != nil {
			us.logger.Errorw(
				"failed to save workout plan",
				"userId", userId,
				"generationId", generationId,
				"error", err,
			)
			us.markFailed(bgCtx, generationId, "failed to save workout plan")
			return
		}
	}

	if err := us.GenerationRepo.UpdateStatus(bgCtx, generationId, domain.GenerationStatusDone, ""); err != nil {
		us.logger.Errorw(
			"failed to update generation status",
			"userId", userId,
			"generationId", generationId,
			"error", err,
		)
	}
}

func (us *UseCase) markFailed(ctx context.Context, generationId, message string) {
	if err := us.GenerationRepo.UpdateStatus(ctx, generationId, domain.GenerationStatusFailed, message); err != nil {
		us.logger.Errorw(
			"failed to mark generation as failed",
			"generationId", generationId,
			"error", err,
		)
	}
}
