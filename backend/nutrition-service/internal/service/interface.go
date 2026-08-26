package service

type Repository interface {
	// TODO: реализовать
	SavePlan(SaveGeneratedPlanRequest) error
	GetDayPlan(GetDayPlanRequest) (GetDayPlanResponse, error)
}
