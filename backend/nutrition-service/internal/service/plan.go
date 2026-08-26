package service

const MINIMAL_WATER_GOAL_ML = 0

type DaysOfWeek int32

func (d DaysOfWeek) validate() bool {
	return d >= Monday && d <= Sunday
}

const (
	Unspecified DaysOfWeek = iota
	Monday
	Tuesday
	Wednesday
	Thurday
	Friday
	Saturday
	Sunday
)

type SaveGeneratedPlanRequest struct {
	UserID       string
	GenerationID string
	PlannedMeals []PlannedMealsRequest
	NutritionFacts
	WaterGoalMl int32
}

func (sgpr SaveGeneratedPlanRequest) validate() error {
	if sgpr.UserID == "" {
		return ErrEmptyUserID
	}

	if sgpr.GenerationID == "" {
		return ErrEmptyGenerationID
	}

	if len(sgpr.PlannedMeals) < 7 {
		return ErrInvalidPlannedMealsCount
	}

	for _, plannedMeals := range sgpr.PlannedMeals {
		if err := plannedMeals.validate(); err != nil {
			return err
		}
	}

	if err := sgpr.NutritionFacts.validate(); err != nil {
		return err
	}

	if sgpr.WaterGoalMl < MINIMAL_WATER_GOAL_ML {
		return ErrWaterGoalTooLow
	}

	return nil
}
