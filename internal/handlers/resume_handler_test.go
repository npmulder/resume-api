package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockResumeService is a mock implementation of the ResumeService interface
type MockResumeService struct {
	mock.Mock
}

func (m *MockResumeService) GetProfile(ctx context.Context) (*models.Profile, error) {
	args := m.Called(ctx)
	profile, _ := args.Get(0).(*models.Profile)
	return profile, args.Error(1)
}

func (m *MockResumeService) GetExperiences(ctx context.Context, filters repository.ExperienceFilters) ([]*models.Experience, error) {
	args := m.Called(ctx, filters)
	experiences, _ := args.Get(0).([]*models.Experience)
	return experiences, args.Error(1)
}

func (m *MockResumeService) GetSkills(ctx context.Context, filters repository.SkillFilters) ([]*models.Skill, error) {
	args := m.Called(ctx, filters)
	skills, _ := args.Get(0).([]*models.Skill)
	return skills, args.Error(1)
}

func (m *MockResumeService) GetAchievements(ctx context.Context, filters repository.AchievementFilters) ([]*models.Achievement, error) {
	args := m.Called(ctx, filters)
	achievements, _ := args.Get(0).([]*models.Achievement)
	return achievements, args.Error(1)
}

func (m *MockResumeService) GetEducation(ctx context.Context, filters repository.EducationFilters) ([]*models.Education, error) {
	args := m.Called(ctx, filters)
	education, _ := args.Get(0).([]*models.Education)
	return education, args.Error(1)
}

func (m *MockResumeService) GetProjects(ctx context.Context, filters repository.ProjectFilters) ([]*models.Project, error) {
	args := m.Called(ctx, filters)
	projects, _ := args.Get(0).([]*models.Project)
	return projects, args.Error(1)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestGetProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup
		router := setupRouter()
		mockService := new(MockResumeService)
		handler := NewResumeHandler(mockService)

		expectedProfile := &models.Profile{
			ID:    1,
			Name:  "John Doe",
			Title: "Software Engineer",
			Email: "john@example.com",
		}

		// Configure mock
		mockService.On("GetProfile", mock.Anything).Return(expectedProfile, nil)

		// Setup route
		router.GET("/api/v1/profile", handler.GetProfile)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Profile
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedProfile.ID, response.ID)
		assert.Equal(t, expectedProfile.Name, response.Name)
		assert.Equal(t, expectedProfile.Title, response.Title)
		assert.Equal(t, expectedProfile.Email, response.Email)

		// Verify mock expectations
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		// Setup
		router := setupRouter()
		mockService := new(MockResumeService)
		handler := NewResumeHandler(mockService)

		// Configure mock
		mockService.On("GetProfile", mock.Anything).Return(nil, repository.ErrNotFound)

		// Setup route
		router.GET("/api/v1/profile", handler.GetProfile)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Profile not found")

		// Verify mock expectations
		mockService.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		// Setup
		router := setupRouter()
		mockService := new(MockResumeService)
		handler := NewResumeHandler(mockService)

		// Configure mock
		mockService.On("GetProfile", mock.Anything).Return(nil, errors.New("database error"))

		// Setup route
		router.GET("/api/v1/profile", handler.GetProfile)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "An unexpected error occurred")

		// Verify mock expectations
		mockService.AssertExpectations(t)
	})
}

func TestGetExperiences(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup
		router := setupRouter()
		mockService := new(MockResumeService)
		handler := NewResumeHandler(mockService)

		expectedExperiences := []*models.Experience{
			{
				ID:       1,
				Company:  "Example Corp",
				Position: "Software Engineer",
			},
		}

		// Configure mock to match any filters
		mockService.On("GetExperiences", mock.Anything, mock.AnythingOfType("repository.ExperienceFilters")).Return(expectedExperiences, nil)

		// Setup route
		router.GET("/api/v1/experiences", handler.GetExperiences)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/experiences?company=Example&limit=10", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.Experience
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, expectedExperiences[0].ID, response[0].ID)
		assert.Equal(t, expectedExperiences[0].Company, response[0].Company)
		assert.Equal(t, expectedExperiences[0].Position, response[0].Position)

		// Verify mock expectations
		mockService.AssertExpectations(t)
	})

	// Note: We're not testing invalid query parameters because Gin's binding
	// behavior for int fields with invalid values is to set them to 0, not fail

	t.Run("not found", func(t *testing.T) {
		// Setup
		router := setupRouter()
		mockService := new(MockResumeService)
		handler := NewResumeHandler(mockService)

		// Configure mock
		mockService.On("GetExperiences", mock.Anything, mock.AnythingOfType("repository.ExperienceFilters")).Return(nil, repository.ErrNotFound)

		// Setup route
		router.GET("/api/v1/experiences", handler.GetExperiences)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/experiences", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "No experiences found matching the criteria")

		// Verify mock expectations
		mockService.AssertExpectations(t)
	})
}

