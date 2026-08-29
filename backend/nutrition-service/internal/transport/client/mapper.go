package client

import pb "github.com/noirbyss/worktrition-app/gen/gamification-service"

func toPBApplyMealRewardRequest(userID string) *pb.ApplyMealRewardRequest {
	return &pb.ApplyMealRewardRequest{
		UserId: userID,
	}
}

func toPBApplyWaterRewardRequest(userID string) *pb.ApplyWaterRewardRequest {
	return &pb.ApplyWaterRewardRequest{
		UserId: userID,
	}
}
