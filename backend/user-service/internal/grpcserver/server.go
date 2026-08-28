package grpcserver

import (
	"context"
	"errors"
	"log/slog"

	"github.com/noirbyss/worktrition-app/backend/user-service/internal/domain"
	userpb "github.com/noirbyss/worktrition-app/gen/user-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUseCase interface {
	CreateUser(ctx context.Context, name, email, plainPassword, birthDate string) (*domain.User, error)
	VerifyCredentials(ctx context.Context, email, plainPassword string) (*domain.User, error)
}

type Server struct {
	userpb.UnimplementedUserServiceServer

	auth AuthUseCase
}

func New(auth AuthUseCase) *Server {
	return &Server{auth: auth}
}

func (s *Server) CreateUser(
	ctx context.Context,
	req *userpb.CreateUserRequest,
) (*userpb.CreateUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	createdUser, err := s.auth.CreateUser(ctx, req.GetName(), req.GetEmail(), req.GetPassword(), req.GetBirthDate())
	if err != nil {
		return nil, toStatusError(err)
	}

	return &userpb.CreateUserResponse{
		UserId: createdUser.ID,
	}, nil
}

func (s *Server) VerifyCredentials(
	ctx context.Context,
	req *userpb.VerifyCredentialsRequest,
) (*userpb.VerifyCredentialsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	foundUser, err := s.auth.VerifyCredentials(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, toStatusError(err)
	}

	return &userpb.VerifyCredentialsResponse{
		UserId:           foundUser.ID,
		ProfileCompleted: foundUser.ProfileCompleted,
	}, nil
}

func toStatusError(err error) error {
	if err == nil {
		return nil
	}

	if status.Code(err) != codes.Unknown {
		return err
	}

	switch {
	case domain.IsValidationError(err):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrProfileNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		slog.Error("internal user-service error", "error", err)
		return status.Error(codes.Internal, "internal error")
	}
}
