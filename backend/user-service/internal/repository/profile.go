package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noirbyss/worktrition-app/backend/user-service/internal/domain"
)

type PostgresProfileRepository struct {
	pool *pgxpool.Pool
}

var _ domain.ProfileRepository = (*PostgresProfileRepository)(nil)

func NewPostgresProfileRepository(pool *pgxpool.Pool) *PostgresProfileRepository {
	return &PostgresProfileRepository{pool: pool}
}

func (r *PostgresProfileRepository) Save(ctx context.Context, profile *domain.Profile) error {
	if profile == nil {
		return domain.NewValidationError("profile", "is required")
	}

	const query = `
		INSERT INTO user_profiles (
			user_id,
			age,
			gender,
			height_cm,
			weight_kg,
			training_level,
			activity_level,
			goal,
			target_weight_kg,
			allergies,
			excluded_foods,
			food_preferences,
			training_location,
			training_days_per_week,
			equipment
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::text[], $11::text[], $12::text[], $13, $14, $15)
		ON CONFLICT (user_id) DO UPDATE SET
			age = EXCLUDED.age,
			gender = EXCLUDED.gender,
			height_cm = EXCLUDED.height_cm,
			weight_kg = EXCLUDED.weight_kg,
			training_level = EXCLUDED.training_level,
			activity_level = EXCLUDED.activity_level,
			goal = EXCLUDED.goal,
			target_weight_kg = EXCLUDED.target_weight_kg,
			allergies = EXCLUDED.allergies,
			excluded_foods = EXCLUDED.excluded_foods,
			food_preferences = EXCLUDED.food_preferences,
			training_location = EXCLUDED.training_location,
			training_days_per_week = EXCLUDED.training_days_per_week,
			equipment = EXCLUDED.equipment
		RETURNING created_at, updated_at
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		profile.UserID,
		profile.Age,
		profile.Gender,
		profile.HeightCM,
		profile.WeightKG,
		profile.TrainingLevel,
		profile.ActivityLevel,
		profile.Goal,
		profile.TargetWeightKG,
		stringList(profile.Allergies),
		stringList(profile.ExcludedFoods),
		stringList(profile.FoodPreferences),
		profile.TrainingLocation,
		profile.TrainingDaysPerWeek,
		profile.Equipment,
	).Scan(&profile.CreatedAt, &profile.UpdatedAt)
	if mappedErr := mapProfileWriteError(err); mappedErr != nil {
		if mappedErr == domain.ErrUserNotFound {
			return mappedErr
		}

		return fmt.Errorf("save profile: %w", mappedErr)
	}

	return nil
}

func (r *PostgresProfileRepository) GetByUserID(ctx context.Context, userID string) (*domain.Profile, error) {
	const query = `
		SELECT
			user_id::text,
			age,
			gender,
			height_cm,
			weight_kg::float8,
			training_level,
			activity_level,
			goal,
			target_weight_kg::float8,
			allergies,
			excluded_foods,
			food_preferences,
			training_location,
			training_days_per_week,
			equipment,
			created_at,
			updated_at
		FROM user_profiles
		WHERE user_id = $1
	`

	var profile domain.Profile
	var targetWeightKG pgtype.Float8
	var gender int16
	var trainingLevel int16
	var activityLevel int16
	var goal int16
	var trainingLocation int16

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&profile.UserID,
		&profile.Age,
		&gender,
		&profile.HeightCM,
		&profile.WeightKG,
		&trainingLevel,
		&activityLevel,
		&goal,
		&targetWeightKG,
		&profile.Allergies,
		&profile.ExcludedFoods,
		&profile.FoodPreferences,
		&trainingLocation,
		&profile.TrainingDaysPerWeek,
		&profile.Equipment,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if isNoRows(err) {
		return nil, domain.ErrProfileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get profile by user id: %w", err)
	}

	profile.Gender = domain.Gender(gender)
	profile.TrainingLevel = domain.TrainingLevel(trainingLevel)
	profile.ActivityLevel = domain.ActivityLevel(activityLevel)
	profile.Goal = domain.FitnessGoal(goal)
	profile.TrainingLocation = domain.TrainingLocation(trainingLocation)
	if targetWeightKG.Valid {
		profile.TargetWeightKG = &targetWeightKG.Float64
	}

	return &profile, nil
}

func stringList(values []string) []string {
	if values == nil {
		return []string{}
	}

	return values
}
