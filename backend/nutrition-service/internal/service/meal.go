package service

type NutritionFacts struct {
	Calories float64
	Protein  float64
	Fat      float64
	Carb     float64
}

func (nf NutritionFacts) validate() error {
	if nf.Calories <= 0 {
		return ErrInvalidCalories
	}

	if nf.Protein <= 0 {
		return ErrInvalidProtein
	}

	if nf.Fat <= 0 {
		return ErrInvalidFat
	}

	if nf.Carb <= 0 {
		return ErrInvalidCarb
	}

	return nil
}

type CreateMealRequest struct {
	Day   DaysOfWeek
	Name  string
	Items []string
	NutritionFacts
}

func (cm CreateMealRequest) validate() error {
	if !cm.Day.isValid() {
		return ErrInvalidDayOfWeek
	}

	if cm.Name == "" {
		return ErrEmptyMealName
	}

	if len(cm.Items) == 0 {
		return ErrEmptyMealItems
	}

	if err := cm.NutritionFacts.validate(); err != nil {
		return err
	}

	return nil
}

type GetMealResponse struct {
	ID    int32
	Name  string
	Items []string
	NutritionFacts
	Completed bool
}

type CompleteMealRequest struct {
	UserID string
	MealID int32
}

func (cmr CompleteMealRequest) validate() error {
	if cmr.UserID == "" {
		return ErrEmptyUserID
	}

	if cmr.MealID <= 0 {
		return ErrInvalidMealID
	}

	return nil
}
