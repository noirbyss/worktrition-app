package domain

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
      Age		          int32
  }