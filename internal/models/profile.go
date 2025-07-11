package models

import (
	"time"
)

// Profile represents the user's personal information and summary
type Profile struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Title     string    `json:"title" db:"title"`
	Email     string    `json:"email" db:"email"`
	Phone     *string   `json:"phone,omitempty" db:"phone"`
	Location  *string   `json:"location,omitempty" db:"location"`
	LinkedIn  *string   `json:"linkedin,omitempty" db:"linkedin"`
	GitHub    *string   `json:"github,omitempty" db:"github"`
	Summary   *string   `json:"summary,omitempty" db:"summary"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}