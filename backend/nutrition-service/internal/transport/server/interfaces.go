package server

import (
	"context"
)

type GamificationClient interface {
	ApplyMealReward(ctx context.Context, userID string) error
	ApplyWaterReward(ctx context.Context, userID string) error
}
