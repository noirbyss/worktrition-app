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
	Register(ctx context.Context, name, email, plainPassword, birthDate string) (*domain.AuthSession, error)
	Login(ctx context.Context, email, plainPassword string) (*domain.AuthSession, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthSession, error)
	Logout(ctx context.Context, refreshToken string) error
}

type ProfileUseCase interface {
	SaveProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error)
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	GetProfile(ctx context.Context, userID string) (*domain.Profile, error)
}

type Server struct {
	userpb.UnimplementedUserServiceServer

	auth    AuthUseCase
	profile ProfileUseCase
}

func New(auth AuthUseCase, profile ProfileUseCase) *Server {
	return &Server{
		auth:    auth,
		profile: profile,
	}
}

func (s *Server) Register(
	ctx context.Context,
	req *userpb.RegisterRequest,
) (*userpb.RegisterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	session, err := s.auth.Register(ctx, req.GetName(), req.GetEmail(), req.GetPassword(), req.GetBirthDate())
	if err != nil {
		return nil, toStatusError(err)
	}

	return authSessionToRegisterProto(session), nil
}

func (s *Server) Login(
	ctx context.Context,
	req *userpb.LoginRequest,
) (*userpb.LoginResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	session, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, toStatusError(err)
	}

	return authSessionToLoginProto(session), nil
}

func (s *Server) RefreshToken(
	ctx context.Context,
	req *userpb.RefreshTokenRequest,
) (*userpb.RefreshTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	session, err := s.auth.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, toStatusError(err)
	}

	return authSessionToRefreshTokenProto(session), nil
}

func (s *Server) Logout(
	ctx context.Context,
	req *userpb.LogoutRequest,
) (*userpb.LogoutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if err := s.auth.Logout(ctx, req.GetRefreshToken()); err != nil {
		return nil, toStatusError(err)
	}

	return &userpb.LogoutResponse{
		Success: true,
	}, nil
}

func (s *Server) SaveProfile(
	ctx context.Context,
	req *userpb.SaveProfileRequest,
) (*userpb.SaveProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	savedProfile, err := s.profile.SaveProfile(ctx, profileFromProto(req))
	if err != nil {
		return nil, toStatusError(err)
	}

	return &userpb.SaveProfileResponse{
		ProfileCompleted: true,
		Bmi:              savedProfile.BMI(),
	}, nil
}

func (s *Server) GetUser(
	ctx context.Context,
	req *userpb.GetUserRequest,
) (*userpb.GetUserResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	user, err := s.profile.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, toStatusError(err)
	}

	return userToProto(user), nil
}

func (s *Server) GetProfile(
	ctx context.Context,
	req *userpb.GetProfileRequest,
) (*userpb.GetProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	profile, err := s.profile.GetProfile(ctx, req.GetUserId())
	if err != nil {
		return nil, toStatusError(err)
	}

	return profileToProto(profile), nil
}

func profileFromProto(req *userpb.SaveProfileRequest) *domain.Profile {
	var targetWeightKG *float64
	if req.TargetWeightKg != nil {
		targetWeightKG = req.TargetWeightKg
	}

	return &domain.Profile{
		UserID:              req.GetUserId(),
		Age:                 int(req.GetAge()),
		Gender:              domain.Gender(req.GetGender()),
		HeightCM:            int(req.GetHeightCm()),
		WeightKG:            req.GetWeightKg(),
		TrainingLevel:       domain.TrainingLevel(req.GetTrainingLevel()),
		ActivityLevel:       domain.ActivityLevel(req.GetActivityLevel()),
		Goal:                domain.FitnessGoal(req.GetGoal()),
		TargetWeightKG:      targetWeightKG,
		Allergies:           cloneStrings(req.GetAllergies()),
		ExcludedFoods:       cloneStrings(req.GetExcludedFoods()),
		FoodPreferences:     cloneStrings(req.GetFoodPreferences()),
		TrainingLocation:    domain.TrainingLocation(req.GetTrainingLocation()),
		TrainingDaysPerWeek: int(req.GetTrainingDaysPerWeek()),
		Equipment:           req.GetEquipment(),
	}
}

func authSessionToRegisterProto(session *domain.AuthSession) *userpb.RegisterResponse {
	if session == nil || session.User == nil {
		return &userpb.RegisterResponse{}
	}

	return &userpb.RegisterResponse{
		UserId:                session.User.ID,
		ProfileCompleted:      session.User.ProfileCompleted,
		AccessToken:           session.AccessToken,
		RefreshToken:          session.RefreshToken,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt.Unix(),
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt.Unix(),
	}
}

func authSessionToLoginProto(session *domain.AuthSession) *userpb.LoginResponse {
	if session == nil || session.User == nil {
		return &userpb.LoginResponse{}
	}

	return &userpb.LoginResponse{
		UserId:                session.User.ID,
		ProfileCompleted:      session.User.ProfileCompleted,
		AccessToken:           session.AccessToken,
		RefreshToken:          session.RefreshToken,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt.Unix(),
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt.Unix(),
	}
}

func authSessionToRefreshTokenProto(session *domain.AuthSession) *userpb.RefreshTokenResponse {
	if session == nil || session.User == nil {
		return &userpb.RefreshTokenResponse{}
	}

	return &userpb.RefreshTokenResponse{
		UserId:                session.User.ID,
		ProfileCompleted:      session.User.ProfileCompleted,
		AccessToken:           session.AccessToken,
		RefreshToken:          session.RefreshToken,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt.Unix(),
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt.Unix(),
	}
}

func userToProto(user *domain.User) *userpb.GetUserResponse {
	return &userpb.GetUserResponse{
		UserId:           user.ID,
		Name:             user.Name,
		Email:            user.Email,
		BirthDate:        user.BirthDate.Format(domain.BirthDateLayout),
		ProfileCompleted: user.ProfileCompleted,
	}
}

func profileToProto(profile *domain.Profile) *userpb.GetProfileResponse {
	var targetWeightKG *float64
	if profile.TargetWeightKG != nil {
		targetWeightKG = profile.TargetWeightKG
	}

	return &userpb.GetProfileResponse{
		UserId:              profile.UserID,
		Age:                 int32(profile.Age),
		Gender:              userpb.Gender(profile.Gender),
		HeightCm:            int32(profile.HeightCM),
		WeightKg:            profile.WeightKG,
		TrainingLevel:       userpb.TrainingLevel(profile.TrainingLevel),
		ActivityLevel:       userpb.ActivityLevel(profile.ActivityLevel),
		Goal:                userpb.FitnessGoal(profile.Goal),
		TargetWeightKg:      targetWeightKG,
		Allergies:           cloneStrings(profile.Allergies),
		ExcludedFoods:       cloneStrings(profile.ExcludedFoods),
		FoodPreferences:     cloneStrings(profile.FoodPreferences),
		TrainingLocation:    userpb.TrainingLocation(profile.TrainingLocation),
		TrainingDaysPerWeek: int32(profile.TrainingDaysPerWeek),
		Equipment:           profile.Equipment,
	}
}

func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}

	cloned := make([]string, len(values))
	copy(cloned, values)

	return cloned
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
	case errors.Is(err, domain.ErrInvalidToken),
		errors.Is(err, domain.ErrTokenExpired),
		errors.Is(err, domain.ErrTokenRevoked):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrProfileNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		slog.Error("internal user-service error", "error", err)
		return status.Error(codes.Internal, "internal error")
	}
}
