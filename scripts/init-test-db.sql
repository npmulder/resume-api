-- Initial test database setup script
-- This file is run when the PostgreSQL container starts

-- Create additional databases if needed
-- (The main test database is created via POSTGRES_DB environment variable)

-- Set up any initial configuration
SET timezone = 'UTC';

-- Log the initialization
\echo 'Test database initialized successfully'
\echo 'Database: resume_api_test'
\echo 'User: dev'
\echo 'Ready for migrations and testing'