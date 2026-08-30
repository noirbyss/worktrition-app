package service

import (
	"testing"
	"time"
)

func TestCurrentStreakDays(t *testing.T) {
	day := func(offset int) time.Time {
		return time.Now().UTC().AddDate(0, 0, offset)
	}

	tests := []struct {
		name  string
		dates []time.Time
		want  int32
	}{
		{name: "empty", dates: nil, want: 0},
		{name: "today only", dates: []time.Time{day(0)}, want: 1},
		{name: "yesterday only", dates: []time.Time{day(-1)}, want: 1},
		{name: "three in a row from today", dates: []time.Time{day(0), day(-1), day(-2)}, want: 3},
		{name: "gap breaks streak", dates: []time.Time{day(0), day(-2), day(-3)}, want: 1},
		{name: "stale streak", dates: []time.Time{day(-3), day(-4)}, want: 0},
		{name: "duplicates ignored", dates: []time.Time{day(0), day(0), day(-1)}, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := currentStreakDays(tt.dates); got != tt.want {
				t.Fatalf("currentStreakDays() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestIsStrength(t *testing.T) {
	tests := map[string]bool{
		"Силовая":            true,
		"силовая тренировка": true,
		"Кардио":             false,
		"CARDIO session":     false,
		"Аэробная":           false,
		"Отдых":              true,
		"":                   true,
	}

	for input, want := range tests {
		if got := IsStrength(input); got != want {
			t.Errorf("IsStrength(%q) = %v, want %v", input, got, want)
		}
	}
}
