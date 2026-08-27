package service

import "context"

type Service struct {
	repo Repository
}

func (s *Service) SavePlan(ctx context.Context, r SaveGeneratedPlanRequest) error {
	if err := r.validate(); err != nil {
		return err
	}

	if err := s.repo.SavePlan(ctx, r); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetDayPlan(ctx context.Context, r GetDayPlanRequest) (GetDayPlanResponse, error) {
	if err := r.validate(); err != nil {
		return GetDayPlanResponse{}, err
	}

	dayPlan, err := s.repo.GetDayPlan(ctx, r)
	if err != nil {
		return GetDayPlanResponse{}, err
	}

	return dayPlan, nil
}

func (s *Service) CompleteMeal(ctx context.Context, r CompleteMealRequest) error {
	if err := r.validate(); err != nil {
		return err
	}

	if err := s.repo.CompleteMeal(ctx, r); err != nil {
		return err
	}

	return nil
}

func (s *Service) CompleteWater(ctx context.Context, r CompleteWaterRequest) error {
	if err := r.validate(); err != nil {
		return err
	}

	if err := s.repo.CompleteWater(ctx, r); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetStats(ctx context.Context, r GetStatsRequest) (GetStatsResponse, error) {
	if err := r.validate(); err != nil {
		return GetStatsResponse{}, err
	}

	nutritionRecords, err := s.repo.GetNutritionHistory(ctx, r.UserID)
	if err != nil {
		return GetStatsResponse{}, err
	}

	waterRecords, err := s.repo.GetWaterHistory(ctx, r.UserID)
	if err != nil {
		return GetStatsResponse{}, err
	}

	completed, total, err := s.repo.GetActivePlanFulfillment(ctx, r.UserID)
	if err != nil {
		return GetStatsResponse{}, err
	}

	return GetStatsResponse{
		PercentageComplianceNutritionFacts: nutritionCompliancePercentage(nutritionRecords),
		PercentagePlanFulfilled:            planFulfillmentPercentage(completed, total),
		PercentageWaterStandardFulfillment: waterCompliancePercentage(waterRecords),
	}, nil
}
