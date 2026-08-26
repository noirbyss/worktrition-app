package service

type Repository interface {
	// TODO: реализовать
	SavePlan(SaveGeneratedPlanRequest) error
}
