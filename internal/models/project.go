package models

import (
	"time"
)

// Project represents notable projects and implementations
type Project struct {
	ID               int       `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Description      *string   `json:"description,omitempty" db:"description"`
	ShortDescription *string   `json:"short_description,omitempty" db:"short_description"`
	Technologies     []string  `json:"technologies,omitempty" db:"technologies"` // JSONB in DB
	GitHubURL        *string   `json:"github_url,omitempty" db:"github_url"`
	DemoURL          *string   `json:"demo_url,omitempty" db:"demo_url"`
	StartDate        *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate          *time.Time `json:"end_date,omitempty" db:"end_date"`
	Status           string    `json:"status" db:"status"` // active, completed, archived, planned
	IsFeatured       bool      `json:"is_featured" db:"is_featured"`
	OrderIndex       int       `json:"order_index" db:"order_index"`
	KeyFeatures      []string  `json:"key_features,omitempty" db:"key_features"` // TEXT[] in DB
	Highlights       []string  `json:"highlights,omitempty" db:"-"` // For interface compatibility
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Project status constants
const (
	ProjectStatusActive    = "active"
	ProjectStatusCompleted = "completed"
	ProjectStatusArchived  = "archived"
	ProjectStatusPlanned   = "planned"
)

// ValidProjectStatuses returns valid project statuses
func ValidProjectStatuses() []string {
	return []string{
		ProjectStatusActive,
		ProjectStatusCompleted,
		ProjectStatusArchived,
		ProjectStatusPlanned,
	}
}

// IsOngoing returns true if the project is currently active (end_date is nil and status is active)
func (p *Project) IsOngoing() bool {
	return p.EndDate == nil && p.Status == ProjectStatusActive
}