func TestGetSkills(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup
		router := setupRouter()
		mockService := new(MockResumeService)
		handler := NewResumeHandler(mockService)

		level := "expert"
		expectedSkills := []*models.Skill{
			{
				ID:       1,
				Name:     "Go",
				Category: "Programming Languages",
				Level:    &level,
			},
		}

		// Configure mock
		mockService.On("GetSkills", mock.Anything, mock.AnythingOfType("repository.SkillFilters")).Return(expectedSkills, nil)

		// Setup route
		router.GET("/api/v1/skills", handler.GetSkills)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/skills?category=Programming+Languages&limit=10", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.Skill
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, expectedSkills[0].ID, response[0].ID)
		assert.Equal(t, expectedSkills[0].Name, response[0].Name)
		assert.Equal(t, expectedSkills[0].Category, response[0].Category)
		assert.Equal(t, expectedSkills[0].Level, response[0].Level)

		// Verify mock expectations
		mockService.AssertExpectations(t)
	})
}

func TestGetAchievements(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup
		router := setupRouter()
		mockService := new(MockResumeService)
		handler := NewResumeHandler(mockService)

		description := "Improved system performance by 50%"
		yearAchieved := 2023
		expectedAchievements := []*models.Achievement{
			{
				ID:           1,
				Title:        "Performance Improvement",
				Description:  &description,
				YearAchieved: &yearAchieved,
			},
		}

		// Configure mock
		mockService.On("GetAchievements", mock.Anything, mock.AnythingOfType("repository.AchievementFilters")).Return(expectedAchievements, nil)

		// Setup route
		router.GET("/api/v1/achievements", handler.GetAchievements)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/achievements?year=2023&limit=10", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.Achievement
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, expectedAchievements[0].ID, response[0].ID)
		assert.Equal(t, expectedAchievements[0].Title, response[0].Title)
		assert.Equal(t, *expectedAchievements[0].Description, *response[0].Description)
		assert.Equal(t, *expectedAchievements[0].YearAchieved, *response[0].YearAchieved)

		// Verify mock expectations
		mockService.AssertExpectations(t)
	})
}

func TestGetEducation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup
		router := setupRouter()
		mockService := new(MockResumeService)
		handler := NewResumeHandler(mockService)

		fieldOfStudy := "Computer Science"
		expectedEducation := []*models.Education{
			{
				ID:                    1,
				Institution:           "University of Example",
				DegreeOrCertification: "Bachelor of Science",
				FieldOfStudy:          &fieldOfStudy,
				Type:                  "education",
				Status:                "completed",
			},
		}

		// Configure mock
		mockService.On("GetEducation", mock.Anything, mock.AnythingOfType("repository.EducationFilters")).Return(expectedEducation, nil)

		// Setup route
		router.GET("/api/v1/education", handler.GetEducation)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/education?type=degree&limit=10", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.Education
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, expectedEducation[0].ID, response[0].ID)
		assert.Equal(t, expectedEducation[0].Institution, response[0].Institution)
		assert.Equal(t, expectedEducation[0].DegreeOrCertification, response[0].DegreeOrCertification)
		assert.Equal(t, *expectedEducation[0].FieldOfStudy, *response[0].FieldOfStudy)
		assert.Equal(t, expectedEducation[0].Type, response[0].Type)
		assert.Equal(t, expectedEducation[0].Status, response[0].Status)

		// Verify mock expectations
		mockService.AssertExpectations(t)
	})
}

func TestGetProjects(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Setup
		router := setupRouter()
		mockService := new(MockResumeService)
		handler := NewResumeHandler(mockService)

		expectedProjects := []*models.Project{
			{
				ID:   1,
				Name: "Resume API",
			},
		}

		// Configure mock
		mockService.On("GetProjects", mock.Anything, mock.AnythingOfType("repository.ProjectFilters")).Return(expectedProjects, nil)

		// Setup route
		router.GET("/api/v1/projects", handler.GetProjects)

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/api/v1/projects?status=active&limit=10", nil)
		w := httptest.NewRecorder()

		// Serve request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.Project
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, expectedProjects[0].ID, response[0].ID)
		assert.Equal(t, expectedProjects[0].Name, response[0].Name)

		// Verify mock expectations
		mockService.AssertExpectations(t)
	})
}
