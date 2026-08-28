package grpcserver

import (
	"context"

	"gamification-service/internal/usecase"

	pb "github.com/noirbyss/worktrition-app/gen/gamification-service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedGamificationServiceServer
	uc     *usecase.GamificationUseCase
	logger *zap.SugaredLogger
}

func NewServer(uc *usecase.GamificationUseCase, logger *zap.SugaredLogger) *Server {
	return &Server{uc: uc, logger: logger}
}

func (s *Server) GetCharacter(ctx context.Context, req *pb.GetCharacterRequest) (*pb.GetCharacterResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	c, err := s.uc.GetCharacter(ctx, req.UserId)
	if err != nil {
		s.logger.Errorf("Failed to get character %s: %v", req.UserId, err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &pb.GetCharacterResponse{Character: c}, nil
}

func (s *Server) ApplyWorkoutReward(ctx context.Context, req *pb.ApplyWorkoutRewardRequest) (*pb.ApplyWorkoutRewardResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	strGain, endGain := int32(0), int32(2)
	if req.IsStrength {
		strGain, endGain = 2, 0
	}

	res, err := s.uc.ApplyReward(ctx, req.UserId, 100, strGain, endGain, 1, 0)
	if err != nil {
		s.logger.Errorf("Failed to apply workout reward: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.ApplyWorkoutRewardResponse{
		Character: res.Character, GainedXp: res.GainedXP, LeveledUp: res.LeveledUp,
	}, nil
}

func (s *Server) ApplyMealReward(ctx context.Context, req *pb.ApplyMealRewardRequest) (*pb.ApplyMealRewardResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	res, err := s.uc.ApplyReward(ctx, req.UserId, 50, 0, 0, 0, 1)
	if err != nil {
		s.logger.Errorf("Failed to apply meal reward: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &pb.ApplyMealRewardResponse{
		Character: res.Character, GainedXp: res.GainedXP, LeveledUp: res.LeveledUp,
	}, nil
}

func (s *Server) ApplyWaterReward(ctx context.Context, req *pb.ApplyWaterRewardRequest) (*pb.ApplyWaterRewardResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	res, err := s.uc.ApplyReward(ctx, req.UserId, 20, 0, 0, 0, 0)
	if err != nil {
		s.logger.Errorf("Failed to apply water reward: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &pb.ApplyWaterRewardResponse{
		Character: res.Character, GainedXp: res.GainedXP, LeveledUp: res.LeveledUp,
	}, nil
}
