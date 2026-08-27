package grpcserver

import (
	"ai-service/internal/domain"
	"context"
	"errors"

	"github.com/noirbyss/worktrition-app/gen/ai-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) StartGeneration(ctx context.Context, req *ai.StartGenerationRequest) (*ai.StartGenerationResponse, error) {
	planType := req.PlanType.String()

	generationId, err := s.us.StartGeneration(ctx, req.UserId, planType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start generation: %v", err)
	}

	return &ai.StartGenerationResponse{
		GenerationId: generationId,
		Status:       ai.GenerationStatus_GENERATION_STATUS_PENDING,
	}, nil
}

func (s *Server) GetGenerationStatus(ctx context.Context, req *ai.GetGenerationStatusRequest) (*ai.GetGenerationStatusResponse, error) {
	generation, err := s.us.GetGenerationStatus(ctx, req.GenerationId)
	if err != nil {
		if errors.Is(err, domain.ErrGenerationNotFound) {
			return nil, status.Error(codes.NotFound, "generation not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get generation status: %v", err)
	}

	return &ai.GetGenerationStatusResponse{
		GenerationId: generation.ID,
		Status:       mapGenerationStatus(generation.Status),
		ErrorMessage: generation.ErrorMessage,
	}, nil
}

func mapGenerationStatus(status domain.GenerationStatus) ai.GenerationStatus {
	switch status {
	case domain.GenerationStatusPending:
		return ai.GenerationStatus_GENERATION_STATUS_PENDING
	case domain.GenerationStatusDone:
		return ai.GenerationStatus_GENERATION_STATUS_DONE
	case domain.GenerationStatusFailed:
		return ai.GenerationStatus_GENERATION_STATUS_FAILED
	default:
		return ai.GenerationStatus_GENERATION_STATUS_UNSPECIFIED
	}
}
