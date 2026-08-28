CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    birth_date DATE NOT NULL,
    profile_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_not_blank CHECK (btrim(email) <> ''),
    CONSTRAINT users_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT users_password_hash_not_blank CHECK (btrim(password_hash) <> '')
);

CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    age INTEGER NOT NULL,
    gender SMALLINT NOT NULL,
    height_cm INTEGER NOT NULL,
    weight_kg NUMERIC(5, 2) NOT NULL,
    training_level SMALLINT NOT NULL,
    activity_level SMALLINT NOT NULL,
    goal SMALLINT NOT NULL,
    target_weight_kg NUMERIC(5, 2),
    allergies TEXT[] NOT NULL DEFAULT '{}',
    excluded_foods TEXT[] NOT NULL DEFAULT '{}',
    food_preferences TEXT[] NOT NULL DEFAULT '{}',
    training_location SMALLINT NOT NULL,
    training_days_per_week INTEGER NOT NULL,
    equipment TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT user_profiles_age_positive CHECK (age > 0),
    CONSTRAINT user_profiles_gender_valid CHECK (gender IN (0, 1, 2)),
    CONSTRAINT user_profiles_height_positive CHECK (height_cm > 0),
    CONSTRAINT user_profiles_weight_positive CHECK (weight_kg > 0),
    CONSTRAINT user_profiles_training_level_valid CHECK (training_level IN (0, 1, 2, 3)),
    CONSTRAINT user_profiles_activity_level_valid CHECK (activity_level IN (0, 1, 2, 3, 4)),
    CONSTRAINT user_profiles_goal_valid CHECK (goal IN (0, 1, 2, 3)),
    CONSTRAINT user_profiles_target_weight_positive CHECK (target_weight_kg IS NULL OR target_weight_kg > 0),
    CONSTRAINT user_profiles_training_location_valid CHECK (training_location IN (0, 1, 2)),
    CONSTRAINT user_profiles_training_days_per_week_valid CHECK (training_days_per_week BETWEEN 0 AND 7)
);

CREATE UNIQUE INDEX users_email_unique_idx ON users (lower(email));
CREATE INDEX user_profiles_training_level_idx ON user_profiles(training_level);
CREATE INDEX user_profiles_activity_level_idx ON user_profiles(activity_level);
CREATE INDEX user_profiles_goal_idx ON user_profiles(goal);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER user_profiles_set_updated_at
BEFORE UPDATE ON user_profiles
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE FUNCTION set_profile_completed()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET profile_completed = TRUE
    WHERE id = NEW.user_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION unset_profile_completed()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET profile_completed = FALSE
    WHERE id = OLD.user_id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_profiles_set_profile_completed
AFTER INSERT OR UPDATE ON user_profiles
FOR EACH ROW
EXECUTE FUNCTION set_profile_completed();

CREATE TRIGGER user_profiles_unset_profile_completed
AFTER DELETE ON user_profiles
FOR EACH ROW
EXECUTE FUNCTION unset_profile_completed();
