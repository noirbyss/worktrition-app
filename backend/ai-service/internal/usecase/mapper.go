package usecase

import (
	"ai-service/internal/domain"
	"fmt"
	"strings"
)

func MapUserProfilePrompt (profile *domain.UserProfile) string {
	return fmt.Sprintf(
		`Данные пользователя:
	 	Пол и возраст: %s, 28 лет. 
		Параметры: Рост %d, вес %d кг.
		Активность: %s.
		Уровень подготовки: %s.
		Цель: %s.
		Желаемый вес: %.1f
		Тренировки: Готов(а) заниматься %s %d раза в неделю.
		Инвентарь: %s
		Аллергия: %s исключить их полностью.
		Терпеть не может: %s исключить их полностью
		Пищевые предпочтения: Нравится %s`,
			profile.Gender,
			profile.HeightCm,
			profile.WeightKg,
			profile.ActivityLevel,
			profile.TrainingLevel,
			profile.Goal,
			profile.TargetWeightKg,
			profile.TrainingLocation,
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
	return strings.Join(items,", ")
}