package service

type NutritionFacts struct {
	Calories float64
	Proteins float64
	Fats     float64
	Carbs    float64
}

func (nf *NutritionFacts) Validate() error {
	if nf.Calories <= 0 {
		return ErrInvalidCalories
	}

	if nf.Proteins <= 0 {
		return ErrInvalidProteins
	}

	if nf.Fats <= 0 {
		return ErrInvalidFats
	}

	if nf.Carbs <= 0 {
		return ErrInvalidCarbs
	}

	return nil
}

type CreateMealRequest struct {
	Day   DaysOfWeek
	Name  string
	Items []string
	NutritionFacts
}

func (cm *CreateMealRequest) Validate() error {
	if !cm.Day.isValid() {
		return ErrInvalidDayOfWeek
	}

	if cm.Name == "" {
		return ErrEmptyMealName
	}

	if len(cm.Items) == 0 {
		return ErrEmptyMealItems
	}

	if err := cm.NutritionFacts.Validate(); err != nil {
		return err
	}

	return nil
}
