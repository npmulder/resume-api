// Package services defines interfaces for the business logic layer
package services

import (
	"context"

	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/repository"
)

// ResumeService defines the business logic for handling resume data.
// It orchestrates calls to the repository layer and implements business rules.
type ResumeService interface {
	GetProfile(ctx context.Context) (*models.Profile, error)
	GetExperiences(ctx context.Context, filters repository.ExperienceFilters) ([]*models.Experience, error)
	GetSkills(ctx context.Context, filters repository.SkillFilters) ([]*models.Skill, error)
	GetAchievements(ctx context.Context, filters repository.AchievementFilters) ([]*models.Achievement, error)
	GetEducation(ctx context.Context, filters repository.EducationFilters) ([]*models.Education, error)
	GetProjects(ctx context.Context, filters repository.ProjectFilters) ([]*models.Project, error)
}
