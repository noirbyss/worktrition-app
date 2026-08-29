package service

import "math"

type GetStatsRequest struct {
	UserID string
}

func (gsr GetStatsRequest) validate() error {
	if gsr.UserID == "" {
		return ErrEmptyUserID
	}

	return nil
}

type GetStatsResponse struct {
	PercentageComplianceNutritionFacts float64
	PercentagePlanFulfilled            float64
	PercentageWaterStandardFulfillment float64
}

const nutritionTolerance = 0.05

type NutritionDayRecord struct {
	Target   NutritionFacts
	Consumed NutritionFacts
}

func (ndr NutritionDayRecord) complianceRatio() float64 {
	calories := complianceRatio(ndr.Consumed.Calories, ndr.Target.Calories, nutritionTolerance)
	protein := complianceRatio(ndr.Consumed.Protein, ndr.Target.Protein, nutritionTolerance)
	fat := complianceRatio(ndr.Consumed.Fat, ndr.Target.Fat, nutritionTolerance)
	carb := complianceRatio(ndr.Consumed.Carb, ndr.Target.Carb, nutritionTolerance)

	return (calories + protein + fat + carb) / 4
}

func nutritionCompliancePercentage(records []NutritionDayRecord) float64 {
	if len(records) == 0 {
		return 0
	}

	var total float64

	for _, rec := range records {
		total += rec.complianceRatio()
	}

	return (total / float64(len(records))) * 100
}

type WaterDayRecord struct {
	GoalMl     int32
	ConsumedMl int32
}

func (wdr WaterDayRecord) complianceRatio() float64 {
	return complianceRatio(float64(wdr.ConsumedMl), float64(wdr.GoalMl), nutritionTolerance)
}

func waterCompliancePercentage(records []WaterDayRecord) float64 {
	if len(records) <= 0 {
		return 0
	}

	var total float64

	for _, rec := range records {
		total += rec.complianceRatio()
	}

	return (total / float64(len(records))) * 100
}

func complianceRatio(consumed, target, tolerance float64) float64 {
	if target <= 0 {
		return 0
	}

	discrepancy := math.Abs(consumed - target)
	permissibleDiscrepancy := target * tolerance
	fine := math.Max(0, discrepancy-permissibleDiscrepancy)
	procentage := 1 - fine/target

	return math.Max(0, procentage)
}

func planFulfillmentPercentage(completed, total int) float64 {
	if total <= 0 {
		return 0
	}

	return (float64(completed) / float64(total)) * 100
}
