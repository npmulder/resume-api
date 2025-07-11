package models

import (
	"time"
)

// Experience represents work history and professional experience
type Experience struct {
	ID          int              `json:"id" db:"id"`
	Company     string           `json:"company" db:"company"`
	Position    string           `json:"position" db:"position"`
	StartDate   time.Time        `json:"start_date" db:"start_date"`
	EndDate     *time.Time       `json:"end_date,omitempty" db:"end_date"`
	Description *string          `json:"description,omitempty" db:"description"`
	Highlights  []string         `json:"highlights,omitempty" db:"highlights"`
	OrderIndex  int              `json:"order_index" db:"order_index"`
	IsCurrent   bool             `json:"is_current" db:"-"` // Computed field based on end_date
	Location    *string          `json:"location,omitempty" db:"location"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

// IsCurrentPosition returns true if this is a current position (end_date is nil)
func (e *Experience) IsCurrentPosition() bool {
	return e.EndDate == nil
}