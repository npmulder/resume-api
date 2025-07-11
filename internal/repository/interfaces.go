// Package repository defines data access interfaces and implementations
package repository

import (
	"context"
	"fmt"

	"github.com/npmulder/resume-api/internal/models"
)

// ProfileRepository defines operations for profile data
type ProfileRepository interface {
	// GetProfile retrieves the user's profile information
	GetProfile(ctx context.Context) (*models.Profile, error)
	
	// UpdateProfile updates the user's profile information
	UpdateProfile(ctx context.Context, profile *models.Profile) error
	
	// CreateProfile creates a new profile (typically only used once)
	CreateProfile(ctx context.Context, profile *models.Profile) error
}

// ExperienceRepository defines operations for work experience data
type ExperienceRepository interface {
	// GetExperiences retrieves all work experiences with optional filtering
	GetExperiences(ctx context.Context, filters ExperienceFilters) ([]*models.Experience, error)
	
	// GetExperienceByID retrieves a specific experience by ID
	GetExperienceByID(ctx context.Context, id int) (*models.Experience, error)
	
	// CreateExperience creates a new experience entry
	CreateExperience(ctx context.Context, experience *models.Experience) error
	
	// UpdateExperience updates an existing experience
	UpdateExperience(ctx context.Context, experience *models.Experience) error
	
	// DeleteExperience deletes an experience by ID
	DeleteExperience(ctx context.Context, id int) error
}

// SkillRepository defines operations for skills data
type SkillRepository interface {
	// GetSkills retrieves all skills with optional filtering
	GetSkills(ctx context.Context, filters SkillFilters) ([]*models.Skill, error)
	
	// GetSkillsByCategory retrieves skills grouped by category
	GetSkillsByCategory(ctx context.Context, category string) ([]*models.Skill, error)
	
	// GetFeaturedSkills retrieves only featured skills
	GetFeaturedSkills(ctx context.Context) ([]*models.Skill, error)
	
	// CreateSkill creates a new skill entry
	CreateSkill(ctx context.Context, skill *models.Skill) error
	
	// UpdateSkill updates an existing skill
	UpdateSkill(ctx context.Context, skill *models.Skill) error
	
	// DeleteSkill deletes a skill by ID
	DeleteSkill(ctx context.Context, id int) error
}

// AchievementRepository defines operations for achievements data
type AchievementRepository interface {
	// GetAchievements retrieves all achievements with optional filtering
	GetAchievements(ctx context.Context, filters AchievementFilters) ([]*models.Achievement, error)
	
	// GetFeaturedAchievements retrieves only featured achievements
	GetFeaturedAchievements(ctx context.Context) ([]*models.Achievement, error)
	
	// CreateAchievement creates a new achievement entry
	CreateAchievement(ctx context.Context, achievement *models.Achievement) error
	
	// UpdateAchievement updates an existing achievement
	UpdateAchievement(ctx context.Context, achievement *models.Achievement) error
	
	// DeleteAchievement deletes an achievement by ID
	DeleteAchievement(ctx context.Context, id int) error
}

// EducationRepository defines operations for education and certification data
type EducationRepository interface {
	// GetEducation retrieves all education entries with optional filtering
	GetEducation(ctx context.Context, filters EducationFilters) ([]*models.Education, error)
	
	// GetEducationByType retrieves education entries by type (education, certification)
	GetEducationByType(ctx context.Context, eduType string) ([]*models.Education, error)
	
	// GetFeaturedEducation retrieves only featured education entries
	GetFeaturedEducation(ctx context.Context) ([]*models.Education, error)
	
	// CreateEducation creates a new education entry
	CreateEducation(ctx context.Context, education *models.Education) error
	
	// UpdateEducation updates an existing education entry
	UpdateEducation(ctx context.Context, education *models.Education) error
	
	// DeleteEducation deletes an education entry by ID
	DeleteEducation(ctx context.Context, id int) error
}

// ProjectRepository defines operations for project data
type ProjectRepository interface {
	// GetProjects retrieves all projects with optional filtering
	GetProjects(ctx context.Context, filters ProjectFilters) ([]*models.Project, error)
	
	// GetProjectByID retrieves a specific project by ID
	GetProjectByID(ctx context.Context, id int) (*models.Project, error)
	
	// GetFeaturedProjects retrieves only featured projects
	GetFeaturedProjects(ctx context.Context) ([]*models.Project, error)
	
	// CreateProject creates a new project entry
	CreateProject(ctx context.Context, project *models.Project) error
	
	// UpdateProject updates an existing project
	UpdateProject(ctx context.Context, project *models.Project) error
	
	// DeleteProject deletes a project by ID
	DeleteProject(ctx context.Context, id int) error
}

// Filter types for repository queries

// ExperienceFilters defines filtering options for experience queries
type ExperienceFilters struct {
	Company    string
	Position   string
	DateFrom   *string // ISO date string
	DateTo     *string // ISO date string
	IsCurrent  *bool   // Filter for current positions (end_date IS NULL)
	Limit      int
	Offset     int
}

// SkillFilters defines filtering options for skill queries
type SkillFilters struct {
	Category string
	Level    string
	Featured *bool
	Limit    int
	Offset   int
}

// AchievementFilters defines filtering options for achievement queries
type AchievementFilters struct {
	Category string
	Year     *int
	Featured *bool
	Limit    int
	Offset   int
}

// EducationFilters defines filtering options for education queries
type EducationFilters struct {
	Type         string // 'education' or 'certification'
	Institution  string
	Status       string // 'completed', 'in_progress', 'planned'
	Featured     *bool
	Limit        int
	Offset       int
}

// ProjectFilters defines filtering options for project queries
type ProjectFilters struct {
	Status       string // 'active', 'completed', 'archived', 'planned'
	Technology   string // Search in technologies JSONB
	Featured     *bool
	Limit        int
	Offset       int
}

// Repositories aggregates all repository interfaces
type Repositories struct {
	Profile     ProfileRepository
	Experience  ExperienceRepository
	Skill       SkillRepository
	Achievement AchievementRepository
	Education   EducationRepository
	Project     ProjectRepository
}

// RepositoryError represents a repository-specific error
type RepositoryError struct {
	Operation string
	Entity    string
	Err       error
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("repository error during %s on %s: %v", e.Operation, e.Entity, e.Err)
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// NewRepositoryError creates a new repository error
func NewRepositoryError(operation, entity string, err error) *RepositoryError {
	return &RepositoryError{
		Operation: operation,
		Entity:    entity,
		Err:       err,
	}
}