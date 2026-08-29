package server

import (
	"context"
	"nutrition-service/internal/service"
	"nutrition-service/internal/transport/client"

	pb "github.com/noirbyss/worktrition-app/gen/nutrition-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NutritionServiceServer struct {
	pb.UnimplementedNutritionServiceServer
	service                   *service.Service
	gamificationServiceClient GamificationClient
}

func NewServer(service *service.Service, client *client.NutritionServiceClient) (*NutritionServiceServer, error) {
	if service == nil {
		return nil, ErrNilPointerService
	}

	if client == nil {
		return nil, ErrNilPointerClient
	}

	return &NutritionServiceServer{service: service, gamificationServiceClient: client}, nil
}

func (nss *NutritionServiceServer) SaveGeneratedPlan(ctx context.Context, r *pb.SaveGeneratedPlanRequest) (*emptypb.Empty, error) {
	serviceSaveGenerationPlanRequest := toServiceSaveGeneratedPlanRequest(r)

	if err := nss.service.SavePlan(ctx, serviceSaveGenerationPlanRequest); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil // TODO: добавить валидацию
}

func (nss *NutritionServiceServer) GetDayPlan(ctx context.Context, r *pb.GetDayPlanRequest) (*pb.GetDayPlanResponse, error) {
	serviceGetDayPlanResponse, err := nss.service.GetDayPlan(ctx, toServiceGetDayPlanRequest(r))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error()) // TODO: добавить валидацию
	}

	return toPBGetDayPlanResponse(serviceGetDayPlanResponse), nil
}

func (nss *NutritionServiceServer) CompleteMeal(ctx context.Context, r *pb.CompleteMealRequest) (*emptypb.Empty, error) {
	if err := nss.service.CompleteMeal(ctx, toServiceCompleteMealRequest(r)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (nss *NutritionServiceServer) CompleteWater(ctx context.Context, r *pb.CompleteWaterRequest) (*emptypb.Empty, error) {
	if err := nss.service.CompleteWater(ctx, toServiceCompleteWaterRequest(r)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error()) // TODO: добавить валидацию
	}

	return &emptypb.Empty{}, nil
}

func (nss *NutritionServiceServer) GetStats(ctx context.Context, r *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	serviceGetStatsResponse, err := nss.service.GetStats(ctx, toServiceGetStatsRequest(r))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return toPBGetStatsResponse(serviceGetStatsResponse), nil
}
