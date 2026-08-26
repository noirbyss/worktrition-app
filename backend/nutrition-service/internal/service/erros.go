package service

import "errors"

var (
	ErrEmptyUserID              = errors.New("value (UserID) can`t be empty")
	ErrEmptyGenerationID        = errors.New("value (GenerationID) can`t be empty")
	ErrInvalidPlannedMealsCount = errors.New("planned meals must contain exactly 7 days")
	ErrWaterGoalTooLow          = errors.New("value (water goal) is below minimal allowed value")
	ErrInvalidDayOfWeek         = errors.New("invalid day of week")
	ErrEmptyMealItems           = errors.New("value (MealItems) can`t be empty")
	ErrCaloriesTooLow           = errors.New("calories is below minimal allowed value")
	ErrProteinTooLow            = errors.New("protein is below minimal allowed value")
	ErrFatTooLow                = errors.New("fat is below minimal allowed value")
	ErrCarbTooLow               = errors.New("carb is below minimal allowed value")
	ErrEmptyName                = errors.New("value (name) can`t be empty")
	ErrEmptyRecipe              = errors.New("value (recipe) cant`t be empty")
)
