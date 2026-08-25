package server

import (
	"context"
	"nutrition-service/internal/service"

	pb "github.com/noirbyss/worktrition-app/gen/nutrition-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NutritionServiceServer struct {
	pb.UnimplementedNutritionServiceServer
	service *service.NutritionService
}

func New(service *service.NutritionService) *NutritionServiceServer {
	return &NutritionServiceServer{service: service}
}

func (ns *NutritionServiceServer) SaveGeneratedPlan(ctx context.Context, r *pb.SaveGeneratedPlanRequest) (*emptypb.Empty, error) {
	if err := ns.service.SavePlan(ctx, toServiceCreatePlanRequest(r)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (ns *NutritionServiceServer) GetDayPlan(ctx context.Context, r *pb.GetDayPlanRequest) (*pb.GetDayPlanResponse, error) {
	serviceDayPlan, err := ns.service.GetDayPlan(ctx, toServiceDayPlanRequest(r))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pbDayPlan := toPBGetDayPlanResponse(&serviceDayPlan)

	return &pbDayPlan, nil
}

func (ns *NutritionServiceServer) CompleteMeal(ctx context.Context, r *pb.CompleteMealRequest) (*emptypb.Empty, error) {
	if err := ns.service.CompleteMeal(ctx, toServiceCompleteMealRequest(r)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (ns *NutritionServiceServer) CompleteWater(ctx context.Context, r *pb.CompleteWaterRequest) (*emptypb.Empty, error) {
	if err := ns.service.CompleteWater(ctx, toServiceCompleteWaterReques(r)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}
