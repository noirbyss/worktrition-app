package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin save profile transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	err = tx.QueryRow(ctx, saveProfileQuery,
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
		if mappedErr == domain.ErrUserNotFound || domain.IsValidationError(mappedErr) {
			return mappedErr
		}

		return fmt.Errorf("save profile: %w", mappedErr)
	}

	measurement := domain.WeightMeasurement{
		MeasuredOn: normalizeDate(profile.UpdatedAt),
		UserID:     profile.UserID,
		WeightKG:   profile.WeightKG,
	}
	if err := upsertWeightMeasurement(ctx, tx, &measurement); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit save profile transaction: %w", err)
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
	if mappedErr := mapUUIDValidationError(err, "user_id"); mappedErr != nil {
		return nil, mappedErr
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

func (r *PostgresProfileRepository) SaveWeightMeasurement(
	ctx context.Context,
	userID string,
	weightKG float64,
	measuredOn time.Time,
) (*domain.WeightMeasurement, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin save weight measurement transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	commandTag, err := tx.Exec(ctx, `
		UPDATE user_profiles
		SET weight_kg = $2
		WHERE user_id = $1
	`, userID, weightKG)
	if mappedErr := mapUUIDValidationError(err, "user_id"); mappedErr != nil {
		return nil, mappedErr
	}
	if err != nil {
		return nil, fmt.Errorf("update profile weight: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return nil, domain.ErrProfileNotFound
	}

	measurement := domain.WeightMeasurement{
		UserID:     userID,
		MeasuredOn: measuredOn,
		WeightKG:   weightKG,
	}
	if err := upsertWeightMeasurement(ctx, tx, &measurement); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit save weight measurement transaction: %w", err)
	}

	return &measurement, nil
}

func (r *PostgresProfileRepository) ListWeightMeasurements(
	ctx context.Context,
	userID string,
	limit int,
) ([]domain.WeightMeasurement, error) {
	const query = `
		SELECT
			user_id::text,
			measured_on,
			weight_kg::float8,
			created_at,
			updated_at
		FROM (
			SELECT user_id, measured_on, weight_kg, created_at, updated_at
			FROM user_weight_measurements
			WHERE user_id = $1
			ORDER BY measured_on DESC, updated_at DESC
			LIMIT $2
		) measurements
		ORDER BY measured_on ASC, updated_at ASC
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if mappedErr := mapUUIDValidationError(err, "user_id"); mappedErr != nil {
		return nil, mappedErr
	}
	if err != nil {
		return nil, fmt.Errorf("list weight measurements: %w", err)
	}
	defer rows.Close()

	measurements := make([]domain.WeightMeasurement, 0, limit)
	for rows.Next() {
		var measurement domain.WeightMeasurement
		if err := rows.Scan(
			&measurement.UserID,
			&measurement.MeasuredOn,
			&measurement.WeightKG,
			&measurement.CreatedAt,
			&measurement.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan weight measurement: %w", err)
		}

		measurements = append(measurements, measurement)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate weight measurements: %w", err)
	}

	return measurements, nil
}

const saveProfileQuery = `
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

func upsertWeightMeasurement(
	ctx context.Context,
	tx pgx.Tx,
	measurement *domain.WeightMeasurement,
) error {
	if measurement == nil {
		return domain.NewValidationError("weight_measurement", "is required")
	}

	measurement.MeasuredOn = normalizeDate(measurement.MeasuredOn)

	err := tx.QueryRow(ctx, `
		INSERT INTO user_weight_measurements (user_id, measured_on, weight_kg)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, measured_on) DO UPDATE SET
			weight_kg = EXCLUDED.weight_kg
		RETURNING created_at, updated_at
	`, measurement.UserID, measurement.MeasuredOn, measurement.WeightKG).Scan(&measurement.CreatedAt, &measurement.UpdatedAt)
	if mappedErr := mapProfileWriteError(err); mappedErr != nil {
		if mappedErr == domain.ErrUserNotFound || domain.IsValidationError(mappedErr) {
			return mappedErr
		}

		return fmt.Errorf("save weight measurement: %w", mappedErr)
	}

	return nil
}

func normalizeDate(value time.Time) time.Time {
	if value.IsZero() {
		return value
	}

	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}

func stringList(values []string) []string {
	if values == nil {
		return []string{}
	}

	return values
}
