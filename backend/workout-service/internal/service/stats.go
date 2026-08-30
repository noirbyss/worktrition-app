package service

import "time"

const dateLayout = "2006-01-02"

type GetStatsRequest struct {
	UserID string
}

func (r GetStatsRequest) validate() error {
	if r.UserID == "" {
		return ErrEmptyUserID
	}

	return nil
}

type GetStatsResponse struct {
	PercentagePlanFulfilled  float64
	CurrentStreakDays        int32
	TotalTrainingTimeSeconds int32
}

func planFulfillmentPercentage(completed, total int) float64 {
	if total <= 0 {
		return 0
	}

	return (float64(completed) / float64(total)) * 100
}

// currentStreakDays counts how many consecutive calendar days (UTC), ending today
// or yesterday, contain at least one completed training.
func currentStreakDays(completionDates []time.Time) int32 {
	if len(completionDates) == 0 {
		return 0
	}

	days := make(map[string]struct{}, len(completionDates))
	for _, date := range completionDates {
		days[date.UTC().Format(dateLayout)] = struct{}{}
	}

	cursor := time.Now().UTC()
	if _, ok := days[cursor.Format(dateLayout)]; !ok {
		cursor = cursor.AddDate(0, 0, -1)
	}

	var streak int32
	for {
		if _, ok := days[cursor.Format(dateLayout)]; !ok {
			break
		}

		streak++
		cursor = cursor.AddDate(0, 0, -1)
	}

	return streak
}
