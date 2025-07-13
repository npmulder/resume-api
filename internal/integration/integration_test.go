package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npmulder/resume-api/internal/config"
	"github.com/npmulder/resume-api/internal/database"
	"github.com/npmulder/resume-api/internal/handlers"
	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/repository"
	"github.com/npmulder/resume-api/internal/repository/postgres"
	"github.com/npmulder/resume-api/internal/services"
)

// TestDB represents a test database connection
type TestDB struct {
	*database.DB
	cleanup func()
}

// setupTestDB creates a test database connection with proper cleanup
func setupTestDB(t *testing.T) *TestDB {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	cfg := &config.DatabaseConfig{
		Host:               getTestEnv("TEST_DB_HOST", "localhost"),
		Port:               5432,
		Name:               getTestEnv("TEST_DB_NAME", "resume_api_test"),
		User:               getTestEnv("TEST_DB_USER", "dev"),
		Password:           getTestEnv("TEST_DB_PASSWORD", "devpass"),
		SSLMode:            "disable",
		MaxConnections:     5,
		MaxIdleConnections: 2,
		ConnMaxLifetime:    30 * time.Minute,
		ConnMaxIdleTime:    5 * time.Minute,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.New(ctx, cfg, nil)
	require.NoError(t, err, "Failed to connect to test database")

	// Test the connection
	err = db.Ping(ctx)
	require.NoError(t, err, "Failed to ping test database")

	return &TestDB{
		DB: db,
		cleanup: func() {
			db.Close()
		},
	}
}

// Close cleans up the test database connection
func (tdb *TestDB) Close() {
	tdb.cleanup()
}

// CleanupTables removes all data from tables for clean test state
func (tdb *TestDB) CleanupTables(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Clean tables in correct order due to potential foreign keys
	tables := []string{
		"projects",
		"education",
		"achievements",
		"skills",
		"experiences",
		"profiles",
	}

	for _, table := range tables {
		_, err := tdb.Pool().Exec(ctx, "DELETE FROM "+table)
		require.NoError(t, err, "Failed to clean table: %s", table)
	}
}

// getTestEnv gets environment variable for tests with fallback
func getTestEnv(key, fallback string) string {
	value := getenv(key)
	if value != "" {
		return value
	}
	return fallback
}

// getenv is a wrapper for os.Getenv to make it testable
var getenv = func(key string) string {
	return ""
}

// setupTestApp creates a test application with real repositories, services, and handlers
func setupTestApp(t *testing.T, db *database.DB) (*gin.Engine, *repository.Repositories) {
	// Create repositories
	profileRepo := postgres.NewProfileRepository(db.Pool())
	experienceRepo := postgres.NewExperienceRepository(db.Pool())
	skillRepo := postgres.NewSkillRepository(db.Pool())
	achievementRepo := postgres.NewAchievementRepository(db.Pool())
	educationRepo := postgres.NewEducationRepository(db.Pool())
	projectRepo := postgres.NewProjectRepository(db.Pool())

	repos := &repository.Repositories{
		Profile:     profileRepo,
		Experience:  experienceRepo,
		Skill:       skillRepo,
		Achievement: achievementRepo,
		Education:   educationRepo,
		Project:     projectRepo,
	}

	// Create service
	resumeService := services.NewResumeService(*repos)

	// Create handler
	resumeHandler := handlers.NewResumeHandler(resumeService)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Register routes
	router.GET("/api/v1/profile", resumeHandler.GetProfile)
	router.GET("/api/v1/experiences", resumeHandler.GetExperiences)
	router.GET("/api/v1/skills", resumeHandler.GetSkills)
	router.GET("/api/v1/achievements", resumeHandler.GetAchievements)
	router.GET("/api/v1/education", resumeHandler.GetEducation)
	router.GET("/api/v1/projects", resumeHandler.GetProjects)

	return router, repos
}

// Helper functions for creating test data
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestProfileEndToEnd(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()
	testDB.CleanupTables(t)

	router, repos := setupTestApp(t, testDB.DB)

	// Create test profile
	ctx := context.Background()
	profile := &models.Profile{
		Name:     "John Doe",
		Title:    "Software Engineer",
		Email:    "john.doe@example.com",
		Phone:    stringPtr("+1-555-0123"),
		Location: stringPtr("San Francisco, CA"),
		LinkedIn: stringPtr("https://linkedin.com/in/johndoe"),
		GitHub:   stringPtr("https://github.com/johndoe"),
		Summary:  stringPtr("Experienced software engineer with 5+ years in Go development"),
	}

	err := repos.Profile.CreateProfile(ctx, profile)
	require.NoError(t, err)

	// Test GET /api/v1/profile
	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseProfile models.Profile
	err = json.Unmarshal(w.Body.Bytes(), &responseProfile)
	require.NoError(t, err)

	assert.Equal(t, profile.ID, responseProfile.ID)
	assert.Equal(t, profile.Name, responseProfile.Name)
	assert.Equal(t, profile.Title, responseProfile.Title)
	assert.Equal(t, profile.Email, responseProfile.Email)
	assert.Equal(t, *profile.Phone, *responseProfile.Phone)
	assert.Equal(t, *profile.Location, *responseProfile.Location)
	assert.Equal(t, *profile.LinkedIn, *responseProfile.LinkedIn)
	assert.Equal(t, *profile.GitHub, *responseProfile.GitHub)
	assert.Equal(t, *profile.Summary, *responseProfile.Summary)
}

func TestExperiencesEndToEnd(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()
	testDB.CleanupTables(t)

	router, repos := setupTestApp(t, testDB.DB)

	// Create test experiences
	ctx := context.Background()
	experiences := []*models.Experience{
		{
			Company:     "Google",
			Position:    "Senior Software Engineer",
			StartDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     timePtr(time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)),
			Description: stringPtr("Led development of cloud services"),
			Highlights: []string{
				"Improved system performance by 30%",
				"Mentored junior engineers",
			},
			OrderIndex: 0,
		},
		{
			Company:     "Microsoft",
			Position:    "Principal Engineer",
			StartDate:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     nil, // Current position
			Description: stringPtr("Leading architecture for distributed systems"),
			Highlights: []string{
				"Designed microservices architecture",
				"Implemented CI/CD pipeline",
			},
			OrderIndex: 1,
		},
	}

	for _, exp := range experiences {
		err := repos.Experience.CreateExperience(ctx, exp)
		require.NoError(t, err)
	}

	// Test GET /api/v1/experiences
	req := httptest.NewRequest(http.MethodGet, "/api/v1/experiences", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseExperiences []*models.Experience
	err := json.Unmarshal(w.Body.Bytes(), &responseExperiences)
	require.NoError(t, err)

	assert.Len(t, responseExperiences, 2)
	// Should be ordered by start_date DESC (most recent first)
	assert.Equal(t, "Microsoft", responseExperiences[0].Company)
	assert.Equal(t, "Google", responseExperiences[1].Company)

	// Test filtering by company
	req = httptest.NewRequest(http.MethodGet, "/api/v1/experiences?company=Google", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &responseExperiences)
	require.NoError(t, err)

	assert.Len(t, responseExperiences, 1)
	assert.Equal(t, "Google", responseExperiences[0].Company)

	// Test filtering by current position
	req = httptest.NewRequest(http.MethodGet, "/api/v1/experiences?current=true", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &responseExperiences)
	require.NoError(t, err)

	assert.Len(t, responseExperiences, 1)
	assert.Equal(t, "Microsoft", responseExperiences[0].Company)
	assert.Nil(t, responseExperiences[0].EndDate)
}

func TestSkillsEndToEnd(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()
	testDB.CleanupTables(t)

	router, repos := setupTestApp(t, testDB.DB)

	// Create test skills
	ctx := context.Background()
	skills := []*models.Skill{
		{
			Name:        "Go",
			Category:    "Languages",
			Level:       stringPtr("Expert"),
			YearsExperience:  intPtr(5),
			IsFeatured:  true,
			OrderIndex:  0,
			Description: stringPtr("Primary programming language"),
		},
		{
			Name:        "PostgreSQL",
			Category:    "Databases",
			Level:       stringPtr("Advanced"),
			YearsExperience:  intPtr(4),
			IsFeatured:  true,
			OrderIndex:  1,
			Description: stringPtr("Primary database"),
		},
		{
			Name:        "Docker",
			Category:    "DevOps",
			Level:       stringPtr("Intermediate"),
			YearsExperience:  intPtr(3),
			IsFeatured:  false,
			OrderIndex:  2,
			Description: stringPtr("Containerization"),
		},
	}

	for _, skill := range skills {
		err := repos.Skill.CreateSkill(ctx, skill)
		require.NoError(t, err)
	}

	// Test GET /api/v1/skills
	req := httptest.NewRequest(http.MethodGet, "/api/v1/skills", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseSkills []*models.Skill
	err := json.Unmarshal(w.Body.Bytes(), &responseSkills)
	require.NoError(t, err)

	assert.Len(t, responseSkills, 3)

	// Test filtering by category
	req = httptest.NewRequest(http.MethodGet, "/api/v1/skills?category=Languages", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &responseSkills)
	require.NoError(t, err)

	assert.Len(t, responseSkills, 1)
	assert.Equal(t, "Go", responseSkills[0].Name)

	// Test filtering by featured
	req = httptest.NewRequest(http.MethodGet, "/api/v1/skills?featured=true", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &responseSkills)
	require.NoError(t, err)

	assert.Len(t, responseSkills, 2)
	for _, skill := range responseSkills {
		assert.True(t, skill.IsFeatured)
	}
}

func TestAchievementsEndToEnd(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()
	testDB.CleanupTables(t)

	router, repos := setupTestApp(t, testDB.DB)

	// Create test achievements
	ctx := context.Background()
	achievements := []*models.Achievement{
		{
			Title:        "Performance Optimization",
			Description:  stringPtr("Improved system performance by 50%"),
			YearAchieved: intPtr(2023),
			Category:     stringPtr("Technical"),
			IsFeatured:   true,
			OrderIndex:   0,
		},
		{
			Title:        "Team Leadership",
			Description:  stringPtr("Led a team of 5 engineers"),
			YearAchieved: intPtr(2022),
			Category:     stringPtr("Leadership"),
			IsFeatured:   true,
			OrderIndex:   1,
		},
		{
			Title:        "Conference Speaker",
			Description:  stringPtr("Spoke at GopherCon 2021"),
			YearAchieved: intPtr(2021),
			Category:     stringPtr("Community"),
			IsFeatured:   false,
			OrderIndex:   2,
		},
	}

	for _, achievement := range achievements {
		err := repos.Achievement.CreateAchievement(ctx, achievement)
		require.NoError(t, err)
	}

	// Test GET /api/v1/achievements
	req := httptest.NewRequest(http.MethodGet, "/api/v1/achievements", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseAchievements []*models.Achievement
	err := json.Unmarshal(w.Body.Bytes(), &responseAchievements)
	require.NoError(t, err)

	assert.Len(t, responseAchievements, 3)

	// Test filtering by year
	req = httptest.NewRequest(http.MethodGet, "/api/v1/achievements?year=2023", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &responseAchievements)
	require.NoError(t, err)

	assert.Len(t, responseAchievements, 1)
	assert.Equal(t, "Performance Optimization", responseAchievements[0].Title)

	// Test filtering by featured
	req = httptest.NewRequest(http.MethodGet, "/api/v1/achievements?featured=true", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &responseAchievements)
	require.NoError(t, err)

	assert.Len(t, responseAchievements, 2)
	for _, achievement := range responseAchievements {
		assert.True(t, achievement.IsFeatured)
	}
}

func TestEducationEndToEnd(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()
	testDB.CleanupTables(t)

	router, repos := setupTestApp(t, testDB.DB)

	// Create test education
	ctx := context.Background()
	fieldOfStudy := "Computer Science"
	education := []*models.Education{
		{
			Institution:           "Stanford University",
			DegreeOrCertification: "Bachelor of Science",
			FieldOfStudy:          &fieldOfStudy,
			StartDate:             timePtr(time.Date(2015, 9, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:               timePtr(time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC)),
			Type:                  "education",
			Status:                "completed",
			IsFeatured:            true,
			OrderIndex:            0,
		},
		{
			Institution:           "AWS",
			DegreeOrCertification: "AWS Certified Solutions Architect",
			StartDate:             timePtr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:               timePtr(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
			Type:                  "certification",
			Status:                "active",
			IsFeatured:            true,
			OrderIndex:            1,
		},
		{
			Institution:           "Coursera",
			DegreeOrCertification: "Machine Learning",
			StartDate:             timePtr(time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:               timePtr(time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC)),
			Type:                  "course",
			Status:                "completed",
			IsFeatured:            false,
			OrderIndex:            2,
		},
	}

	for _, edu := range education {
		err := repos.Education.CreateEducation(ctx, edu)
		require.NoError(t, err)
	}

	// Test GET /api/v1/education
	req := httptest.NewRequest(http.MethodGet, "/api/v1/education", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseEducation []*models.Education
	err := json.Unmarshal(w.Body.Bytes(), &responseEducation)
	require.NoError(t, err)

	assert.Len(t, responseEducation, 3)

	// Test filtering by type
	req = httptest.NewRequest(http.MethodGet, "/api/v1/education?type=certification", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &responseEducation)
	require.NoError(t, err)

	assert.Len(t, responseEducation, 1)
	assert.Equal(t, "AWS Certified Solutions Architect", responseEducation[0].DegreeOrCertification)

	// Test filtering by featured
	req = httptest.NewRequest(http.MethodGet, "/api/v1/education?featured=true", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &responseEducation)
	require.NoError(t, err)

	assert.Len(t, responseEducation, 2)
	for _, edu := range responseEducation {
		assert.True(t, edu.IsFeatured)
	}
}

func TestProjectsEndToEnd(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()
	testDB.CleanupTables(t)

	router, repos := setupTestApp(t, testDB.DB)

	// Create test projects
	ctx := context.Background()
	projects := []*models.Project{
		{
			Name:        "Resume API",
			Description: stringPtr("RESTful API for resume data"),
			GitHubURL:   stringPtr("https://github.com/example/resume-api"),
			Status:      "active",
			IsFeatured:  true,
			OrderIndex:  0,
			Technologies: []string{"Go", "PostgreSQL", "Docker", "Kubernetes"},
			KeyFeatures: []string{"RESTful API", "Clean Architecture", "Integration Tests"},
		},
		{
			Name:        "Personal Website",
			Description: stringPtr("Portfolio website"),
			DemoURL:     stringPtr("https://example.com"),
			Status:      "completed",
			IsFeatured:  true,
			OrderIndex:  1,
			Technologies: []string{"JavaScript", "HTML", "CSS", "React"},
			KeyFeatures: []string{"Responsive Design", "Portfolio", "Contact Form"},
		},
		{
			Name:        "Side Project",
			Description: stringPtr("Experimental project"),
			Status:      "planning",
			IsFeatured:  false,
			OrderIndex:  2,
			Technologies: []string{"Rust"},
			KeyFeatures: []string{"Experimental", "Learning Project"},
		},
	}

	for _, project := range projects {
		err := repos.Project.CreateProject(ctx, project)
		require.NoError(t, err)
	}

	// Test GET /api/v1/projects
	req := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseProjects []*models.Project
	err := json.Unmarshal(w.Body.Bytes(), &responseProjects)
	require.NoError(t, err)

	assert.Len(t, responseProjects, 3)

	// Test filtering by status
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects?status=active", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &responseProjects)
	require.NoError(t, err)

	assert.Len(t, responseProjects, 1)
	assert.Equal(t, "Resume API", responseProjects[0].Name)

	// Test filtering by featured
	req = httptest.NewRequest(http.MethodGet, "/api/v1/projects?featured=true", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &responseProjects)
	require.NoError(t, err)

	assert.Len(t, responseProjects, 2)
	for _, project := range responseProjects {
		assert.True(t, project.IsFeatured)
	}
}
