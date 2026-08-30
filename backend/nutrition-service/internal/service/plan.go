package service

const minimalWaterGoalMl = 0

type DaysOfWeek int32

func (d DaysOfWeek) validate() bool {
	return d >= Monday && d <= Sunday
}

const (
	Unspecified DaysOfWeek = iota
	Monday
	Tuesday
	Wednesday
	Thursday
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

	if sgpr.WaterGoalMl < minimalWaterGoalMl {
		return ErrWaterGoalTooLow
	}

	return nil
}

type GetDayPlanRequest struct {
	UserID    string
	DayOfWeek DaysOfWeek
}

func (gdpr GetDayPlanRequest) validate() error {
	if gdpr.UserID == "" {
		return ErrEmptyUserID
	}

	if !gdpr.DayOfWeek.validate() {
		return ErrInvalidDayOfWeek
	}

	return nil
}

type MealItemsResponse struct {
	ID          int32
	Name        string
	Recipe      string
	IsCompleted bool
	NutritionFacts
}

type GetDayPlanResponse struct {
	MealItems []MealItemsResponse
	NutritionFacts
	WaterGoalMl int32
}
