-- Drop projects table
DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
DROP TABLE IF EXISTS projects;