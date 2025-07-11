// Package services implements the business logic for the Resume API.
package services

import (
	"context"

	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/repository"
)

// resumeService is the implementation of the ResumeService interface.
// It uses the repository interfaces to access the data layer.
type resumeService struct {
	repos repository.Repositories
}

// NewResumeService creates a new instance of the resumeService.
// It takes the repository interfaces as dependencies.
func NewResumeService(repos repository.Repositories) ResumeService {
	return &resumeService{
		repos: repos,
	}
}

// GetProfile retrieves the user's profile.
func (s *resumeService) GetProfile(ctx context.Context) (*models.Profile, error) {
	return s.repos.Profile.GetProfile(ctx)
}

// GetExperiences retrieves work experiences with optional filtering.
func (s *resumeService) GetExperiences(ctx context.Context, filters repository.ExperienceFilters) ([]*models.Experience, error) {
	return s.repos.Experience.GetExperiences(ctx, filters)
}

// GetSkills retrieves skills with optional filtering.
func (s *resumeService) GetSkills(ctx context.Context, filters repository.SkillFilters) ([]*models.Skill, error) {
	return s.repos.Skill.GetSkills(ctx, filters)
}

// GetAchievements retrieves achievements with optional filtering.
func (s *resumeService) GetAchievements(ctx context.Context, filters repository.AchievementFilters) ([]*models.Achievement, error) {
	return s.repos.Achievement.GetAchievements(ctx, filters)
}

// GetEducation retrieves education entries with optional filtering.
func (s *resumeService) GetEducation(ctx context.Context, filters repository.EducationFilters) ([]*models.Education, error) {
	return s.repos.Education.GetEducation(ctx, filters)
}

// GetProjects retrieves projects with optional filtering.
func (s *resumeService) GetProjects(ctx context.Context, filters repository.ProjectFilters) ([]*models.Project, error) {
	return s.repos.Project.GetProjects(ctx, filters)
}
