package service

import "errors"

var (
	ErrEmptyUserID              = errors.New("user id must not be empty")
	ErrEmptyGenerationID        = errors.New("generation id must not be empty")
	ErrInvalidPlannedMealsCount = errors.New("planned meals must contain exactly 7 days")
	ErrWaterGoalTooLow          = errors.New("water goal must be at least the minimal allowed value")
	ErrInvalidDayOfWeek         = errors.New("day of week is invalid")
	ErrEmptyMealItems           = errors.New("meal items must not be empty")
	ErrCaloriesTooLow           = errors.New("calories must be at least the minimal allowed value")
	ErrProteinTooLow            = errors.New("protein must be at least the minimal allowed value")
	ErrFatTooLow                = errors.New("fat must be at least the minimal allowed value")
	ErrCarbTooLow               = errors.New("carb must be at least the minimal allowed value")
	ErrEmptyName                = errors.New("name must not be empty")
	ErrEmptyRecipe              = errors.New("recipe must not be empty")
	ErrAmountMlTooLow           = errors.New("amount ml must be at least the minimal allowed value")
)
