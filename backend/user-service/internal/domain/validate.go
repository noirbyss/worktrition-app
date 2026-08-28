package domain

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	BirthDateLayout = "2006-01-02"

	MinNameLen     = 2
	MaxNameLen     = 100
	MinPasswordLen = 8
	MaxPasswordLen = 72
	MaxEmailLen    = 254

	MinAge      = 13
	MaxAge      = 120
	MinHeightCM = 80
	MaxHeightCM = 250
	MinWeightKG = 25
	MaxWeightKG = 400

	MaxStringListItems      = 50
	MaxStringListItemLength = 80
	MaxEquipmentLength      = 500
)

const (
	nameField                = "name"
	emailField               = "email"
	passwordField            = "password"
	birthDateField           = "birth_date"
	ageField                 = "age"
	genderField              = "gender"
	heightField              = "height_cm"
	weightField              = "weight_kg"
	trainingLevelField       = "training_level"
	activityLevelField       = "activity_level"
	goalField                = "goal"
	targetWeightField        = "target_weight_kg"
	trainingLocationField    = "training_location"
	trainingDaysPerWeekField = "training_days_per_week"
	equipmentField           = "equipment"

	msgRequired = "is required"
	msgTooLong  = "is too long"
)

func ValidateCreateUser(name, email, password, birthDate string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	if _, err := ParseBirthDate(birthDate); err != nil {
		return err
	}

	return nil
}

func ValidateCredentials(email, password string) error {
	if strings.TrimSpace(email) == "" {
		return NewValidationError(emailField, msgRequired)
	}
	if password == "" {
		return NewValidationError(passwordField, msgRequired)
	}

	return nil
}

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	length := utf8.RuneCountInString(name)

	if length == 0 {
		return NewValidationError(nameField, msgRequired)
	}
	if length < MinNameLen {
		return NewValidationError(nameField, fmt.Sprintf("must be at least %d characters long", MinNameLen))
	}
	if length > MaxNameLen {
		return NewValidationError(nameField, msgTooLong)
	}

	return nil
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return NewValidationError(emailField, msgRequired)
	}
	if utf8.RuneCountInString(email) > MaxEmailLen {
		return NewValidationError(emailField, msgTooLong)
	}
	if strings.ContainsAny(email, " \t\r\n") {
		return NewValidationError(emailField, "is not a valid email address")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return NewValidationError(emailField, "is not a valid email address")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return NewValidationError(passwordField, msgRequired)
	}
	if len(password) < MinPasswordLen {
		return NewValidationError(
			passwordField,
			fmt.Sprintf("must be at least %d characters long", MinPasswordLen),
		)
	}
	if len(password) > MaxPasswordLen {
		return NewValidationError(passwordField, msgTooLong)
	}
	if !hasLetterAndDigit(password) {
		return NewValidationError(passwordField, "must contain at least one letter and one digit")
	}

	return nil
}

func ParseBirthDate(value string) (time.Time, error) {
	birthDate, err := time.Parse(BirthDateLayout, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, NewValidationError(birthDateField, "must use YYYY-MM-DD format")
	}
	if err := ValidateBirthDate(birthDate); err != nil {
		return time.Time{}, err
	}

	return birthDate, nil
}

func ValidateBirthDate(birthDate time.Time) error {
	if birthDate.IsZero() {
		return NewValidationError(birthDateField, msgRequired)
	}
	if birthDate.After(time.Now()) {
		return NewValidationError(birthDateField, "cannot be in the future")
	}

	age := AgeFromBirthDate(birthDate, time.Now())
	if age < MinAge || age > MaxAge {
		return NewValidationError(birthDateField, fmt.Sprintf("age must be between %d and %d", MinAge, MaxAge))
	}

	return nil
}

func ValidateProfile(profile *Profile) error {
	if profile == nil {
		return NewValidationError("profile", msgRequired)
	}
	if err := validateRange(ageField, profile.Age, MinAge, MaxAge); err != nil {
		return err
	}
	if profile.Gender != GenderMale && profile.Gender != GenderFemale {
		return NewValidationError(genderField, "must be specified")
	}
	if err := validateRange(heightField, profile.HeightCM, MinHeightCM, MaxHeightCM); err != nil {
		return err
	}
	if err := validateFloatRange(weightField, profile.WeightKG, MinWeightKG, MaxWeightKG); err != nil {
		return err
	}
	if profile.TrainingLevel < TrainingLevelBeginner || profile.TrainingLevel > TrainingLevelAdvanced {
		return NewValidationError(trainingLevelField, "must be specified")
	}
	if profile.ActivityLevel < ActivityLevelSedentary || profile.ActivityLevel > ActivityLevelHigh {
		return NewValidationError(activityLevelField, "must be specified")
	}
	if profile.Goal < FitnessGoalLoseWeight || profile.Goal > FitnessGoalGainMuscle {
		return NewValidationError(goalField, "must be specified")
	}
	if profile.TargetWeightKG != nil {
		if err := validateFloatRange(targetWeightField, *profile.TargetWeightKG, MinWeightKG, MaxWeightKG); err != nil {
			return err
		}
	}
	if profile.TrainingLocation != TrainingLocationHome && profile.TrainingLocation != TrainingLocationGym {
		return NewValidationError(trainingLocationField, "must be specified")
	}
	if err := validateRange(trainingDaysPerWeekField, profile.TrainingDaysPerWeek, 0, 7); err != nil {
		return err
	}
	if err := validateStringList("allergies", profile.Allergies); err != nil {
		return err
	}
	if err := validateStringList("excluded_foods", profile.ExcludedFoods); err != nil {
		return err
	}
	if err := validateStringList("food_preferences", profile.FoodPreferences); err != nil {
		return err
	}
	if utf8.RuneCountInString(strings.TrimSpace(profile.Equipment)) > MaxEquipmentLength {
		return NewValidationError(equipmentField, msgTooLong)
	}

	return nil
}

func AgeFromBirthDate(birthDate, now time.Time) int {
	age := now.Year() - birthDate.Year()
	if now.Month() < birthDate.Month() ||
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age
}

func hasLetterAndDigit(value string) bool {
	hasLetter := false
	hasDigit := false

	for _, r := range value {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}

func validateRange(field string, value, minValue, maxValue int) error {
	if value < minValue || value > maxValue {
		return NewValidationError(field, fmt.Sprintf("must be between %d and %d", minValue, maxValue))
	}

	return nil
}

func validateFloatRange(field string, value, minValue, maxValue float64) error {
	if value < minValue || value > maxValue {
		return NewValidationError(field, fmt.Sprintf("must be between %.0f and %.0f", minValue, maxValue))
	}

	return nil
}

func validateStringList(field string, values []string) error {
	if len(values) > MaxStringListItems {
		return NewValidationError(field, fmt.Sprintf("must contain no more than %d items", MaxStringListItems))
	}

	for _, value := range values {
		item := strings.TrimSpace(value)
		if item == "" {
			return NewValidationError(field, "items cannot be empty")
		}
		if utf8.RuneCountInString(item) > MaxStringListItemLength {
			return NewValidationError(field, fmt.Sprintf("items must contain no more than %d characters", MaxStringListItemLength))
		}
	}

	return nil
}
