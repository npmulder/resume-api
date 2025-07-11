// Package postgres provides PostgreSQL implementations of repository interfaces
package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/npmulder/resume-api/internal/repository"
)

// Repositories implements repository.Repositories with PostgreSQL
type Repositories struct {
	Profile     repository.ProfileRepository
	Experience  repository.ExperienceRepository
	Skill       repository.SkillRepository
	Achievement repository.AchievementRepository
	Education   repository.EducationRepository
	Project     repository.ProjectRepository
}

// NewRepositories creates a new set of PostgreSQL repositories
func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		Profile:     NewProfileRepository(db),
		Experience:  NewExperienceRepository(db),
		Skill:       NewSkillRepository(db),
		Achievement: NewAchievementRepository(db),
		Education:   NewEducationRepository(db),
		Project:     NewProjectRepository(db),
	}
}

// Close gracefully shuts down all repositories
// Currently no cleanup needed for PostgreSQL repositories
func (r *Repositories) Close() error {
	// No cleanup needed for PostgreSQL repositories as they use the shared pool
	return nil
}