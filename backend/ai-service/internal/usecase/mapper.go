package usecase

import (
	"ai-service/internal/domain"
	"fmt"
	"strings"
)

func MapUserProfilePrompt(profile *domain.UserProfile) string {
	targetWeight := "не указан"
	if profile.TargetWeightKg != nil {
		targetWeight = fmt.Sprintf("%.1f кг", *profile.TargetWeightKg)
	}

	nutritionTargets := "не рассчитаны"
	if profile.NutritionTargets != nil {
		nutritionTargets = fmt.Sprintf(
			"калории %.0f ккал, белки %.0f г, жиры %.0f г, углеводы %.0f г",
			profile.NutritionTargets.Calories,
			profile.NutritionTargets.Protein,
			profile.NutritionTargets.Fat,
			profile.NutritionTargets.Carbs,
		)
	}

	waterGoal := "не рассчитана"
	if profile.WaterGoalMl > 0 {
		waterGoal = fmt.Sprintf("%d мл", profile.WaterGoalMl)
	}

	bmi := "не рассчитан"
	if profile.Bmi > 0 {
		bmi = fmt.Sprintf("%.1f", profile.Bmi)
	}

	return fmt.Sprintf(
		`Данные пользователя:
		Пол и возраст: %s, %d лет.
		Параметры: Рост %d см, вес %.1f кг, ИМТ %s.
		Активность: %s.
		Уровень подготовки: %s.
		Цель: %s.
		Желаемый вес: %s.
		Рассчитанные цели питания: %s.
		Норма воды: %s.
		Тренировки: Готов(а) заниматься %s %d раза в неделю.
		Инвентарь: %s.
		Аллергия: %s — исключить полностью.
		Терпеть не может: %s — исключить полностью.
		Пищевые предпочтения: Нравится %s.`,
		mapGender(profile.Gender),
		profile.Age,
		profile.HeightCm,
		profile.WeightKg,
		bmi,
		mapActivityLevel(profile.ActivityLevel),
		mapTrainingLevel(profile.TrainingLevel),
		mapFitnessGoal(profile.Goal),
		targetWeight,
		nutritionTargets,
		waterGoal,
		mapTrainingLocation(profile.TrainingLocation),
		profile.TrainingDaysPerWeek,
		profile.Equipment,
		JoinOrDefault(profile.Allergies, "нет"),
		JoinOrDefault(profile.ExcludedFoods, "нет"),
		JoinOrDefault(profile.FoodPreferences, "нет"),
	)
}

func JoinOrDefault(items []string, defaultVal string) string {
	if len(items) == 0 {
		return defaultVal
	}
	return strings.Join(items, ", ")
}

func mapGender(gender string) string {
	switch gender {
	case "GENDER_MALE":
		return "Мужчина"
	case "GENDER_FEMALE":
		return "Женщина"
	default:
		return gender
	}
}

func mapActivityLevel(level string) string {
	switch level {
	case "ACTIVITY_LEVEL_SEDENTARY":
		return "Работа в офисе (почти не двигаюсь)"
	case "ACTIVITY_LEVEL_LIGHT":
		return "Лёгкая активность"
	case "ACTIVITY_LEVEL_MODERATE":
		return "Средняя активность"
	case "ACTIVITY_LEVEL_HIGH":
		return "Высокая активность"
	default:
		return level
	}
}

func mapTrainingLevel(level string) string {
	switch level {
	case "TRAINING_LEVEL_BEGINNER":
		return "начинающий"
	case "TRAINING_LEVEL_INTERMEDIATE":
		return "средний"
	case "TRAINING_LEVEL_ADVANCED":
		return "продвинутый"
	default:
		return level
	}
}

func mapFitnessGoal(goal string) string {
	switch goal {
	case "FITNESS_GOAL_LOSE_WEIGHT":
		return "Похудение"
	case "FITNESS_GOAL_MAINTAIN_WEIGHT":
		return "Поддержание веса"
	case "FITNESS_GOAL_GAIN_MUSCLE":
		return "Наращивание мышечной массы"
	default:
		return goal
	}
}

func mapTrainingLocation(location string) string {
	switch location {
	case "TRAINING_LOCATION_HOME":
		return "дома"
	case "TRAINING_LOCATION_GYM":
		return "в зале"
	default:
		return location
	}
}
