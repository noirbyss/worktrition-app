package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/config"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/grpcserver"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/migrator"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/repository"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/service"
	userpb "github.com/noirbyss/worktrition-app/gen/user-service"
	"google.golang.org/grpc"
)

const grpcShutdownTimeout = 10 * time.Second

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := setupLogger(cfg); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := openPostgres(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := migrator.Run(ctx, pool); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	listener, err := net.Listen("tcp", cfg.GRPC.Address())
	if err != nil {
		return fmt.Errorf("listen on %s: %w", cfg.GRPC.Address(), err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	userRepository := repository.NewPostgresUserRepository(pool)
	profileRepository := repository.NewPostgresProfileRepository(pool)
	authService := service.NewAuthService(userRepository)
	profileService := service.NewProfileService(userRepository, profileRepository)

	userpb.RegisterUserServiceServer(grpcServer, grpcserver.New(authService, profileService))

	return serveGRPC(ctx, cfg, grpcServer, listener)
}

func openPostgres(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	connectCtx, cancel := context.WithTimeout(ctx, cfg.Database.ConnectTimeout)
	defer cancel()

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect to PostgreSQL: %w", err)
	}

	slog.Info(
		"connected to PostgreSQL",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"database", cfg.Database.Name,
	)

	return pool, nil
}

func serveGRPC(ctx context.Context, cfg *config.Config, grpcServer *grpc.Server, listener net.Listener) error {
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- grpcServer.Serve(listener)
	}()

	slog.Info(
		"user-service started",
		"environment", cfg.Environment,
		"grpc_address", cfg.GRPC.Address(),
		"log_level", cfg.Log.Level,
	)

	select {
	case err := <-serveErr:
		return wrapServeError(err)
	case <-ctx.Done():
		slog.Info("stopping user-service")
		stopGRPCServer(grpcServer)
		return wrapServeError(<-serveErr)
	}
}

func setupLogger(cfg *config.Config) error {
	level, err := parseLogLevel(cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("parse log level %q: %w", cfg.Log.Level, err)
	}

	handlerOptions := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if strings.EqualFold(cfg.Environment, "local") {
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	}

	slog.SetDefault(slog.New(handler))

	return nil
}

func parseLogLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unsupported level")
	}
}

func stopGRPCServer(grpcServer *grpc.Server) {
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	timer := time.NewTimer(grpcShutdownTimeout)
	defer timer.Stop()

	select {
	case <-stopped:
	case <-timer.C:
		grpcServer.Stop()
	}
}

func wrapServeError(err error) error {
	if err == nil || errors.Is(err, grpc.ErrServerStopped) {
		return nil
	}

	return fmt.Errorf("serve gRPC: %w", err)
}
