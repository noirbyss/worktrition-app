package domain

import "math"

func CalculateBMI(heightCM int, weightKG float64) float64 {
	heightM := float64(heightCM) / 100
	bmi := weightKG / (heightM * heightM)

	return math.Round(bmi*10) / 10
}

func (p Profile) BMI() float64 {
	return CalculateBMI(p.HeightCM, p.WeightKG)
}
