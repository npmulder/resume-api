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

// ExperienceRepository implements repository.ExperienceRepository for PostgreSQL
type ExperienceRepository struct {
	db *pgxpool.Pool
}

// NewExperienceRepository creates a new PostgreSQL experience repository
func NewExperienceRepository(db *pgxpool.Pool) *ExperienceRepository {
	return &ExperienceRepository{db: db}
}

// GetExperiences retrieves all work experiences with optional filtering
func (r *ExperienceRepository) GetExperiences(ctx context.Context, filters repository.ExperienceFilters) ([]*models.Experience, error) {
	query := `
		SELECT id, company, position, start_date, end_date, description, 
		       highlights, order_index, created_at, updated_at
		FROM experiences`
	
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Apply filters
	if filters.Company != "" {
		conditions = append(conditions, fmt.Sprintf("company ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Company+"%")
		argIndex++
	}

	if filters.Position != "" {
		conditions = append(conditions, fmt.Sprintf("position ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Position+"%")
		argIndex++
	}

	if filters.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", argIndex))
		args = append(args, *filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("start_date <= $%d", argIndex))
		args = append(args, *filters.DateTo)
		argIndex++
	}

	if filters.IsCurrent != nil {
		if *filters.IsCurrent {
			conditions = append(conditions, "end_date IS NULL")
		} else {
			conditions = append(conditions, "end_date IS NOT NULL")
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY start_date DESC"

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
		return nil, repository.NewRepositoryError("get", "experiences", err)
	}
	defer rows.Close()

	var experiences []*models.Experience
	for rows.Next() {
		var exp models.Experience
		err := rows.Scan(
			&exp.ID,
			&exp.Company,
			&exp.Position,
			&exp.StartDate,
			&exp.EndDate,
			&exp.Description,
			&exp.Highlights,
			&exp.OrderIndex,
			&exp.CreatedAt,
			&exp.UpdatedAt,
		)
		if err != nil {
			return nil, repository.NewRepositoryError("scan", "experience", err)
		}
		experiences = append(experiences, &exp)
	}

	if err := rows.Err(); err != nil {
		return nil, repository.NewRepositoryError("iterate", "experiences", err)
	}

	return experiences, nil
}

// GetExperienceByID retrieves a specific experience by ID
func (r *ExperienceRepository) GetExperienceByID(ctx context.Context, id int) (*models.Experience, error) {
	query := `
		SELECT id, company, position, start_date, end_date, description, 
		       highlights, order_index, created_at, updated_at
		FROM experiences 
		WHERE id = $1`

	var exp models.Experience
	err := r.db.QueryRow(ctx, query, id).Scan(
		&exp.ID,
		&exp.Company,
		&exp.Position,
		&exp.StartDate,
		&exp.EndDate,
		&exp.Description,
		&exp.Highlights,
		&exp.OrderIndex,
		&exp.CreatedAt,
		&exp.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.NewRepositoryError("get", "experience", fmt.Errorf("experience with id %d not found", id))
		}
		return nil, repository.NewRepositoryError("get", "experience", err)
	}

	return &exp, nil
}

// CreateExperience creates a new experience entry
func (r *ExperienceRepository) CreateExperience(ctx context.Context, experience *models.Experience) error {
	query := `
		INSERT INTO experiences (company, position, start_date, end_date, description, 
		                        highlights, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		experience.Company,
		experience.Position,
		experience.StartDate,
		experience.EndDate,
		experience.Description,
		experience.Highlights,
		experience.OrderIndex,
	).Scan(&experience.ID, &experience.CreatedAt, &experience.UpdatedAt)

	if err != nil {
		return repository.NewRepositoryError("create", "experience", err)
	}

	return nil
}

// UpdateExperience updates an existing experience
func (r *ExperienceRepository) UpdateExperience(ctx context.Context, experience *models.Experience) error {
	query := `
		UPDATE experiences 
		SET company = $2, position = $3, start_date = $4, end_date = $5, 
		    description = $6, highlights = $7, order_index = $8,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		experience.ID,
		experience.Company,
		experience.Position,
		experience.StartDate,
		experience.EndDate,
		experience.Description,
		experience.Highlights,
		experience.OrderIndex,
	).Scan(&experience.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return repository.NewRepositoryError("update", "experience", fmt.Errorf("experience with id %d not found", experience.ID))
		}
		return repository.NewRepositoryError("update", "experience", err)
	}

	return nil
}

// DeleteExperience deletes an experience by ID
func (r *ExperienceRepository) DeleteExperience(ctx context.Context, id int) error {
	query := `DELETE FROM experiences WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return repository.NewRepositoryError("delete", "experience", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return repository.NewRepositoryError("delete", "experience", fmt.Errorf("experience with id %d not found", id))
	}

	return nil
}