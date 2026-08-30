package service

import "errors"

var (
	ErrEmptyUserID             = errors.New("user id must not be empty")
	ErrEmptyGenerationID       = errors.New("generation id must not be empty")
	ErrEmptyType               = errors.New("training type must not be empty")
	ErrEmptyExercise           = errors.New("exercise name must not be empty")
	ErrInvalidWorkoutDaysCount = errors.New("workout days count must be between 1 and 7")
	ErrDuplicateDayOfWeek      = errors.New("workout days must not contain duplicate day of week")
	ErrInvalidDayOfWeek        = errors.New("day of week is invalid")
	ErrInvalidDuration         = errors.New("duration seconds must be greater than zero")
)
