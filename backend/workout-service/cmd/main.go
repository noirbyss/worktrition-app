package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"workout-service/internal/repository"
	"workout-service/internal/service"
	"workout-service/internal/transport/client"
	"workout-service/internal/transport/server"

	pb "github.com/noirbyss/worktrition-app/gen/workout-service"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	conn, err := createDBConnection(ctx)
	if err != nil {
		log.Fatalf("failed connect to db: %v", err)
	}
	defer conn.Close()

	postgres, err := repository.New(conn)
	if err != nil {
		log.Fatalf("failed to create repo: %v", err)
	}

	svc := service.New(postgres)

	gamificationAddr := os.Getenv("GAMIFICATION_SERVICE_ADDR")
	if gamificationAddr == "" {
		log.Fatal("GAMIFICATION_SERVICE_ADDR is empty")
	}

	clientConn, err := grpc.NewClient(gamificationAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create client connection: %v", err)
	}
	defer clientConn.Close()

	serviceClient, err := client.New(clientConn)
	if err != nil {
		log.Fatalf("failed create client: %v", err)
	}

	serviceServer, err := server.New(svc, serviceClient)
	if err != nil {
		log.Fatalf("failed create server: %v", err)
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		log.Fatal("GRPC_PORT is empty")
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed start listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWorkoutServiceServer(s, serviceServer)

	log.Printf("workout-service listening on :%s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func createDBConnection(ctx context.Context) (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty")
	}

	conn, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error ping db: %w", err)
	}

	return conn, nil
}
