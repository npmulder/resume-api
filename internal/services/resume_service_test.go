package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/repository"
)

// Mock Repositories (as before)

type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) GetProfile(ctx context.Context) (*models.Profile, error) {
	args := m.Called(ctx)
	profile, _ := args.Get(0).(*models.Profile)
	return profile, args.Error(1)
}

func (m *MockProfileRepository) UpdateProfile(ctx context.Context, profile *models.Profile) error {
	return m.Called(ctx, profile).Error(0)
}

func (m *MockProfileRepository) CreateProfile(ctx context.Context, profile *models.Profile) error {
	return m.Called(ctx, profile).Error(0)
}

type MockExperienceRepository struct {
	mock.Mock
}

func (m *MockExperienceRepository) GetExperiences(ctx context.Context, filters repository.ExperienceFilters) ([]*models.Experience, error) {
	args := m.Called(ctx, filters)
	experiences, _ := args.Get(0).([]*models.Experience)
	return experiences, args.Error(1)
}

func (m *MockExperienceRepository) GetExperienceByID(ctx context.Context, id int) (*models.Experience, error) {
	args := m.Called(ctx, id)
	experience, _ := args.Get(0).(*models.Experience)
	return experience, args.Error(1)
}

func (m *MockExperienceRepository) CreateExperience(ctx context.Context, experience *models.Experience) error {
	return m.Called(ctx, experience).Error(0)
}

func (m *MockExperienceRepository) UpdateExperience(ctx context.Context, experience *models.Experience) error {
	return m.Called(ctx, experience).Error(0)
}

