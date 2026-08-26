package service

type NutritionFacts struct {
	Calories float32
	Protein  float32
	Fat      float32
	Carb     float32
}

const (
	MinimalCalories float32 = 1
	MinimalProtein  float32 = 1
	MinimalFat      float32 = 1
	MinimalCarb     float32 = 1
)

func (nf NutritionFacts) validate() error {
	if nf.Calories < MinimalCalories {
		return ErrCaloriesTooLow
	}

	if nf.Protein < MinimalProtein {
		return ErrProteinTooLow
	}

	if nf.Fat < MinimalFat {
		return ErrFatTooLow
	}

	if nf.Carb < MinimalCarb {
		return ErrCarbTooLow
	}

	return nil
}

type MealItemsRequest struct {
	Name   string
	Recipe string
	NutritionFacts
}

func (mir MealItemsRequest) validate() error {
	if mir.Name == "" {
		return ErrEmptyName
	}

	if mir.Recipe == "" {
		return ErrEmptyRecipe
	}

	if err := mir.NutritionFacts.validate(); err != nil {
		return err
	}

	return nil
}

type PlannedMealsRequest struct {
	DayOfWeek DaysOfWeek
	MealItems []MealItemsRequest
	NutritionFacts
}

func (pmr PlannedMealsRequest) validate() error {
	if !pmr.DayOfWeek.validate() {
		return ErrInvalidDayOfWeek
	}

	if pmr.MealItems == nil {
		return ErrEmptyMealItems
	}

	for _, mealItem := range pmr.MealItems {
		if err := mealItem.validate(); err != nil {
			return err
		}
	}

	if err := pmr.NutritionFacts.validate(); err != nil {
		return err
	}

	return nil
}
