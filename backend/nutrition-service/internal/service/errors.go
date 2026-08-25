package service

import "errors"

// Validation errors
var (
	ErrEmptyUserID       = errors.New("value (UserID) is empty")
	ErrEmptyGenerationID = errors.New("value (GenerationID) is empty")
	ErrEmptyMeals        = errors.New("value (Meals) is empty")
	ErrInvalidDayOfWeek  = errors.New("value (Meal) is invalid")
	ErrEmptyMealName     = errors.New("value (Name) is empty")
	ErrEmptyMealItems    = errors.New("value (Items) is empty")
	ErrInvalidWaterGoal  = errors.New("value (WaterGoalMl) must be greater than zero")
	ErrInvalidCalories   = errors.New("value (Calories) must be greater than zero")
	ErrInvalidProtein    = errors.New("value (Protein) must be greater than zero")
	ErrInvalidFat        = errors.New("value (Fat) must be greater than zero")
	ErrInvalidCarb       = errors.New("value (Carb) must be greater than zero")
	ErrInvalidMealID     = errors.New("value (MealID) must be greater than zero")
	ErrInvalidWaterMl    = errors.New("value (WaterMl) must be greater than zero")
)
