package repository

import "errors"

var (
	ErrInvalidPool       = errors.New("*pgxpool.Pool is nil")
	ErrPlanAlreadyExists = errors.New("plan already exists")
	ErrPlanNotFound      = errors.New("plan for user not found")
	ErrTrainingNotFound  = errors.New("training for the requested day not found")
)
