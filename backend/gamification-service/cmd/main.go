package main

import (
	"fmt"
	"net"
	"os"

	"gamification-service/internal/database"
	"gamification-service/internal/grpcserver"
	"gamification-service/internal/usecase"

	pb "github.com/noirbyss/worktrition-app/gen/gamification-service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50055"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/gamification_db?sslmode=disable"
	}

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	db, err := database.NewPostgres(dbURL)
	if err != nil {
		sugar.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	sugar.Info("Successfully connected to PostgreSQL and executed migrations")

	uc := usecase.NewGamificationUseCase(db, sugar)
	handler := grpcserver.NewServer(uc, sugar)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		sugar.Fatalf("Failed to listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterGamificationServiceServer(grpcSrv, handler)

	reflection.Register(grpcSrv)

	sugar.Infof("Gamification MVP service listening on port %s", port)
	if err := grpcSrv.Serve(lis); err != nil {
		sugar.Fatalf("Serve failed: %v", err)
	}
}
