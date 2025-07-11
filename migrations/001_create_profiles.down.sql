-- Drop profiles table and related objects
DROP TRIGGER IF EXISTS update_profiles_updated_at ON profiles;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS profiles;