package grpcclient

import (
	"ai-service/internal/domain"

	"github.com/noirbyss/worktrition-app/gen/user-service"
)

func mapProtoToUserProfile(resp *user.GetProfileResponse) *domain.UserProfile {
     return &domain.UserProfile{
         Gender:		 		resp.Gender.String(),
         HeightCm:            	resp.HeightCm,
         WeightKg:            	resp.WeightKg,
         TrainingLevel:       	resp.TrainingLevel.String(),
         ActivityLevel:       	resp.ActivityLevel.String(),
         Goal:					resp.Goal.String(),
         TargetWeightKg:      	resp.TargetWeightKg,
         Allergies:           	resp.Allergies,
         ExcludedFoods:       	resp.ExcludedFoods,
         FoodPreferences:     	resp.FoodPreferences,
         TrainingLocation:    	resp.TrainingLocation.String(),
         TrainingDaysPerWeek:	resp.TrainingDaysPerWeek,
         Equipment:           	resp.Equipment,
         Age:           		resp.Age,
     }
 }