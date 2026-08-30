package server

import (
	"context"
	"log"

	"workout-service/internal/service"
	"workout-service/internal/transport/client"

	pb "github.com/noirbyss/worktrition-app/gen/workout-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WorkoutServiceServer struct {
	pb.UnimplementedWorkoutServiceServer
	service            *service.Service
	gamificationClient GamificationClient
}

func New(service *service.Service, client *client.WorkoutServiceClient) (*WorkoutServiceServer, error) {
	if service == nil {
		return nil, ErrNilPointerService
	}

	if client == nil {
		return nil, ErrNilPointerClient
	}

	return &WorkoutServiceServer{service: service, gamificationClient: client}, nil
}

func (s *WorkoutServiceServer) SaveGeneratedPlan(ctx context.Context, r *pb.SaveGeneratedPlanRequest) (*emptypb.Empty, error) {
	if err := s.service.SavePlan(ctx, toServiceSaveGeneratedPlanRequest(r)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *WorkoutServiceServer) GetDayPlan(ctx context.Context, r *pb.GetDayPlanRequest) (*pb.GetDayPlanResponse, error) {
	dayPlan, err := s.service.GetDayPlan(ctx, toServiceGetDayPlanRequest(r))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return toPBGetDayPlanResponse(dayPlan), nil
}

func (s *WorkoutServiceServer) CompleteTraining(ctx context.Context, r *pb.CompleteTrainingRequest) (*emptypb.Empty, error) {
	trainingType, err := s.service.CompleteTraining(ctx, toServiceCompleteTrainingRequest(r))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.gamificationClient.ApplyWorkoutReward(ctx, r.GetUserId(), service.IsStrength(trainingType)); err != nil {
		// The training is already recorded; the gamification reward is best-effort.
		log.Printf("workout-service: apply workout reward for user_id=%s failed: %v", r.GetUserId(), err)
	}

	return &emptypb.Empty{}, nil
}

func (s *WorkoutServiceServer) GetStats(ctx context.Context, r *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	stats, err := s.service.GetStats(ctx, toServiceGetStatsRequest(r))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return toPBGetStatsResponse(stats), nil
}
