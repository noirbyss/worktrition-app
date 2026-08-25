package server

import (
	"nutrition-service/internal/service"

	pb "github.com/noirbyss/worktrition-app/gen/nutrition-service"
)

func toServiceNutritionFacts(pbr *pb.NutritionFacts) service.NutritionFacts {
	if pbr == nil {
		return service.NutritionFacts{}
	}

	return service.NutritionFacts{
		Calories: pbr.GetCalories(),
		Protein:  pbr.GetProtein(),
		Fat:      pbr.GetFat(),
		Carb:     pbr.GetCarbs(),
	}
}

func toServiceCreateMealRequest(pbr *pb.PlannedMeal) service.CreateMealRequest {
	return service.CreateMealRequest{
		Day:            service.DaysOfWeek(pbr.GetDayOfWeek()),
		Name:           pbr.GetName(),
		Items:          pbr.GetMealItems(),
		NutritionFacts: toServiceNutritionFacts(pbr.GetNutritionFacts()),
	}
}

func toServiceCreatePlanRequest(pbr *pb.SaveGeneratedPlanRequest) service.CreatePlanRequest {
	plannedMealPb := pbr.GetMeals()
	plannedMealService := make([]service.CreateMealRequest, 0, len(plannedMealPb))

	for _, meal := range plannedMealPb {
		plannedMealService = append(plannedMealService, toServiceCreateMealRequest(meal))
	}

	return service.CreatePlanRequest{
		UserID:         pbr.GetUserId(),
		GenerationID:   pbr.GetGenerationId(),
		Meals:          plannedMealService,
		NutritionFacts: toServiceNutritionFacts(pbr.NutritionFacts),
		WaterGoalMl:    pbr.GetWaterGoalMl(),
	}
}

func toServiceDayPlanRequest(pbr *pb.GetDayPlanRequest) service.GetDayPlanRequest {
	return service.GetDayPlanRequest{
		UserID: pbr.GetUserId(),
		Day:    service.DaysOfWeek(pbr.GetDayOfWeek()),
	}
}

func toServiceCompleteMealRequest(pbr *pb.CompleteMealRequest) service.CompleteMealRequest {
	return service.CompleteMealRequest{
		UserID: pbr.GetUserId(),
		MealID: pbr.GetMealId(),
	}
}

func toServiceCompleteWaterReques(pbr *pb.CompleteWaterRequest) service.CompleteWaterRequest {
	return service.CompleteWaterRequest{
		UserID:   pbr.GetUserId(),
		AmountMl: pbr.GetAmountMl(),
	}
}

func toPBNutritionFacts(sr *service.NutritionFacts) *pb.NutritionFacts {
	return &pb.NutritionFacts{
		Calories: sr.Calories,
		Protein:  sr.Protein,
		Fat:      sr.Fat,
		Carbs:    sr.Carb,
	}
}

func toPBDayMeal(sr *service.GetMealResponse) *pb.DayMeal {
	return &pb.DayMeal{
		Id:             sr.ID,
		Name:           sr.Name,
		MealItems:      sr.Items,
		NutritionFacts: toPBNutritionFacts(&sr.NutritionFacts),
		Completed:      sr.Completed,
	}
}

func toPBGetDayPlanResponse(sr *service.GetDayPlanResponse) pb.GetDayPlanResponse {
	pbMeals := make([]*pb.DayMeal, 0, len(sr.Meals))
	for _, serviceMeal := range sr.Meals {
		pbMeals = append(pbMeals, toPBDayMeal(&serviceMeal))
	}

	return pb.GetDayPlanResponse{
		Meals: pbMeals,
	}
}
