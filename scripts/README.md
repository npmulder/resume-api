# Database Seeding Scripts

This directory contains scripts for populating the resume database with data.

## Files

- **`seed.go`** - Main seeding script that reads JSON data and populates the database
- **`seed-data.example.json`** - Example seed data structure (safe for GitHub)
- **`seed-data.json`** - Actual seed data (gitignored, contains personal information)

## Usage

### 1. Prepare Your Seed Data

Copy the example file and customize it with your information:
```bash
cp scripts/seed-data.example.json scripts/seed-data.json
# Edit scripts/seed-data.json with your actual resume data
```

### 2. Set Up Database

Ensure PostgreSQL is running and create the database:
```bash
# Using Docker
docker run --name resume-postgres \
  -e POSTGRES_DB=resume_api_dev \
  -e POSTGRES_USER=dev \
  -e POSTGRES_PASSWORD=devpass \
  -p 5432:5432 -d postgres:15

# Or using local PostgreSQL
createdb resume_api_dev
```

### 3. Run Migrations

Apply database schema:
```bash
go run cmd/migrate/main.go up
```

### 4. Seed the Database

Run the seeding script:
```bash
# Uses scripts/seed-data.json if it exists, otherwise falls back to example
go run scripts/seed.go

# Specify custom seed file
SEED_FILE=scripts/my-custom-data.json go run scripts/seed.go

# Use custom database URL
DATABASE_URL=postgres://user:pass@host:port/dbname go run scripts/seed.go
```

## Seed Data Structure

The JSON file contains the following sections:

### Profile
Personal information and contact details:
```json
{
  "profile": {
    "name": "Your Name",
    "title": "Your Job Title",
    "email": "your.email@example.com",
    "phone": "+1 (555) 123-4567",
    "location": "City, Country",
    "linkedin": "https://linkedin.com/in/yourname",
    "github": "https://github.com/yourname",
    "summary": "Your professional summary..."
  }
}
```

### Experiences
Work history with highlights:
```json
{
  "experiences": [
    {
      "company": "Company Name",
      "position": "Job Title",
      "start_date": "2023-01-01",
      "end_date": null,
      "description": "Brief job description",
      "highlights": [
        "Key achievement 1",
        "Key achievement 2"
      ],
      "order": 1
    }
  ]
}
```

### Skills
Categorized skills with levels:
```json
{
  "skills": [
    {
      "category": "Programming Languages",
      "name": "Go",
      "level": "intermediate",
      "order": 1,
      "featured": true
    }
  ]
}
```

### Achievements
Key accomplishments with impact metrics:
```json
{
  "achievements": [
    {
      "title": "Performance Optimization",
      "description": "Led initiative to optimize application performance",
      "category": "performance",
      "impact": "50% improvement in response time",
      "year": 2023,
      "order": 1,
      "featured": true
    }
  ]
}
```

### Education
Education and certifications:
```json
{
  "education": [
    {
      "institution": "University Name",
      "degree": "Degree Name",
      "field": "Field of Study",
      "year_completed": 2020,
      "year_started": 2016,
      "description": "Description of studies",
      "type": "education",
      "status": "completed",
      "credential_id": "",
      "credential_url": "",
      "order": 1,
      "featured": true
    }
  ]
}
```

### Projects
Notable projects and portfolio items:
```json
{
  "projects": [
    {
      "name": "Project Name",
      "description": "Detailed project description",
      "short_description": "Brief summary",
      "technologies": ["Go", "PostgreSQL", "Docker"],
      "github_url": "https://github.com/user/project",
      "demo_url": "https://demo.example.com",
      "start_date": "2023-01-01",
      "end_date": null,
      "status": "active",
      "is_featured": true,
      "order": 1,
      "key_features": [
        "Feature 1",
        "Feature 2"
      ]
    }
  ]
}
```

## Environment Variables

- **`DATABASE_URL`** - PostgreSQL connection string (default: `postgres://dev:devpass@localhost:5432/resume_api_dev?sslmode=disable`)
- **`SEED_FILE`** - Path to seed data file (default: `scripts/seed-data.json`)

## Notes

- The script automatically handles transactions - if any part fails, all changes are rolled back
- Existing data is cleared before seeding (except profiles, which are upserted by email)
- If `seed-data.json` doesn't exist, the script automatically falls back to the example file
- Date fields use ISO format: `YYYY-MM-DD`
- Arrays are stored as PostgreSQL arrays in the database
- Technologies for projects are stored as JSONB