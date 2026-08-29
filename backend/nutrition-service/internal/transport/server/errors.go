package server

import "errors"

var (
	ErrNilPointerService = errors.New("*service.Service is nil")
	ErrNilPointerClient  = errors.New("*client.NutritionServiceClient is nil")
)
