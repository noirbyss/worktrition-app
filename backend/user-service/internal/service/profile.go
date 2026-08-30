package service

import (
	"context"
	"strings"
	"time"

	"github.com/noirbyss/worktrition-app/backend/user-service/internal/domain"
)

type ProfileService struct {
	users    domain.UserRepository
	profiles domain.ProfileRepository
	now      func() time.Time
}

func NewProfileService(users domain.UserRepository, profiles domain.ProfileRepository) *ProfileService {
	return &ProfileService{
		users:    users,
		profiles: profiles,
		now:      time.Now,
	}
}

func (s *ProfileService) SaveProfile(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	if profile == nil {
		return nil, domain.NewValidationError("profile", "is required")
	}

	userID := strings.TrimSpace(profile.UserID)
	if userID == "" {
		return nil, domain.NewValidationError("user_id", "is required")
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	normalizedProfile := normalizeProfile(profile)
	normalizedProfile.UserID = userID
	normalizedProfile.Age = domain.AgeFromBirthDate(user.BirthDate, s.now())

	if err := domain.ValidateProfile(normalizedProfile); err != nil {
		return nil, err
	}
	if err := s.profiles.Save(ctx, normalizedProfile); err != nil {
		return nil, err
	}
	if err := s.users.SetProfileCompleted(ctx, normalizedProfile.UserID, true); err != nil {
		return nil, err
	}

	return normalizedProfile, nil
}

func (s *ProfileService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, domain.NewValidationError("user_id", "is required")
	}

	return s.users.GetByID(ctx, userID)
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, domain.NewValidationError("user_id", "is required")
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile, err := s.profiles.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile.Age = domain.AgeFromBirthDate(user.BirthDate, s.now())

	return profile, nil
}

func (s *ProfileService) SaveWeightMeasurement(
	ctx context.Context,
	userID string,
	weightKG float64,
) (*domain.WeightMeasurement, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, domain.NewValidationError("user_id", "is required")
	}
	if err := domain.ValidateWeightMeasurement(weightKG); err != nil {
		return nil, err
	}
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	return s.profiles.SaveWeightMeasurement(ctx, userID, weightKG, s.now())
}

func (s *ProfileService) GetWeightHistory(
	ctx context.Context,
	userID string,
	limit int,
) ([]domain.WeightMeasurement, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, domain.NewValidationError("user_id", "is required")
	}

	if limit <= 0 {
		limit = domain.DefaultWeightHistoryLimit
	}
	if err := domain.ValidateWeightHistoryLimit(limit); err != nil {
		return nil, err
	}
	if _, err := s.users.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	return s.profiles.ListWeightMeasurements(ctx, userID, limit)
}

func normalizeProfile(profile *domain.Profile) *domain.Profile {
	if profile == nil {
		return nil
	}

	normalized := *profile
	normalized.Allergies = normalizeStringList(profile.Allergies)
	normalized.ExcludedFoods = normalizeStringList(profile.ExcludedFoods)
	normalized.FoodPreferences = normalizeStringList(profile.FoodPreferences)
	normalized.Equipment = strings.TrimSpace(profile.Equipment)

	return &normalized
}

func normalizeStringList(values []string) []string {
	if values == nil {
		return nil
	}

	normalized := make([]string, len(values))
	for i, value := range values {
		normalized[i] = strings.TrimSpace(value)
	}

	return normalized
}
