package domain

type GeneratedPlanDTO struct {
	Workouts 	[]WorkoutPlanDTO		`json:"workouts"` 		
	Nutrition 	[]NutritionPlanDTO		`json:"nutrition"`
	WaterMl     int 					`json:"water_ml"` 
}

type NutritionPlanDTO struct {
	Day			 	int 			`json:"day"`
	Meals 			[]Meal			`json:"meals"`
	NutritionFacts 	NutritionFacts 	`json:"nutrition_facts"`
} 

type WorkoutPlanDTO struct {
	Day 		int 		`json:"day"`
	Type 		string		`json:"type"`
	Exercises 	[]string	`json:"exercises"`
}

type Meal struct {
	Name	 		string			`json:"meal_name"`
	Recipe 			string			`json:"recipe"`
	NutritionFacts	NutritionFacts	`json:"nutrition_facts"`  
}

type NutritionFacts struct {
	Calories 		float64 	`json:"calories"`
	Protein 		float64		`json:"protein"` 
	Fat 			float64 	`json:"fat"`
	Carbohydrates 	float64 	`json:"carbohydrates"`
}