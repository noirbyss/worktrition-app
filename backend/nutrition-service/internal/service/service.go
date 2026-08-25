package service

import "context"

type NutritionService struct {
	repo Repository
}

func New(repository Repository) *NutritionService {
	return &NutritionService{repo: repository}
}

func (ns *NutritionService) SavePlan(ctx context.Context, plan CreatePlanRequest) error {
	if err := plan.validate(); err != nil {
		return err
	}

	if err := ns.repo.SavePlan(ctx, plan); err != nil {
		return err
	}

	return nil
}

func (ns *NutritionService) GetDayPlan(ctx context.Context, request GetDayPlanRequest) (GetDayPlanResponse, error) {
	if err := request.validate(); err != nil {
		return GetDayPlanResponse{}, err
	}

	meals, err := ns.repo.GetDayPlan(ctx, request.UserID, request.Day)
	if err != nil {
		return GetDayPlanResponse{}, err
	}

	dayPlan := GetDayPlanResponse{Meals: meals}

	return dayPlan, nil
}

func (ns *NutritionService) CompleteMeal(ctx context.Context, request CompleteMealRequest) error {
	if err := request.validate(); err != nil {
		return err
	}

	return ns.repo.CompleteMeal(ctx, request.UserID, request.MealID)
}

func (ns *NutritionService) CompleteWater(ctx context.Context, request CompleteWaterRequest) error {
	if err := request.validate(); err != nil {
		return err
	}

	return ns.repo.CompleteWater(ctx, request.UserID, request.AmountMl)
}
