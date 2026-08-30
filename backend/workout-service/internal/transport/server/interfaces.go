package server

import "context"

type GamificationClient interface {
	ApplyWorkoutReward(ctx context.Context, userID string, isStrength bool) error
}
