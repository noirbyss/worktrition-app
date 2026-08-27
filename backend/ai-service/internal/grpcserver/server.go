package grpcserver

import (
	"ai-service/internal/usecase"

	"github.com/noirbyss/worktrition-app/gen/ai-service"
	"go.uber.org/zap"
)

type Server struct {
	ai.UnimplementedAiServiceServer
	us *usecase.UseCase
	logger *zap.SugaredLogger
}

func NewServer (us *usecase.UseCase,logger *zap.SugaredLogger) *Server {
	return &Server{
		us: us,
		logger: logger,
	}
}