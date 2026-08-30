package client

import pb "github.com/noirbyss/worktrition-app/gen/gamification-service"

func toPBApplyWorkoutRewardRequest(userID string, isStrength bool) *pb.ApplyWorkoutRewardRequest {
	return &pb.ApplyWorkoutRewardRequest{
		UserId:     userID,
		IsStrength: isStrength,
	}
}
