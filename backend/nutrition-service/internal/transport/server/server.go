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

func (ns *NutritionServiceServer) SaveGeneratedPlan(ctx context.Context, planPb *pb.SaveGeneratedPlanRequest) (*emptypb.Empty, error) {
	servicePlan := toServiceCreatePlanRequest(planPb)

	if err := ns.service.SavePlan(ctx, servicePlan); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}