func (m *MockExperienceRepository) DeleteExperience(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

type MockSkillRepository struct {
	mock.Mock
}

func (m *MockSkillRepository) GetSkills(ctx context.Context, filters repository.SkillFilters) ([]*models.Skill, error) {
	args := m.Called(ctx, filters)
	skills, _ := args.Get(0).([]*models.Skill)
	return skills, args.Error(1)
}

func (m *MockSkillRepository) GetSkillsByCategory(ctx context.Context, category string) ([]*models.Skill, error) {
	args := m.Called(ctx, category)
	skills, _ := args.Get(0).([]*models.Skill)
	return skills, args.Error(1)
}

func (m *MockSkillRepository) GetFeaturedSkills(ctx context.Context) ([]*models.Skill, error) {
	args := m.Called(ctx)
	skills, _ := args.Get(0).([]*models.Skill)
	return skills, args.Error(1)
}

func (m *MockSkillRepository) CreateSkill(ctx context.Context, skill *models.Skill) error {
	return m.Called(ctx, skill).Error(0)
}

func (m *MockSkillRepository) UpdateSkill(ctx context.Context, skill *models.Skill) error {
	return m.Called(ctx, skill).Error(0)
}

func (m *MockSkillRepository) DeleteSkill(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

type MockAchievementRepository struct {
	mock.Mock
}

func (m *MockAchievementRepository) GetAchievements(ctx context.Context, filters repository.AchievementFilters) ([]*models.Achievement, error) {
	args := m.Called(ctx, filters)
	achievements, _ := args.Get(0).([]*models.Achievement)
	return achievements, args.Error(1)
}

func (m *MockAchievementRepository) GetFeaturedAchievements(ctx context.Context) ([]*models.Achievement, error) {
	args := m.Called(ctx)
	achievements, _ := args.Get(0).([]*models.Achievement)
	return achievements, args.Error(1)
}

func (m *MockAchievementRepository) CreateAchievement(ctx context.Context, achievement *models.Achievement) error {
	return m.Called(ctx, achievement).Error(0)
}

func (m *MockAchievementRepository) UpdateAchievement(ctx context.Context, achievement *models.Achievement) error {
	return m.Called(ctx, achievement).Error(0)
}

func (m *MockAchievementRepository) DeleteAchievement(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

type MockEducationRepository struct {
	mock.Mock
}

func (m *MockEducationRepository) GetEducation(ctx context.Context, filters repository.EducationFilters) ([]*models.Education, error) {
	args := m.Called(ctx, filters)
	education, _ := args.Get(0).([]*models.Education)
	return education, args.Error(1)
}

func (m *MockEducationRepository) GetEducationByType(ctx context.Context, eduType string) ([]*models.Education, error) {
	args := m.Called(ctx, eduType)
	education, _ := args.Get(0).([]*models.Education)
	return education, args.Error(1)
}

func (m *MockEducationRepository) GetFeaturedEducation(ctx context.Context) ([]*models.Education, error) {
	args := m.Called(ctx)
	education, _ := args.Get(0).([]*models.Education)
	return education, args.Error(1)
}

func (m *MockEducationRepository) CreateEducation(ctx context.Context, education *models.Education) error {
	return m.Called(ctx, education).Error(0)
}

func (m *MockEducationRepository) UpdateEducation(ctx context.Context, education *models.Education) error {
	return m.Called(ctx, education).Error(0)
}

func (m *MockEducationRepository) DeleteEducation(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) GetProjects(ctx context.Context, filters repository.ProjectFilters) ([]*models.Project, error) {
	args := m.Called(ctx, filters)
	projects, _ := args.Get(0).([]*models.Project)
	return projects, args.Error(1)
}

func (m *MockProjectRepository) GetProjectByID(ctx context.Context, id int) (*models.Project, error) {
	args := m.Called(ctx, id)
	project, _ := args.Get(0).(*models.Project)
	return project, args.Error(1)
}

func (m *MockProjectRepository) GetFeaturedProjects(ctx context.Context) ([]*models.Project, error) {
	args := m.Called(ctx)
	projects, _ := args.Get(0).([]*models.Project)
	return projects, args.Error(1)
}

func (m *MockProjectRepository) CreateProject(ctx context.Context, project *models.Project) error {
	return m.Called(ctx, project).Error(0)
}

func (m *MockProjectRepository) UpdateProject(ctx context.Context, project *models.Project) error {
	return m.Called(ctx, project).Error(0)
}

func (m *MockProjectRepository) DeleteProject(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func TestResumeService(t *testing.T) {
	ctx := context.Background()

	t.Run("GetProfile_Success", func(t *testing.T) {
		mockProfileRepo := new(MockProfileRepository)
		mockRepos := repository.Repositories{Profile: mockProfileRepo}
		service := NewResumeService(mockRepos)

		expectedProfile := &models.Profile{ID: 1, Name: "Test User"}
		mockProfileRepo.On("GetProfile", ctx).Return(expectedProfile, nil)

		profile, err := service.GetProfile(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedProfile, profile)
		mockProfileRepo.AssertExpectations(t)
	})

	t.Run("GetProfile_Error", func(t *testing.T) {
		mockProfileRepo := new(MockProfileRepository)
		mockRepos := repository.Repositories{Profile: mockProfileRepo}
		service := NewResumeService(mockRepos)

		expectedError := errors.New("database error")
		mockProfileRepo.On("GetProfile", ctx).Return(nil, expectedError)

		profile, err := service.GetProfile(ctx)

		assert.Error(t, err)
		assert.Nil(t, profile)
		assert.Equal(t, expectedError, err)
		mockProfileRepo.AssertExpectations(t)
	})

	t.Run("GetExperiences_Success", func(t *testing.T) {
		mockExperienceRepo := new(MockExperienceRepository)
		mockRepos := repository.Repositories{Experience: mockExperienceRepo}
		service := NewResumeService(mockRepos)

		filters := repository.ExperienceFilters{Limit: 10}
		expectedExperiences := []*models.Experience{{ID: 1, Company: "Test Co"}}
		mockExperienceRepo.On("GetExperiences", ctx, filters).Return(expectedExperiences, nil)

		experiences, err := service.GetExperiences(ctx, filters)

		assert.NoError(t, err)
		assert.Equal(t, expectedExperiences, experiences)
		mockExperienceRepo.AssertExpectations(t)
	})

	t.Run("GetSkills_Success", func(t *testing.T) {
		mockSkillRepo := new(MockSkillRepository)
		mockRepos := repository.Repositories{Skill: mockSkillRepo}
		service := NewResumeService(mockRepos)

		filters := repository.SkillFilters{Limit: 10}
		expectedSkills := []*models.Skill{{ID: 1, Name: "Go"}}
		mockSkillRepo.On("GetSkills", ctx, filters).Return(expectedSkills, nil)

		skills, err := service.GetSkills(ctx, filters)

		assert.NoError(t, err)
		assert.Equal(t, expectedSkills, skills)
		mockSkillRepo.AssertExpectations(t)
	})

	t.Run("GetAchievements_Success", func(t *testing.T) {
		mockAchievementRepo := new(MockAchievementRepository)
		mockRepos := repository.Repositories{Achievement: mockAchievementRepo}
		service := NewResumeService(mockRepos)

		filters := repository.AchievementFilters{Limit: 10}
		description := "Test Achievement"
		expectedAchievements := []*models.Achievement{{ID: 1, Description: &description}}
		mockAchievementRepo.On("GetAchievements", ctx, filters).Return(expectedAchievements, nil)

		achievements, err := service.GetAchievements(ctx, filters)

		assert.NoError(t, err)
		assert.Equal(t, expectedAchievements, achievements)
		mockAchievementRepo.AssertExpectations(t)
	})

	t.Run("GetEducation_Success", func(t *testing.T) {
		mockEducationRepo := new(MockEducationRepository)
		mockRepos := repository.Repositories{Education: mockEducationRepo}
		service := NewResumeService(mockRepos)

		filters := repository.EducationFilters{Limit: 10}
		expectedEducation := []*models.Education{{ID: 1, Institution: "Test University"}}
		mockEducationRepo.On("GetEducation", ctx, filters).Return(expectedEducation, nil)

		education, err := service.GetEducation(ctx, filters)

		assert.NoError(t, err)
		assert.Equal(t, expectedEducation, education)
		mockEducationRepo.AssertExpectations(t)
	})

	t.Run("GetProjects_Success", func(t *testing.T) {
		mockProjectRepo := new(MockProjectRepository)
		mockRepos := repository.Repositories{Project: mockProjectRepo}
		service := NewResumeService(mockRepos)

		filters := repository.ProjectFilters{Limit: 10}
		expectedProjects := []*models.Project{{ID: 1, Name: "Test Project"}}
		mockProjectRepo.On("GetProjects", ctx, filters).Return(expectedProjects, nil)

		projects, err := service.GetProjects(ctx, filters)

		assert.NoError(t, err)
		assert.Equal(t, expectedProjects, projects)
		mockProjectRepo.AssertExpectations(t)
	})
}
