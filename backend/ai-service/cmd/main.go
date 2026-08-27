package main

import (
	"ai-service/internal/config"
	"ai-service/internal/grpcclient"
	"ai-service/internal/grpcserver"
	"ai-service/internal/logger"
	"ai-service/internal/provider"
	"ai-service/internal/usecase"
	"log"
	"net"

	"github.com/noirbyss/worktrition-app/gen/ai-service"
	"github.com/noirbyss/worktrition-app/gen/nutrition-service"
	"github.com/noirbyss/worktrition-app/gen/user-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main () {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	zapLog,closelog,err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer closelog()
	
	prov := provider.NewOpenRouterProvider(cfg.ApiKey,cfg.AiBaseURL,cfg.AiModel,zapLog)

	userConn,err := grpc.NewClient("localhost:50053",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect user-service: %v",err)
	}
	defer userConn.Close()
	userClient := grpcclient.NewUserClient(user.NewUserServiceClient(userConn),zapLog)
	
	nutrConn, err := grpc.NewClient("localhost:50051",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect nutrition-service: %v",err)
	}
	defer nutrConn.Close()
	nutrClient := grpcclient.NewNutritionClient(nutrition.NewNutritionServiceClient(nutrConn),zapLog)

	
	
	us := usecase.NewUseCase(prov,userClient,nutrClient,nil,zapLog)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	grpcSrv := grpc.NewServer()	
	ai.RegisterAiServiceServer(grpcSrv,grpcserver.NewServer(us,zapLog))

	zapLog.Infof("ai-service listening on :%s", cfg.GRPCPort)
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}