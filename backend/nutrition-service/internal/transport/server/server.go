package server

import (
	"nutrition-service/internal/service"

	pb "github.com/noirbyss/worktrition-app/gen/nutrition-service"
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

// func (nss *NutririonServiceServer) SaveGeneratedPlan(ctx context.Context, r *pb.SaveGeneratedPlanRequest) (*emptypb.Empty, error) {

// }
