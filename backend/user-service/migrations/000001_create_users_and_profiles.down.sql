DROP TRIGGER IF EXISTS user_profiles_unset_profile_completed ON user_profiles;
DROP TRIGGER IF EXISTS user_profiles_set_profile_completed ON user_profiles;
DROP TRIGGER IF EXISTS user_profiles_set_updated_at ON user_profiles;
DROP TRIGGER IF EXISTS users_set_updated_at ON users;

DROP FUNCTION IF EXISTS unset_profile_completed();
DROP FUNCTION IF EXISTS set_profile_completed();
DROP FUNCTION IF EXISTS set_updated_at();

DROP INDEX IF EXISTS user_profiles_goal_idx;
DROP INDEX IF EXISTS user_profiles_activity_level_idx;
DROP INDEX IF EXISTS user_profiles_training_level_idx;
DROP INDEX IF EXISTS users_email_unique_idx;

DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS users;
