package workoutapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	authn "github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/auth"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/httpx"
	workoutpb "github.com/noirbyss/worktrition-app/gen/workout-service"
)

type successResponse struct {
	Success bool `json:"success"`
}

type dayOfWeekValue struct {
	value workoutpb.DaysOfWeek
}

type workoutDayRequest struct {
	DayOfWeek dayOfWeekValue `json:"day_of_week"`
	Type      string         `json:"type"`
	Exercises []string       `json:"exercises"`
}

type saveGeneratedPlanRequest struct {
	GenerationID string              `json:"generation_id"`
	WorkoutDays  []workoutDayRequest `json:"workout_days"`
}

type completeTrainingRequest struct {
	DayOfWeek       dayOfWeekValue `json:"day_of_week"`
	DurationSeconds int32          `json:"duration_seconds"`
}

func (h *Handler) HandlePlan(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetDayPlan(w, r)
	case http.MethodPost:
		h.handleSaveGeneratedPlan(w, r)
	default:
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) HandleCompleteTraining(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req completeTrainingRequest
	if !httpx.DecodeJSONBody(w, r, &req) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	if _, err := h.workoutClient.CompleteTraining(ctx, &workoutpb.CompleteTrainingRequest{
		UserId:          authn.UserIDFromContext(r.Context()),
		DayOfWeek:       req.DayOfWeek.value,
		DurationSeconds: req.DurationSeconds,
	}); err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, successResponse{Success: true})
}

func (h *Handler) HandleGetStats(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.workoutClient.GetStats(ctx, &workoutpb.GetStatsRequest{
		UserId: authn.UserIDFromContext(r.Context()),
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleGetDayPlan(w http.ResponseWriter, r *http.Request) {
	dayOfWeek, err := parseDayOfWeek(r.URL.Query().Get("day_of_week"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.workoutClient.GetDayPlan(ctx, &workoutpb.GetDayPlanRequest{
		UserId:    authn.UserIDFromContext(r.Context()),
		DayOfWeek: dayOfWeek,
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleSaveGeneratedPlan(w http.ResponseWriter, r *http.Request) {
	var req saveGeneratedPlanRequest
	if !httpx.DecodeJSONBody(w, r, &req) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	if _, err := h.workoutClient.SaveGeneratedPlan(ctx, &workoutpb.SaveGeneratedPlanRequest{
		UserId:       authn.UserIDFromContext(r.Context()),
		GenerationId: req.GenerationID,
		WorkoutDays:  toProtoWorkoutDays(req.WorkoutDays),
	}); err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, successResponse{Success: true})
}

func toProtoWorkoutDays(workoutDays []workoutDayRequest) []*workoutpb.WorkoutDayRequest {
	result := make([]*workoutpb.WorkoutDayRequest, 0, len(workoutDays))

	for _, workoutDay := range workoutDays {
		result = append(result, &workoutpb.WorkoutDayRequest{
			DayOfWeek: workoutDay.DayOfWeek.value,
			Type:      workoutDay.Type,
			Exercises: workoutDay.Exercises,
		})
	}

	return result
}

func (d *dayOfWeekValue) UnmarshalJSON(data []byte) error {
	var textValue string
	if err := json.Unmarshal(data, &textValue); err == nil {
		dayOfWeek, parseErr := parseDayOfWeek(textValue)
		if parseErr != nil {
			return parseErr
		}

		d.value = dayOfWeek
		return nil
	}

	var numericValue int32
	if err := json.Unmarshal(data, &numericValue); err == nil {
		dayOfWeek, parseErr := parseDayOfWeek(strconv.FormatInt(int64(numericValue), 10))
		if parseErr != nil {
			return parseErr
		}

		d.value = dayOfWeek
		return nil
	}

	return fmt.Errorf("day_of_week must be a string or integer from 1 to 7")
}

func parseDayOfWeek(raw string) (workoutpb.DaysOfWeek, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if normalized == "" {
		return workoutpb.DaysOfWeek_DAYS_OF_WEEK_UNSPECIFIED, fmt.Errorf("day_of_week is required")
	}

	if numericValue, err := strconv.Atoi(normalized); err == nil {
		if numericValue < int(workoutpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY) || numericValue > int(workoutpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY) {
			return workoutpb.DaysOfWeek_DAYS_OF_WEEK_UNSPECIFIED, fmt.Errorf("day_of_week must be between 1 and 7")
		}

		return workoutpb.DaysOfWeek(numericValue), nil
	}

	if dayOfWeek, ok := supportedDaysOfWeek[normalized]; ok {
		return dayOfWeek, nil
	}

	return workoutpb.DaysOfWeek_DAYS_OF_WEEK_UNSPECIFIED, fmt.Errorf("unsupported day_of_week %q", raw)
}

var supportedDaysOfWeek = map[string]workoutpb.DaysOfWeek{
	"1":                      workoutpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"2":                      workoutpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"3":                      workoutpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"4":                      workoutpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
	"5":                      workoutpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"6":                      workoutpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"7":                      workoutpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"days_of_week_friday":    workoutpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"days_of_week_monday":    workoutpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"days_of_week_saturday":  workoutpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"days_of_week_sunday":    workoutpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"days_of_week_thursday":  workoutpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
	"days_of_week_tuesday":   workoutpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"days_of_week_wednesday": workoutpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"fri":                    workoutpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"friday":                 workoutpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"mon":                    workoutpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"monday":                 workoutpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"sat":                    workoutpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"saturday":               workoutpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"sun":                    workoutpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"sunday":                 workoutpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"thu":                    workoutpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
	"thursday":               workoutpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
	"tue":                    workoutpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"tuesday":                workoutpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"wed":                    workoutpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"wednesday":              workoutpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"воскресенье":            workoutpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"вс":                     workoutpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"пн":                     workoutpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"понедельник":            workoutpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"пт":                     workoutpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"пятница":                workoutpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"сб":                     workoutpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"среда":                  workoutpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"ср":                     workoutpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"суббота":                workoutpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"вт":                     workoutpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"вторник":                workoutpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"четверг":                workoutpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
	"чт":                     workoutpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
}
