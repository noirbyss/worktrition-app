package grpc

import (
	"ai-service/internal/usecase"

	"github.com/noirbyss/worktrition-app/gen/ai-service"
)

type Server struct {
	ai.UnimplementedAiServiceServer
	us *usecase.UseCase
}

func NewServer (us *usecase.UseCase) *Server {
	return &Server{us: us}
}