package usecase

import "context"

type AIProvider interface {
	GeneratePlan(ctx context.Context, systemPrompt, userPromt string) (string, error)
}