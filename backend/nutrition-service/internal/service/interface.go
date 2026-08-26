package service

type Repository interface {
	// TODO: реализовать
	SavePlan(SaveGeneratedPlanRequest) error
	GetDayPlan(GetDayPlanRequest) (GetDayPlanResponse, error)
	CompleteMeal(CompleteMealRequest) error
	CompleteWater(CompleteWaterRequest) error
}
