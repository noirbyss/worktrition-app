package domain

type GenerationStatus string

const (
	GenerationStatusPending GenerationStatus = "PENDING"
	GenerationStatusDone    GenerationStatus = "DONE"
	GenerationStatusFailed  GenerationStatus = "FAILED"
)

type Generation struct {
	ID           string
	UserID       string
	PlanType     string
	Status       GenerationStatus
	ErrorMessage string
}

func ShouldSaveNutrition(planType string) bool {
	return planType == "PLAN_TYPE_ALL" || planType == "PLAN_TYPE_NUTRITION" || planType == "PLAN_TYPE_UNSPECIFIED"
}

func ShouldSaveWorkout(planType string) bool {
	return planType == "PLAN_TYPE_ALL" || planType == "PLAN_TYPE_WORKOUT" || planType == "PLAN_TYPE_UNSPECIFIED"
}

type UserProfile struct {
	Gender              string
	HeightCm            int32
	WeightKg            float64
	TrainingLevel       string
	ActivityLevel       string
	Goal                string
	TargetWeightKg      *float64
	Allergies           []string
	ExcludedFoods       []string
	FoodPreferences     []string
	TrainingLocation    string
	TrainingDaysPerWeek int32
	Equipment           string
	Age                 int32
	Bmi                 float64
	NutritionTargets    *NutritionTargets
	WaterGoalMl         int32
}

type NutritionTargets struct {
	Calories float64
	Protein  float64
	Fat      float64
	Carbs    float64
}