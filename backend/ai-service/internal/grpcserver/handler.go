package grpcserver

import (
	"context"
	"github.com/noirbyss/worktrition-app/gen/ai-service"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)



func (s *Server) StartGeneration(ctx context.Context, req *ai.StartGenerationRequest) (*ai.StartGenerationResponse, error) {
	userId := req.UserId
	planType := req.PlanType.Enum().String()  
	
	generation_id, err := s.us.StartGeneration(ctx, userId, planType)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start generation: %v", err)
	}
	return &ai.StartGenerationResponse{
		GenerationId: generation_id,
		Status: ai.GenerationStatus_GENERATION_STATUS_PENDING,
	}, nil
}