package domain

type GeneratedPlanDTO struct {
	Workouts 	[]WorkoutPlanDTO		`json:"workouts"` 		
	Nutrition 	[]NutriationPlanDTO		`json:"nutrition"`
}

type NutriationPlanDTO struct {
	Day 	int 	`json:"day"`
	Meals 	[]Meal	`json:"meals"`
} 

type WorkoutPlanDTO struct {
	Day 		int 		`json:"day"`
	Type 		string		`json:"type"`
	Exercises 	[]string	`json:"exercises"`
}

type Meal struct {
	MealName 		string	`json:"meal_name"`
	Description 	string	`json:"description"`
	Calories 		int 	`json:"calories"`
	Protein 		int		`json:"protein"` 
	Fat 			int 	`json:"fat"`
	Carbohydrates 	int 	`json:"carbohydrates"`
}

