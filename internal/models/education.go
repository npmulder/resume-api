package models

import (
	"time"
)

// Education represents education and certifications
type Education struct {
	ID                    int        `json:"id" db:"id"`
	Institution           string     `json:"institution" db:"institution"`
	DegreeOrCertification string     `json:"degree_or_certification" db:"degree_or_certification"`
	FieldOfStudy          *string    `json:"field_of_study,omitempty" db:"field_of_study"`
	YearCompleted         *int       `json:"year_completed,omitempty" db:"year_completed"`
	YearStarted           *int       `json:"year_started,omitempty" db:"year_started"`
	Description           *string    `json:"description,omitempty" db:"description"`
	Type                  string     `json:"type" db:"type"` // education or certification
	Status                string     `json:"status" db:"status"` // completed, in_progress, planned
	CredentialID          *string    `json:"credential_id,omitempty" db:"credential_id"`
	CredentialURL         *string    `json:"credential_url,omitempty" db:"credential_url"`
	ExpiryDate            *time.Time `json:"expiry_date,omitempty" db:"expiry_date"`
	OrderIndex            int        `json:"order_index" db:"order_index"`
	IsFeatured            bool       `json:"is_featured" db:"is_featured"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`

	// Computed fields for compatibility with interface
	DegreeTitle  string     `json:"degree_title" db:"-"`
	StartDate    *time.Time `json:"start_date,omitempty" db:"-"`
	EndDate      *time.Time `json:"end_date,omitempty" db:"-"`
	Grade        *string    `json:"grade,omitempty" db:"-"`
}

// Education type constants
const (
	EducationTypeEducation     = "education"
	EducationTypeCertification = "certification"
)

// Education status constants
const (
	EducationStatusCompleted  = "completed"
	EducationStatusInProgress = "in_progress"
	EducationStatusPlanned    = "planned"
)

// ValidEducationTypes returns valid education types
func ValidEducationTypes() []string {
	return []string{EducationTypeEducation, EducationTypeCertification}
}

// ValidEducationStatuses returns valid education statuses
func ValidEducationStatuses() []string {
	return []string{
		EducationStatusCompleted,
		EducationStatusInProgress,
		EducationStatusPlanned,
	}
}