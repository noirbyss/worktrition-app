package service

const (
	minWorkoutDays = 1
	maxWorkoutDays = 7
)

type DaysOfWeek int32

const (
	Unspecified DaysOfWeek = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (d DaysOfWeek) validate() bool {
	return d >= Monday && d <= Sunday
}

type WorkoutDayRequest struct {
	DayOfWeek DaysOfWeek
	Type      string
	Exercises []string
}

func (w WorkoutDayRequest) validate() error {
	if !w.DayOfWeek.validate() {
		return ErrInvalidDayOfWeek
	}

	if w.Type == "" {
		return ErrEmptyType
	}

	for _, exercise := range w.Exercises {
		if exercise == "" {
			return ErrEmptyExercise
		}
	}

	return nil
}

type SaveGeneratedPlanRequest struct {
	UserID       string
	GenerationID string
	WorkoutDays  []WorkoutDayRequest
}

func (r SaveGeneratedPlanRequest) validate() error {
	if r.UserID == "" {
		return ErrEmptyUserID
	}

	if r.GenerationID == "" {
		return ErrEmptyGenerationID
	}

	if len(r.WorkoutDays) < minWorkoutDays || len(r.WorkoutDays) > maxWorkoutDays {
		return ErrInvalidWorkoutDaysCount
	}

	seen := make(map[DaysOfWeek]struct{}, len(r.WorkoutDays))
	for _, day := range r.WorkoutDays {
		if err := day.validate(); err != nil {
			return err
		}

		if _, ok := seen[day.DayOfWeek]; ok {
			return ErrDuplicateDayOfWeek
		}
		seen[day.DayOfWeek] = struct{}{}
	}

	return nil
}

type GetDayPlanRequest struct {
	UserID    string
	DayOfWeek DaysOfWeek
}

func (r GetDayPlanRequest) validate() error {
	if r.UserID == "" {
		return ErrEmptyUserID
	}

	if !r.DayOfWeek.validate() {
		return ErrInvalidDayOfWeek
	}

	return nil
}

type GetDayPlanResponse struct {
	DayOfWeek DaysOfWeek
	Type      string
	Exercises []string
}
