// Package postgres provides PostgreSQL implementations of repository interfaces
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/repository"
)

// ProfileRepository implements repository.ProfileRepository for PostgreSQL
type ProfileRepository struct {
	db *pgxpool.Pool
}

// NewProfileRepository creates a new PostgreSQL profile repository
func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// GetProfile retrieves the user's profile information
func (r *ProfileRepository) GetProfile(ctx context.Context) (*models.Profile, error) {
	query := `
		SELECT id, name, title, email, phone, location, linkedin, github, 
		       summary, created_at, updated_at
		FROM profiles 
		ORDER BY created_at DESC 
		LIMIT 1`

	var profile models.Profile
	err := r.db.QueryRow(ctx, query).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Title,
		&profile.Email,
		&profile.Phone,
		&profile.Location,
		&profile.LinkedIn,
		&profile.GitHub,
		&profile.Summary,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, repository.NewRepositoryError("get", "profile", err)
	}

	return &profile, nil
}

// CreateProfile creates a new profile (typically only used once)
func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *models.Profile) error {
	query := `
		INSERT INTO profiles (name, title, email, phone, location, linkedin, 
		                     github, summary)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		profile.Name,
		profile.Title,
		profile.Email,
		profile.Phone,
		profile.Location,
		profile.LinkedIn,
		profile.GitHub,
		profile.Summary,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		return repository.NewRepositoryError("create", "profile", err)
	}

	return nil
}

// UpdateProfile updates the user's profile information
func (r *ProfileRepository) UpdateProfile(ctx context.Context, profile *models.Profile) error {
	query := `
		UPDATE profiles 
		SET name = $2, title = $3, email = $4, phone = $5, location = $6, 
		    linkedin = $7, github = $8, summary = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		profile.ID,
		profile.Name,
		profile.Title,
		profile.Email,
		profile.Phone,
		profile.Location,
		profile.LinkedIn,
		profile.GitHub,
		profile.Summary,
	).Scan(&profile.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return repository.NewRepositoryError("update", "profile", err)
	}

	return nil
}
