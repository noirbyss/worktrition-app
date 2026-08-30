package service

import "context"

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SavePlan(ctx context.Context, r SaveGeneratedPlanRequest) error {
	if err := r.validate(); err != nil {
		return err
	}

	return s.repo.SavePlan(ctx, r)
}

func (s *Service) GetDayPlan(ctx context.Context, r GetDayPlanRequest) (GetDayPlanResponse, error) {
	if err := r.validate(); err != nil {
		return GetDayPlanResponse{}, err
	}

	return s.repo.GetDayPlan(ctx, r)
}

// CompleteTraining marks the training of the requested day as done and returns the
// training type so the transport layer can grant the matching gamification reward.
func (s *Service) CompleteTraining(ctx context.Context, r CompleteTrainingRequest) (string, error) {
	if err := r.validate(); err != nil {
		return "", err
	}

	return s.repo.CompleteTraining(ctx, r)
}

func (s *Service) GetStats(ctx context.Context, r GetStatsRequest) (GetStatsResponse, error) {
	if err := r.validate(); err != nil {
		return GetStatsResponse{}, err
	}

	completionDates, err := s.repo.GetCompletionDates(ctx, r.UserID)
	if err != nil {
		return GetStatsResponse{}, err
	}

	totalTrainingTime, err := s.repo.GetTotalTrainingTimeSeconds(ctx, r.UserID)
	if err != nil {
		return GetStatsResponse{}, err
	}

	completed, total, err := s.repo.GetActivePlanFulfillment(ctx, r.UserID)
	if err != nil {
		return GetStatsResponse{}, err
	}

	return GetStatsResponse{
		PercentagePlanFulfilled:  planFulfillmentPercentage(completed, total),
		CurrentStreakDays:        currentStreakDays(completionDates),
		TotalTrainingTimeSeconds: totalTrainingTime,
	}, nil
}
