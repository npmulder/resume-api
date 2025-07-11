package models

import (
	"time"
)

// Achievement represents key accomplishments and achievements
type Achievement struct {
	ID           int       `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Description  *string   `json:"description,omitempty" db:"description"`
	Category     *string   `json:"category,omitempty" db:"category"`
	ImpactMetric *string   `json:"impact_metric,omitempty" db:"impact_metric"`
	YearAchieved *int      `json:"year_achieved,omitempty" db:"year_achieved"`
	OrderIndex   int       `json:"order_index" db:"order_index"`
	IsFeatured   bool      `json:"is_featured" db:"is_featured"`
	DateAchieved *time.Time `json:"date_achieved,omitempty" db:"-"` // For interface compatibility
	Organization *string   `json:"organization,omitempty" db:"-"`  // For interface compatibility
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Achievement category constants
const (
	AchievementCategoryPerformance = "performance"
	AchievementCategorySecurity    = "security"
	AchievementCategoryLeadership  = "leadership"
	AchievementCategoryInnovation  = "innovation"
	AchievementCategoryEfficiency  = "efficiency"
	AchievementCategoryTeamwork    = "teamwork"
)