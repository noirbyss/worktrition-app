package server

import (
	"context"
	"errors"
	"log"

	"workout-service/internal/repository"
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
		return nil, toStatusError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *WorkoutServiceServer) GetDayPlan(ctx context.Context, r *pb.GetDayPlanRequest) (*pb.GetDayPlanResponse, error) {
	dayPlan, err := s.service.GetDayPlan(ctx, toServiceGetDayPlanRequest(r))
	if err != nil {
		return nil, toStatusError(err)
	}

	return toPBGetDayPlanResponse(dayPlan), nil
}

func (s *WorkoutServiceServer) CompleteTraining(ctx context.Context, r *pb.CompleteTrainingRequest) (*emptypb.Empty, error) {
	trainingType, err := s.service.CompleteTraining(ctx, toServiceCompleteTrainingRequest(r))
	if err != nil {
		return nil, toStatusError(err)
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
		return nil, toStatusError(err)
	}

	return toPBGetStatsResponse(stats), nil
}

func toStatusError(err error) error {
	switch {
	case errors.Is(err, repository.ErrPlanAlreadyExists),
		errors.Is(err, repository.ErrTrainingAlreadyCompleted):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, repository.ErrPlanNotFound),
		errors.Is(err, repository.ErrTrainingNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.InvalidArgument, err.Error())
	}
}
