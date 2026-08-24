package service

import "context"

type NutritionService struct {
	repo Repository
}

func New(repository Repository) *NutritionService {
	return &NutritionService{repo: repository}
}

func (ns *NutritionService) SavePlan(ctx context.Context, plan CreatePlanRequest) error {
	if err := plan.Validate(); err != nil {
		return err
	}

	if err := ns.repo.SavePlan(ctx, plan); err != nil {
		return err
	}

	return nil
}

type DaysOfWeek int32

func (d DaysOfWeek) isValid() bool {
	return d >= MONDAY && d <= SUNDAY
}

const (
	UNSPECIFIED DaysOfWeek = iota
	MONDAY
	TUESDAY
	WEDNESDAY
	THURSDAY
	FRIDAY
	SATURDAY
	SUNDAY
)

type CreatePlanRequest struct {
	UserID string
	GenID  string
	Meals  []CreateMealRequest
	NutritionFacts
	WaterGoalMl int32
}

func (cp *CreatePlanRequest) Validate() error {
	if cp.UserID == "" {
		return ErrEmptyUserID
	}

	if cp.GenID == "" {
		return ErrEmptyGenerateID
	}

	if len(cp.Meals) == 0 {
		return ErrEmptyMeals
	}

	for _, meal := range cp.Meals {
		if err := meal.Validate(); err != nil {
			return err
		}
	}

	if err := cp.NutritionFacts.Validate(); err != nil {
		return err
	}

	if cp.WaterGoalMl <= 0 {
		return ErrInvalidWaterGoal
	}

	return nil
}
