package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/repository"
)

// EducationRepository implements repository.EducationRepository for PostgreSQL
type EducationRepository struct {
	db *pgxpool.Pool
}

// NewEducationRepository creates a new PostgreSQL education repository
func NewEducationRepository(db *pgxpool.Pool) *EducationRepository {
	return &EducationRepository{db: db}
}

// GetEducation retrieves all education entries with optional filtering
func (r *EducationRepository) GetEducation(ctx context.Context, filters repository.EducationFilters) ([]*models.Education, error) {
	query := `
		SELECT id, institution, degree_or_certification, field_of_study, year_completed, 
		       year_started, description, type, status, credential_id, credential_url, 
		       expiry_date, order_index, is_featured, created_at, updated_at
		FROM education`
	
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Apply filters
	if filters.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, filters.Type)
		argIndex++
	}

	if filters.Institution != "" {
		conditions = append(conditions, fmt.Sprintf("institution ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Institution+"%")
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.Featured != nil {
		conditions = append(conditions, fmt.Sprintf("is_featured = $%d", argIndex))
		args = append(args, *filters.Featured)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY type, year_completed DESC, order_index"

	// Apply pagination
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, repository.NewRepositoryError("get", "education", err)
	}
	defer rows.Close()

	var educations []*models.Education
	for rows.Next() {
		var edu models.Education
		err := rows.Scan(
			&edu.ID,
			&edu.Institution,
			&edu.DegreeOrCertification,
			&edu.FieldOfStudy,
			&edu.YearCompleted,
			&edu.YearStarted,
			&edu.Description,
			&edu.Type,
			&edu.Status,
			&edu.CredentialID,
			&edu.CredentialURL,
			&edu.ExpiryDate,
			&edu.OrderIndex,
			&edu.IsFeatured,
			&edu.CreatedAt,
			&edu.UpdatedAt,
		)
		if err != nil {
			return nil, repository.NewRepositoryError("scan", "education", err)
		}
		educations = append(educations, &edu)
	}

	if err := rows.Err(); err != nil {
		return nil, repository.NewRepositoryError("iterate", "education", err)
	}

	return educations, nil
}

// GetEducationByType retrieves education entries by type (education, certification)
func (r *EducationRepository) GetEducationByType(ctx context.Context, eduType string) ([]*models.Education, error) {
	filters := repository.EducationFilters{
		Type: eduType,
	}
	return r.GetEducation(ctx, filters)
}

// GetFeaturedEducation retrieves only featured education entries
func (r *EducationRepository) GetFeaturedEducation(ctx context.Context) ([]*models.Education, error) {
	featured := true
	filters := repository.EducationFilters{
		Featured: &featured,
	}
	return r.GetEducation(ctx, filters)
}

// CreateEducation creates a new education entry
func (r *EducationRepository) CreateEducation(ctx context.Context, education *models.Education) error {
	query := `
		INSERT INTO education (institution, degree_or_certification, field_of_study, 
		                      year_completed, year_started, description, type, status, 
		                      credential_id, credential_url, expiry_date, order_index, is_featured)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		education.Institution,
		education.DegreeOrCertification,
		education.FieldOfStudy,
		education.YearCompleted,
		education.YearStarted,
		education.Description,
		education.Type,
		education.Status,
		education.CredentialID,
		education.CredentialURL,
		education.ExpiryDate,
		education.OrderIndex,
		education.IsFeatured,
	).Scan(&education.ID, &education.CreatedAt, &education.UpdatedAt)

	if err != nil {
		return repository.NewRepositoryError("create", "education", err)
	}

	return nil
}

// UpdateEducation updates an existing education entry
func (r *EducationRepository) UpdateEducation(ctx context.Context, education *models.Education) error {
	query := `
		UPDATE education 
		SET institution = $2, degree_or_certification = $3, field_of_study = $4, 
		    year_completed = $5, year_started = $6, description = $7, type = $8, 
		    status = $9, credential_id = $10, credential_url = $11, expiry_date = $12, 
		    order_index = $13, is_featured = $14, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		education.ID,
		education.Institution,
		education.DegreeOrCertification,
		education.FieldOfStudy,
		education.YearCompleted,
		education.YearStarted,
		education.Description,
		education.Type,
		education.Status,
		education.CredentialID,
		education.CredentialURL,
		education.ExpiryDate,
		education.OrderIndex,
		education.IsFeatured,
	).Scan(&education.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return repository.NewRepositoryError("update", "education", fmt.Errorf("education with id %d not found", education.ID))
		}
		return repository.NewRepositoryError("update", "education", err)
	}

	return nil
}

// DeleteEducation deletes an education entry by ID
func (r *EducationRepository) DeleteEducation(ctx context.Context, id int) error {
	query := `DELETE FROM education WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return repository.NewRepositoryError("delete", "education", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return repository.NewRepositoryError("delete", "education", fmt.Errorf("education with id %d not found", id))
	}

	return nil
}