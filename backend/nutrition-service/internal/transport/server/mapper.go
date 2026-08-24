package server

import (
	"nutrition-service/internal/service"

	pb "github.com/noirbyss/worktrition-app/gen/nutrition-service"
)

func toServiceNutritionFacts(nutritionFactsPb *pb.NutritionFacts) service.NutritionFacts {
	if nutritionFactsPb == nil {
		return service.NutritionFacts{}
	}

	return service.NutritionFacts{
		Calories: nutritionFactsPb.GetCalories(),
		Proteins: nutritionFactsPb.GetProtein(),
		Fats:     nutritionFactsPb.GetFat(),
		Carbs:    nutritionFactsPb.GetCarbs(),
	}
}

func toServiceCreateMealRequest(plannedMealPb *pb.PlannedMeal) service.CreateMealRequest {
	if plannedMealPb == nil {
		return service.CreateMealRequest{}
	}

	return service.CreateMealRequest{
		Day:            service.DaysOfWeek(plannedMealPb.GetDayOfWeek()),
		Name:           plannedMealPb.GetName(),
		Items:          plannedMealPb.GetMealItems(),
		NutritionFacts: toServiceNutritionFacts(plannedMealPb.GetNutritionFacts()),
	}
}

func toServiceCreatePlanRequest(planPb *pb.SaveGeneratedPlanRequest) service.CreatePlanRequest {
	if planPb == nil {
		return service.CreatePlanRequest{}
	}

	plannedMealPb := planPb.GetMeals()
	plannedMealService := make([]service.CreateMealRequest, 0, len(plannedMealPb))

	for _, meal := range plannedMealPb {
		plannedMealService = append(plannedMealService, toServiceCreateMealRequest(meal))
	}

	return service.CreatePlanRequest{
		UserID:         planPb.GetUserId(),
		GenID:          planPb.GetGenerationId(),
		Meals:          plannedMealService,
		NutritionFacts: toServiceNutritionFacts(planPb.NutritionFacts),
		WaterGoalMl:    planPb.GetWaterGoalMl(),
	}
}
