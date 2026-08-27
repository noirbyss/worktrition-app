package repository

import "errors"

var (
	ErrInvalidPool       = errors.New("*pgx.pool is nil")
	ErrPlanAlreadyExists = errors.New("plan already exists")
	ErrPlanNotFound      = errors.New("plan for user not found")
	ErrMealItemNotFoun   = errors.New("meal item not found")
)
