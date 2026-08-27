package grpcclient

import (
	"ai-service/internal/domain"

	"github.com/noirbyss/worktrition-app/gen/nutrition-service"
	"github.com/noirbyss/worktrition-app/gen/user-service"
)

func mapProtoToUserProfile(resp *user.GetProfileResponse) *domain.UserProfile {
	profile := &domain.UserProfile{
		Gender:              resp.Gender.String(),
		HeightCm:            resp.HeightCm,
		WeightKg:            resp.WeightKg,
		TrainingLevel:       resp.TrainingLevel.String(),
		ActivityLevel:       resp.ActivityLevel.String(),
		Goal:                resp.Goal.String(),
		TargetWeightKg:      resp.TargetWeightKg,
		Allergies:           resp.Allergies,
		ExcludedFoods:       resp.ExcludedFoods,
		FoodPreferences:     resp.FoodPreferences,
		TrainingLocation:    resp.TrainingLocation.String(),
		TrainingDaysPerWeek: resp.TrainingDaysPerWeek,
		Equipment:           resp.Equipment,
		Age:                 resp.Age,
		Bmi:                 resp.Bmi,
		WaterGoalMl:         resp.WaterGoalMl,
	}

	if resp.NutritionTargets != nil {
		profile.NutritionTargets = &domain.NutritionTargets{
			Calories: resp.NutritionTargets.CaloriesKcal,
			Protein:  resp.NutritionTargets.ProteinG,
			Fat:      resp.NutritionTargets.FatG,
			Carbs:    resp.NutritionTargets.CarbsG,
		}
	}

	return profile
}

func mapPlanToProto(plans []domain.NutritionPlanDTO) []*nutrition.PlannedMealsRequest{
	planReq := []*nutrition.PlannedMealsRequest{}
	for _, p := range plans {
		planReq = append(planReq, &nutrition.PlannedMealsRequest{
			DayOfWeek: nutrition.DaysOfWeek(p.Day),
			MealItems: mapMealsToProto(p.Meals),
			NutritionFacts: &nutrition.NutritionFacts{
				Calories: 	p.NutritionFacts.Calories,
				Protein: 	p.NutritionFacts.Protein,
				Fat: 		p.NutritionFacts.Fat,
				Carb: 		p.NutritionFacts.Carbohydrates,
			},
		})
	}
	return  planReq
}

func mapMealsToProto(meals []domain.Meal) []*nutrition.MealItemsRequest {
	mealItems := []*nutrition.MealItemsRequest{}
	for _, m := range meals{
		mealItems = append(mealItems, &nutrition.MealItemsRequest{
			Name: m.Name,
			Recipe: m.Recipe,
			NutritionFacts: &nutrition.NutritionFacts{
				Calories: 	m.NutritionFacts.Calories,
				Protein: 	m.NutritionFacts.Protein,
				Fat: 		m.NutritionFacts.Fat,
				Carb: 		m.NutritionFacts.Carbohydrates,
			},
		})
	}
	return  mealItems
}

func calcTotalNutritionFacts(plan []domain.NutritionPlanDTO) *nutrition.NutritionFacts {
	var cal, prot, fat, carb float64
	for _, p := range plan {
		cal  += p.NutritionFacts.Calories
		prot += p.NutritionFacts.Protein
		fat  += p.NutritionFacts.Fat
		carb += p.NutritionFacts.Carbohydrates
	}
	return &nutrition.NutritionFacts{
		Calories: 	cal,
		Protein: 	prot,
		Fat: 		fat,
		Carb: 	carb,
	}
}