package client

import (
	"context"

	pb "github.com/noirbyss/worktrition-app/gen/gamification-service"
	"google.golang.org/grpc"
)

type WorkoutServiceClient struct {
	client pb.GamificationServiceClient
}

func New(conn grpc.ClientConnInterface) (*WorkoutServiceClient, error) {
	if conn == nil {
		return nil, ErrNilPointerConn
	}

	return &WorkoutServiceClient{
		client: pb.NewGamificationServiceClient(conn),
	}, nil
}

func (c *WorkoutServiceClient) ApplyWorkoutReward(ctx context.Context, userID string, isStrength bool) error {
	if _, err := c.client.ApplyWorkoutReward(ctx, toPBApplyWorkoutRewardRequest(userID, isStrength)); err != nil {
		return err
	}

	return nil
}
