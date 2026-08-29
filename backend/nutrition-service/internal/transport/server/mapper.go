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
		DayOfWeek:      service.DaysOfWeek(pbr.GetDayOfWeek()),
		MealItems:      serviceMealItems,
		NutritionFacts: toServiceNutritionFacts(pbr.GetNutritionFacts()),
	}
}

func toServiceSaveGeneratedPlanRequest(pbr *pb.SaveGeneratedPlanRequest) service.SaveGeneratedPlanRequest {
	servicePlannedMeals := make([]service.PlannedMealsRequest, 0, len(pbr.GetPlannedMeals()))

	for _, plannedMeals := range pbr.GetPlannedMeals() {
		servicePlannedMeals = append(servicePlannedMeals, toServicePlannedMeals(plannedMeals))
	}

	return service.SaveGeneratedPlanRequest{
		UserID:         pbr.GetUserId(),
		GenerationID:   pbr.GetGenerationId(),
		PlannedMeals:   servicePlannedMeals,
		NutritionFacts: toServiceNutritionFacts(pbr.GetNutritionFacts()),
		WaterGoalMl:    pbr.GetWaterGoalMl(),
	}
}

func toServiceGetDayPlanRequest(pbr *pb.GetDayPlanRequest) service.GetDayPlanRequest {
	return service.GetDayPlanRequest{
		UserID:    pbr.GetUserId(),
		DayOfWeek: service.DaysOfWeek(pbr.GetDayOfWeek()),
	}
}

func toServiceCompleteMealRequest(pbr *pb.CompleteMealRequest) service.CompleteMealRequest {
	return service.CompleteMealRequest{
		UserID:     pbr.GetUserId(),
		MealItemID: pbr.GetMealItemId(),
	}
}

func toServiceCompleteWaterRequest(pbr *pb.CompleteWaterRequest) service.CompleteWaterRequest {
	return service.CompleteWaterRequest{
		UserID:   pbr.GetUserId(),
		AmountMl: pbr.GetAmountMl(),
	}
}

func toServiceGetStatsRequest(pbr *pb.GetStatsRequest) service.GetStatsRequest {
	return service.GetStatsRequest{
		UserID: pbr.GetUserId(),
	}
}

func toPBGetStatsResponse(sr service.GetStatsResponse) *pb.GetStatsResponse {
	return &pb.GetStatsResponse{
		PercentageComplianceNutritionFacts: sr.PercentageComplianceNutritionFacts,
		PercentagePlanFulfilled:            sr.PercentagePlanFulfilled,
		PercentageWaterStandardFulfillment: sr.PercentageWaterStandardFulfillment,
	}
}

func toPBNutritionFacts(sr service.NutritionFacts) *pb.NutritionFacts {
	return &pb.NutritionFacts{
		Calories: sr.Calories,
		Protein:  sr.Protein,
		Fat:      sr.Fat,
		Carb:     sr.Carb,
	}
}

func toPBMealItemsResponse(sr service.MealItemsResponse) *pb.MealItemsResponse {
	return &pb.MealItemsResponse{
		Id:             sr.ID,
		Name:           sr.Name,
		Recipe:         sr.Recipe,
		NutritionFacts: toPBNutritionFacts(sr.NutritionFacts),
	}
}

func toPBGetDayPlanResponse(sr service.GetDayPlanResponse) *pb.GetDayPlanResponse {
	pbMealItems := make([]*pb.MealItemsResponse, 0, len(sr.MealItems))

	for _, item := range sr.MealItems {
		pbMealItems = append(pbMealItems, toPBMealItemsResponse(item))
	}

	return &pb.GetDayPlanResponse{
		MealItems:      pbMealItems,
		NutritionFacts: toPBNutritionFacts(sr.NutritionFacts),
		WaterGoalMl:    sr.WaterGoalMl,
	}
}
