package domain

import "time"

type Gender int32
type TrainingLevel int32
type ActivityLevel int32
type FitnessGoal int32
type TrainingLocation int32

const (
	GenderUnspecified Gender = 0
	GenderMale        Gender = 1
	GenderFemale      Gender = 2
)

const (
	TrainingLevelUnspecified  TrainingLevel = 0
	TrainingLevelBeginner     TrainingLevel = 1
	TrainingLevelIntermediate TrainingLevel = 2
	TrainingLevelAdvanced     TrainingLevel = 3
)

const (
	ActivityLevelUnspecified ActivityLevel = 0
	ActivityLevelSedentary   ActivityLevel = 1
	ActivityLevelLight       ActivityLevel = 2
	ActivityLevelModerate    ActivityLevel = 3
	ActivityLevelHigh        ActivityLevel = 4
)

const (
	FitnessGoalUnspecified    FitnessGoal = 0
	FitnessGoalLoseWeight     FitnessGoal = 1
	FitnessGoalMaintainWeight FitnessGoal = 2
	FitnessGoalGainMuscle     FitnessGoal = 3
)

const (
	TrainingLocationUnspecified TrainingLocation = 0
	TrainingLocationHome        TrainingLocation = 1
	TrainingLocationGym         TrainingLocation = 2
)

type User struct {
	ID               string    `db:"id"`
	Name             string    `db:"name"`
	Email            string    `db:"email"`
	PasswordHash     string    `db:"password_hash"`
	BirthDate        time.Time `db:"birth_date"`
	ProfileCompleted bool      `db:"profile_completed"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type Profile struct {
	UserID              string           `db:"user_id"`
	Age                 int              `db:"age"`
	Gender              Gender           `db:"gender"`
	HeightCM            int              `db:"height_cm"`
	WeightKG            float64          `db:"weight_kg"`
	TrainingLevel       TrainingLevel    `db:"training_level"`
	ActivityLevel       ActivityLevel    `db:"activity_level"`
	Goal                FitnessGoal      `db:"goal"`
	TargetWeightKG      *float64         `db:"target_weight_kg"`
	Allergies           []string         `db:"allergies"`
	ExcludedFoods       []string         `db:"excluded_foods"`
	FoodPreferences     []string         `db:"food_preferences"`
	TrainingLocation    TrainingLocation `db:"training_location"`
	TrainingDaysPerWeek int              `db:"training_days_per_week"`
	Equipment           string           `db:"equipment"`
	CreatedAt           time.Time        `db:"created_at"`
	UpdatedAt           time.Time        `db:"updated_at"`
}

type RefreshToken struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	RevokedAt *time.Time `db:"revoked_at"`
	CreatedAt time.Time  `db:"created_at"`
}

type AuthSession struct {
	User                  *User
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}
