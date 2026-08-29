package client

import (
	"context"

	pb "github.com/noirbyss/worktrition-app/gen/gamification-service"
	"google.golang.org/grpc"
)

type NutritionServiceClient struct {
	client pb.GamificationServiceClient
}

func New(conn grpc.ClientConnInterface) (*NutritionServiceClient, error) {
	if conn == nil {
		return nil, ErrNilPointerConn
	}

	return &NutritionServiceClient{
		client: pb.NewGamificationServiceClient(conn),
	}, nil
}

func (nsc *NutritionServiceClient) ApplyMealReward(ctx context.Context, userID string) error {
	if _, err := nsc.client.ApplyMealReward(ctx, toPBApplyMealRewardRequest(userID)); err != nil {
		return err
	}

	return nil
}

func (nsc *NutritionServiceClient) ApplyWaterReward(ctx context.Context, userID string) error {
	if _, err := nsc.client.ApplyWaterReward(ctx, toPBApplyWaterRewardRequest(userID)); err != nil {
		return err
	}

	return nil
}
