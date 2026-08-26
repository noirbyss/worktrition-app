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
