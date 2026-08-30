package nutritionapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	authn "github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/auth"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/httpx"
	nutritionpb "github.com/noirbyss/worktrition-app/gen/nutrition-service"
)

type successResponse struct {
	Success bool `json:"success"`
}

type dayOfWeekValue struct {
	value nutritionpb.DaysOfWeek
}

type nutritionFactsRequest struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carb     float64 `json:"carb"`
}

type mealItemRequest struct {
	Name           string                `json:"name"`
	Recipe         string                `json:"recipe"`
	NutritionFacts nutritionFactsRequest `json:"nutrition_facts"`
}

type plannedMealRequest struct {
	DayOfWeek      dayOfWeekValue        `json:"day_of_week"`
	MealItems      []mealItemRequest     `json:"meal_items"`
	NutritionFacts nutritionFactsRequest `json:"nutrition_facts"`
}

type saveGeneratedPlanRequest struct {
	GenerationID   string                `json:"generation_id"`
	PlannedMeals   []plannedMealRequest  `json:"planned_meals"`
	NutritionFacts nutritionFactsRequest `json:"nutrition_facts"`
	WaterGoalMl    int32                 `json:"water_goal_ml"`
}

type completeMealRequest struct {
	MealItemID int32 `json:"meal_item_id"`
}

type completeWaterRequest struct {
	AmountMl int32 `json:"amount_ml"`
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

func (h *Handler) HandleCompleteMeal(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req completeMealRequest
	if !httpx.DecodeJSONBody(w, r, &req) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	if _, err := h.nutritionClient.CompleteMeal(ctx, &nutritionpb.CompleteMealRequest{
		UserId:     authn.UserIDFromContext(r.Context()),
		MealItemId: req.MealItemID,
	}); err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, successResponse{Success: true})
}

func (h *Handler) HandleCompleteWater(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req completeWaterRequest
	if !httpx.DecodeJSONBody(w, r, &req) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	if _, err := h.nutritionClient.CompleteWater(ctx, &nutritionpb.CompleteWaterRequest{
		UserId:   authn.UserIDFromContext(r.Context()),
		AmountMl: req.AmountMl,
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

	resp, err := h.nutritionClient.GetStats(ctx, &nutritionpb.GetStatsRequest{
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

	resp, err := h.nutritionClient.GetDayPlan(ctx, &nutritionpb.GetDayPlanRequest{
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

	if _, err := h.nutritionClient.SaveGeneratedPlan(ctx, &nutritionpb.SaveGeneratedPlanRequest{
		UserId:         authn.UserIDFromContext(r.Context()),
		GenerationId:   req.GenerationID,
		PlannedMeals:   toProtoPlannedMeals(req.PlannedMeals),
		NutritionFacts: req.NutritionFacts.toProto(),
		WaterGoalMl:    req.WaterGoalMl,
	}); err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, successResponse{Success: true})
}

func (r nutritionFactsRequest) toProto() *nutritionpb.NutritionFacts {
	return &nutritionpb.NutritionFacts{
		Calories: r.Calories,
		Protein:  r.Protein,
		Fat:      r.Fat,
		Carb:     r.Carb,
	}
}

func toProtoPlannedMeals(plannedMeals []plannedMealRequest) []*nutritionpb.PlannedMealsRequest {
	result := make([]*nutritionpb.PlannedMealsRequest, 0, len(plannedMeals))

	for _, plannedMeal := range plannedMeals {
		result = append(result, &nutritionpb.PlannedMealsRequest{
			DayOfWeek:      plannedMeal.DayOfWeek.value,
			MealItems:      toProtoMealItems(plannedMeal.MealItems),
			NutritionFacts: plannedMeal.NutritionFacts.toProto(),
		})
	}

	return result
}

func toProtoMealItems(mealItems []mealItemRequest) []*nutritionpb.MealItemsRequest {
	result := make([]*nutritionpb.MealItemsRequest, 0, len(mealItems))

	for _, mealItem := range mealItems {
		result = append(result, &nutritionpb.MealItemsRequest{
			Name:           mealItem.Name,
			Recipe:         mealItem.Recipe,
			NutritionFacts: mealItem.NutritionFacts.toProto(),
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

func parseDayOfWeek(raw string) (nutritionpb.DaysOfWeek, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if normalized == "" {
		return nutritionpb.DaysOfWeek_DAYS_OF_WEEK_UNSPECIFIED, fmt.Errorf("day_of_week is required")
	}

	if numericValue, err := strconv.Atoi(normalized); err == nil {
		if numericValue < int(nutritionpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY) || numericValue > int(nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY) {
			return nutritionpb.DaysOfWeek_DAYS_OF_WEEK_UNSPECIFIED, fmt.Errorf("day_of_week must be between 1 and 7")
		}

		return nutritionpb.DaysOfWeek(numericValue), nil
	}

	if dayOfWeek, ok := supportedDaysOfWeek[normalized]; ok {
		return dayOfWeek, nil
	}

	return nutritionpb.DaysOfWeek_DAYS_OF_WEEK_UNSPECIFIED, fmt.Errorf("unsupported day_of_week %q", raw)
}

var supportedDaysOfWeek = map[string]nutritionpb.DaysOfWeek{
	"1":                      nutritionpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"2":                      nutritionpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"3":                      nutritionpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"4":                      nutritionpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
	"5":                      nutritionpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"6":                      nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"7":                      nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"days_of_week_friday":    nutritionpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"days_of_week_monday":    nutritionpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"days_of_week_saturday":  nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"days_of_week_sunday":    nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"days_of_week_thursday":  nutritionpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
	"days_of_week_tuesday":   nutritionpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"days_of_week_wednesday": nutritionpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"fri":                    nutritionpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"friday":                 nutritionpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"mon":                    nutritionpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"monday":                 nutritionpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"sat":                    nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"saturday":               nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"sun":                    nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"sunday":                 nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"thu":                    nutritionpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
	"thursday":               nutritionpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
	"tue":                    nutritionpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"tuesday":                nutritionpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"wed":                    nutritionpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"wednesday":              nutritionpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"воскресенье":            nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"вс":                     nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SUNDAY,
	"пн":                     nutritionpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"понедельник":            nutritionpb.DaysOfWeek_DAYS_OF_WEEK_MONDAY,
	"пт":                     nutritionpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"пятница":                nutritionpb.DaysOfWeek_DAYS_OF_WEEK_FRIDAY,
	"сб":                     nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"среда":                  nutritionpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"ср":                     nutritionpb.DaysOfWeek_DAYS_OF_WEEK_WEDNESDAY,
	"суббота":                nutritionpb.DaysOfWeek_DAYS_OF_WEEK_SATURDAY,
	"вт":                     nutritionpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"вторник":                nutritionpb.DaysOfWeek_DAYS_OF_WEEK_TUESDAY,
	"четверг":                nutritionpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
	"чт":                     nutritionpb.DaysOfWeek_DAYS_OF_WEEK_THURSDAY,
}
