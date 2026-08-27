package app

import (
	"ai-service/internal/config"
	"ai-service/internal/database"
	"ai-service/internal/grpcclient"
	"ai-service/internal/grpcserver"
	"ai-service/internal/provider"
	"ai-service/internal/repository"
	"ai-service/internal/usecase"
	"fmt"
	"net"

	"github.com/noirbyss/worktrition-app/gen/ai-service"
	"github.com/noirbyss/worktrition-app/gen/nutrition-service"
	"github.com/noirbyss/worktrition-app/gen/user-service"
	"github.com/noirbyss/worktrition-app/gen/workout-service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func New(cfg *config.Config, logger *zap.SugaredLogger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run() error {
	if a.cfg.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := database.NewPostgres(a.cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	prov := provider.NewOpenRouterProvider(a.cfg.ApiKey, a.cfg.AiBaseURL, a.cfg.AiModel, a.logger)

	userConn, err := grpc.NewClient(a.cfg.UserServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("connect user-service: %w", err)
	}
	defer userConn.Close()
	userClient := grpcclient.NewUserClient(user.NewUserServiceClient(userConn), a.logger)

	nutrConn, err := grpc.NewClient(a.cfg.NutritionServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("connect nutrition-service: %w", err)
	}
	defer nutrConn.Close()
	nutrClient := grpcclient.NewNutritionClient(nutrition.NewNutritionServiceClient(nutrConn), a.logger)

	workoutConn, err := grpc.NewClient(a.cfg.WorkoutServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("connect workout-service: %w", err)
	}
	defer workoutConn.Close()
	workoutClient := grpcclient.NewWorkoutClient(workout.NewWorkoutServiceClient(workoutConn), a.logger)

	generationRepo := repository.NewPostgresRepository(db)
	uc := usecase.NewUseCase(prov, userClient, nutrClient, workoutClient, generationRepo, a.logger)

	lis, err := net.Listen("tcp", ":"+a.cfg.GRPCPort)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	grpcSrv := grpc.NewServer()
	ai.RegisterAiServiceServer(grpcSrv, grpcserver.NewServer(uc, a.logger))

	a.logger.Infof("ai-service listening on :%s", a.cfg.GRPCPort)
	if err := grpcSrv.Serve(lis); err != nil {
		return fmt.Errorf("serve grpc: %w", err)
	}

	return nil
}
