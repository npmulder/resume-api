-- Drop achievements table
DROP TRIGGER IF EXISTS update_achievements_updated_at ON achievements;
DROP TABLE IF EXISTS achievements;