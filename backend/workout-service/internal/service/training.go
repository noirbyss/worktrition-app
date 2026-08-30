package service

import "strings"

type CompleteTrainingRequest struct {
	UserID          string
	DayOfWeek       DaysOfWeek
	DurationSeconds int32
}

func (r CompleteTrainingRequest) validate() error {
	if r.UserID == "" {
		return ErrEmptyUserID
	}

	if !r.DayOfWeek.validate() {
		return ErrInvalidDayOfWeek
	}

	if r.DurationSeconds <= 0 {
		return ErrInvalidDuration
	}

	return nil
}

var cardioKeywords = []string{"кардио", "cardio", "аэроб", "aerobic", "бег", "running", "run", "выносл"}

// IsStrength classifies a free-form training type into strength (true) or cardio
// (false). The AI generates the type as an arbitrary human-readable label, so the
// classification is keyword based and defaults to strength.
func IsStrength(trainingType string) bool {
	normalized := strings.ToLower(trainingType)
	for _, keyword := range cardioKeywords {
		if strings.Contains(normalized, keyword) {
			return false
		}
	}

	return true
}
