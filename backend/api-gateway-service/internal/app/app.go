package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/aiapi"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/config"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/gateway"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/nutritionapi"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/userapi"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/workoutapi"
	aipb "github.com/noirbyss/worktrition-app/gen/ai-service"
	nutritionpb "github.com/noirbyss/worktrition-app/gen/nutrition-service"
	userpb "github.com/noirbyss/worktrition-app/gen/user-service"
	workoutpb "github.com/noirbyss/worktrition-app/gen/workout-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const httpShutdownTimeout = 10 * time.Second

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	userServiceConn, err := grpc.NewClient(
		cfg.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("create user-service gRPC client: %w", err)
	}
	defer userServiceConn.Close()

	nutritionServiceConn, err := grpc.NewClient(
		cfg.NutritionServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("create nutrition-service gRPC client: %w", err)
	}
	defer nutritionServiceConn.Close()

	workoutServiceConn, err := grpc.NewClient(
		cfg.WorkoutServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("create workout-service gRPC client: %w", err)
	}
	defer workoutServiceConn.Close()

	aiServiceConn, err := grpc.NewClient(
		cfg.AIServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("create ai-service gRPC client: %w", err)
	}
	defer aiServiceConn.Close()

	userHandlers := userapi.New(userpb.NewUserServiceClient(userServiceConn), userapi.Config{
		RefreshTokenCookieName: cfg.RefreshTokenCookieName,
		RequestTimeout:         cfg.UpstreamRequestTimeout,
		SecureCookies:          cfg.SecureCookies(),
	})
	nutritionHandlers := nutritionapi.New(nutritionpb.NewNutritionServiceClient(nutritionServiceConn), nutritionapi.Config{
		RequestTimeout: cfg.UpstreamRequestTimeout,
	})
	workoutHandlers := workoutapi.New(workoutpb.NewWorkoutServiceClient(workoutServiceConn), workoutapi.Config{
		RequestTimeout: cfg.UpstreamRequestTimeout,
	})
	aiHandlers := aiapi.New(aipb.NewAiServiceClient(aiServiceConn), aiapi.Config{
		RequestTimeout: cfg.UpstreamRequestTimeout,
	})

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddress(),
		Handler:           gateway.New(cfg, userHandlers, nutritionHandlers, workoutHandlers, aiHandlers),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return serveHTTP(ctx, cfg, httpServer)
}

func serveHTTP(ctx context.Context, cfg *config.Config, httpServer *http.Server) error {
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- httpServer.ListenAndServe()
	}()

	slog.Info(
		"api-gateway-service started",
		"environment", cfg.Environment,
		"http_address", cfg.HTTPAddress(),
		"upstream_request_timeout", cfg.UpstreamRequestTimeout,
		"user_service_addr", cfg.UserServiceAddr,
		"nutrition_service_addr", cfg.NutritionServiceAddr,
		"workout_service_addr", cfg.WorkoutServiceAddr,
		"ai_service_addr", cfg.AIServiceAddr,
	)

	select {
	case err := <-serveErr:
		return wrapServeError(err)
	case <-ctx.Done():
		slog.Info("stopping api-gateway-service")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		return wrapServeError(<-serveErr)
	}
}

func wrapServeError(err error) error {
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return fmt.Errorf("serve HTTP: %w", err)
}
