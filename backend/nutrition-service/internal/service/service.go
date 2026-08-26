package service

import "context"

type Service struct {
	repo Repository
}

func (s *Service) SavePlan(ctx context.Context, r SaveGeneratedPlanRequest) error {
	if err := r.validate(); err != nil {
		return err
	}

	if err := s.repo.SavePlan(r); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetDayPlan(ctx context.Context, r GetDayPlanRequest) (GetDayPlanResponse, error) {
	if err := r.validate(); err != nil {
		return GetDayPlanResponse{}, err
	}

	dayPlan, err := s.repo.GetDayPlan(r)
	if err != nil {
		return GetDayPlanResponse{}, err
	}

	return dayPlan, nil
}
