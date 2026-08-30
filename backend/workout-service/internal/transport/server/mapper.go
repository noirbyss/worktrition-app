package server

import (
	"workout-service/internal/service"

	pb "github.com/noirbyss/worktrition-app/gen/workout-service"
)

func toServiceWorkoutDay(pbr *pb.WorkoutDayRequest) service.WorkoutDayRequest {
	return service.WorkoutDayRequest{
		DayOfWeek: service.DaysOfWeek(pbr.GetDayOfWeek()),
		Type:      pbr.GetType(),
		Exercises: pbr.GetExercises(),
	}
}

func toServiceSaveGeneratedPlanRequest(pbr *pb.SaveGeneratedPlanRequest) service.SaveGeneratedPlanRequest {
	workoutDays := make([]service.WorkoutDayRequest, 0, len(pbr.GetWorkoutDays()))
	for _, day := range pbr.GetWorkoutDays() {
		workoutDays = append(workoutDays, toServiceWorkoutDay(day))
	}

	return service.SaveGeneratedPlanRequest{
		UserID:       pbr.GetUserId(),
		GenerationID: pbr.GetGenerationId(),
		WorkoutDays:  workoutDays,
	}
}

func toServiceGetDayPlanRequest(pbr *pb.GetDayPlanRequest) service.GetDayPlanRequest {
	return service.GetDayPlanRequest{
		UserID:    pbr.GetUserId(),
		DayOfWeek: service.DaysOfWeek(pbr.GetDayOfWeek()),
	}
}

func toServiceCompleteTrainingRequest(pbr *pb.CompleteTrainingRequest) service.CompleteTrainingRequest {
	return service.CompleteTrainingRequest{
		UserID:          pbr.GetUserId(),
		DayOfWeek:       service.DaysOfWeek(pbr.GetDayOfWeek()),
		DurationSeconds: pbr.GetDurationSeconds(),
	}
}

func toServiceGetStatsRequest(pbr *pb.GetStatsRequest) service.GetStatsRequest {
	return service.GetStatsRequest{
		UserID: pbr.GetUserId(),
	}
}

func toPBGetDayPlanResponse(sr service.GetDayPlanResponse) *pb.GetDayPlanResponse {
	return &pb.GetDayPlanResponse{
		DayOfWeek:   pb.DaysOfWeek(sr.DayOfWeek),
		Type:        sr.Type,
		Exercises:   sr.Exercises,
		IsCompleted: sr.IsCompleted,
	}
}

func toPBGetStatsResponse(sr service.GetStatsResponse) *pb.GetStatsResponse {
	return &pb.GetStatsResponse{
		PercentagePlanFulfilled:  sr.PercentagePlanFulfilled,
		CurrentStreakDays:        sr.CurrentStreakDays,
		TotalTrainingTimeSeconds: sr.TotalTrainingTimeSeconds,
	}
}
