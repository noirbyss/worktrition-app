package repository

import "errors"

var ErrInvalidPool = errors.New("*pgx.pool is nil")
