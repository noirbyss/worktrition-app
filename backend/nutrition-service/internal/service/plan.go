package service

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
	UserID       string
	GenerationID string
	Meals        []CreateMealRequest
	NutritionFacts
	WaterGoalMl int32
}

func (cp CreatePlanRequest) validate() error {
	if cp.UserID == "" {
		return ErrEmptyUserID
	}

	if cp.GenerationID == "" {
		return ErrEmptyGenerationID
	}

	if len(cp.Meals) == 0 {
		return ErrEmptyMeals
	}

	for _, meal := range cp.Meals {
		if err := meal.validate(); err != nil {
			return err
		}
	}

	if err := cp.NutritionFacts.validate(); err != nil {
		return err
	}

	if cp.WaterGoalMl <= 0 {
		return ErrInvalidWaterGoal
	}

	return nil
}

type GetDayPlanRequest struct {
	UserID string
	Day    DaysOfWeek
}

func (dpr GetDayPlanRequest) validate() error {
	if dpr.UserID == "" {
		return ErrEmptyUserID
	}

	if !dpr.Day.isValid() {
		return ErrInvalidDayOfWeek
	}

	return nil
}

type GetDayPlanResponse struct {
	Meals []GetMealResponse
}
