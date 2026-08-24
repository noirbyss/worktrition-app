package service

import "errors"

// Validation errors
var (
	ErrEmptyUserID      = errors.New("value (CreatePlanRequest.UserID) is empty")
	ErrEmptyGenerateID  = errors.New("value (CreatePlanRequest.GenID) is empty")
	ErrEmptyMeals       = errors.New("value ([]CreatePlanRequest.Meals) is empty")
	ErrInvalidDayOfWeek = errors.New("value (CreateMealRequest.Meal) is invalid")
	ErrEmptyMealName    = errors.New("value (CreateMealRequest.Name) is empty")
	ErrEmptyMealItems   = errors.New("value ([]CreateMealRequest.Items) is empty")
	ErrInvalidWaterGoal = errors.New("CreatePlanRequest.WaterGoalMl must be greater than zero")
	ErrInvalidCalories  = errors.New("NutritionFacts.Calories must be greater than zero")
	ErrInvalidProteins  = errors.New("NutritionFacts.Proteins must be greater than zero")
	ErrInvalidFats      = errors.New("NutritionFacts.Fats must be greater than zero")
	ErrInvalidCarbs     = errors.New("NutritionFacts.Carbs must be greater than zero")
)
