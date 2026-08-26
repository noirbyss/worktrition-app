package server

import (
	"nutrition-service/internal/service"

	pb "github.com/noirbyss/worktrition-app/gen/nutrition-service"
)

func toServiceNutritionFacts(pbr *pb.NutritionFacts) service.NutritionFacts {
	return service.NutritionFacts{
		Calories: pbr.GetCalories(),
		Protein:  pbr.GetProtein(),
		Fat:      pbr.GetFat(),
		Carb:     pbr.GetCarb(),
	}
}

func toServiceMealItemsRequest(pbr *pb.MealItemsRequest) service.MealItemsRequest {
	return service.MealItemsRequest{
		Name:           pbr.GetName(),
		Recipe:         pbr.GetRecipe(),
		NutritionFacts: toServiceNutritionFacts(pbr.GetNutritionFacts()),
	}
}

func toServicePlannedMeals(pbr *pb.PlannedMealsRequest) service.PlannedMealsRequest {
	serviceMealItems := make([]service.MealItemsRequest, 0, len(pbr.GetMealItems()))

	for _, item := range pbr.GetMealItems() {
		serviceMealItems = append(serviceMealItems, toServiceMealItemsRequest(item))
	}

	return service.PlannedMealsRequest{
		DayOfWeek: service.DaysOfWeek(pbr.GetDayOfWeek()),
		MealItems: serviceMealItems,
	}
}

func toServiceSaveGeneratedPlanRequest(pbr *pb.SaveGeneratedPlanRequest) service.SaveGeneratedPlanRequest {
	servicePlannedMeals := make([]service.PlannedMealsRequest, 0, len(pbr.GetPlannedMeals()))

	for _, plannedMeals := range pbr.GetPlannedMeals() {
		servicePlannedMeals = append(servicePlannedMeals, toServicePlannedMeals(plannedMeals))
	}

	return service.SaveGeneratedPlanRequest{
		UserID:       pbr.GetUserId(),
		GenerationID: pbr.GetGenerationId(),
		PlannedMeals: servicePlannedMeals,
	}
}
