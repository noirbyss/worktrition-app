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

	return &emptypb.Empty{}, nil // TODO: добавить валидацию
}

func (nss *NutririonServiceServer) GetDayPlan(ctx context.Context, r *pb.GetDayPlanRequest) (*pb.GetDayPlanResponse, error) {
	serviceGetDayPlanResponse, err := nss.service.GetDayPlan(ctx, toServiceGetDayPlanRequest(r))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error()) // TODO: добавить валидацию
	}

	return toPBGetDayPlanResponse(serviceGetDayPlanResponse), nil
}

func (nss *NutririonServiceServer) CompleteWater(ctx context.Context, r *pb.CompleteWaterRequest) (*emptypb.Empty, error) {
	if err := nss.service.CompleteWater(ctx, toServiceCompleteWaterRequest(r)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error()) // TODO: добавить валидацию
	}

	return &emptypb.Empty{}, nil
}
