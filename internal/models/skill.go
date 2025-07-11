package models

import (
	"time"
)

// Skill represents technical and soft skills with categories and levels
type Skill struct {
	ID              int       `json:"id" db:"id"`
	Category        string    `json:"category" db:"category"`
	Name            string    `json:"name" db:"name"`
	Level           *string   `json:"level,omitempty" db:"level"` // beginner, intermediate, advanced, expert
	YearsExperience *int      `json:"years_experience,omitempty" db:"years_experience"`
	OrderIndex      int       `json:"order_index" db:"order_index"`
	IsFeatured      bool      `json:"is_featured" db:"is_featured"`
	Description     *string   `json:"description,omitempty" db:"description"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// SkillLevel constants for validation
const (
	SkillLevelBeginner     = "beginner"
	SkillLevelIntermediate = "intermediate"
	SkillLevelAdvanced     = "advanced"
	SkillLevelExpert       = "expert"
)

// ValidSkillLevels returns a slice of valid skill levels
func ValidSkillLevels() []string {
	return []string{
		SkillLevelBeginner,
		SkillLevelIntermediate,
		SkillLevelAdvanced,
		SkillLevelExpert,
	}
}