package server

import (
	"context"
	"nutrition-service/internal/service"

	pb "github.com/noirbyss/worktrition-app/gen/nutrition-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NutririonServiceServer struct {
	pb.UnimplementedNutritionServiceServer
	service *service.Service
}

func NewServer(service *service.Service) (*NutririonServiceServer, error) {
	if service == nil {
		return nil, ErrNilPointerService
	}

	return &NutririonServiceServer{service: service}, nil
}

func (nss *NutririonServiceServer) SaveGeneratedPlan(ctx context.Context, r *pb.SaveGeneratedPlanRequest) (*emptypb.Empty, error) {
	serviceSaveGenerationPlanRequest := toServiceSaveGeneratedPlanRequest(r)

	if err := nss.service.SavePlan(ctx, serviceSaveGenerationPlanRequest); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}
