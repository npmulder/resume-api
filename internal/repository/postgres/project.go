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

// ProjectRepository implements repository.ProjectRepository for PostgreSQL
type ProjectRepository struct {
	db *pgxpool.Pool
}

// NewProjectRepository creates a new PostgreSQL project repository
func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// GetProjects retrieves all projects with optional filtering
func (r *ProjectRepository) GetProjects(ctx context.Context, filters repository.ProjectFilters) ([]*models.Project, error) {
	query := `
		SELECT id, name, description, short_description, technologies, github_url, 
		       demo_url, start_date, end_date, status, is_featured, order_index, 
		       key_features, created_at, updated_at
		FROM projects`
	
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Apply filters
	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.Technology != "" {
		// Search in JSONB technologies array
		conditions = append(conditions, fmt.Sprintf("technologies ? $%d", argIndex))
		args = append(args, filters.Technology)
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

	query += " ORDER BY start_date DESC, order_index"

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
		return nil, repository.NewRepositoryError("get", "projects", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.ShortDescription,
			&project.Technologies,
			&project.GitHubURL,
			&project.DemoURL,
			&project.StartDate,
			&project.EndDate,
			&project.Status,
			&project.IsFeatured,
			&project.OrderIndex,
			&project.KeyFeatures,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, repository.NewRepositoryError("scan", "project", err)
		}
		projects = append(projects, &project)
	}

	if err := rows.Err(); err != nil {
		return nil, repository.NewRepositoryError("iterate", "projects", err)
	}

	return projects, nil
}

// GetProjectByID retrieves a specific project by ID
func (r *ProjectRepository) GetProjectByID(ctx context.Context, id int) (*models.Project, error) {
	query := `
		SELECT id, name, description, short_description, technologies, github_url, 
		       demo_url, start_date, end_date, status, is_featured, order_index, 
		       key_features, created_at, updated_at
		FROM projects 
		WHERE id = $1`

	var project models.Project
	err := r.db.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.ShortDescription,
		&project.Technologies,
		&project.GitHubURL,
		&project.DemoURL,
		&project.StartDate,
		&project.EndDate,
		&project.Status,
		&project.IsFeatured,
		&project.OrderIndex,
		&project.KeyFeatures,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.NewRepositoryError("get", "project", fmt.Errorf("project with id %d not found", id))
		}
		return nil, repository.NewRepositoryError("get", "project", err)
	}

	return &project, nil
}

// GetFeaturedProjects retrieves only featured projects
func (r *ProjectRepository) GetFeaturedProjects(ctx context.Context) ([]*models.Project, error) {
	featured := true
	filters := repository.ProjectFilters{
		Featured: &featured,
	}
	return r.GetProjects(ctx, filters)
}

// CreateProject creates a new project entry
func (r *ProjectRepository) CreateProject(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO projects (name, description, short_description, technologies, 
		                     github_url, demo_url, start_date, end_date, status, 
		                     is_featured, order_index, key_features)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		project.Name,
		project.Description,
		project.ShortDescription,
		project.Technologies,
		project.GitHubURL,
		project.DemoURL,
		project.StartDate,
		project.EndDate,
		project.Status,
		project.IsFeatured,
		project.OrderIndex,
		project.KeyFeatures,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return repository.NewRepositoryError("create", "project", err)
	}

	return nil
}

// UpdateProject updates an existing project
func (r *ProjectRepository) UpdateProject(ctx context.Context, project *models.Project) error {
	query := `
		UPDATE projects 
		SET name = $2, description = $3, short_description = $4, technologies = $5, 
		    github_url = $6, demo_url = $7, start_date = $8, end_date = $9, 
		    status = $10, is_featured = $11, order_index = $12, key_features = $13,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		project.ID,
		project.Name,
		project.Description,
		project.ShortDescription,
		project.Technologies,
		project.GitHubURL,
		project.DemoURL,
		project.StartDate,
		project.EndDate,
		project.Status,
		project.IsFeatured,
		project.OrderIndex,
		project.KeyFeatures,
	).Scan(&project.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return repository.NewRepositoryError("update", "project", fmt.Errorf("project with id %d not found", project.ID))
		}
		return repository.NewRepositoryError("update", "project", err)
	}

	return nil
}

// DeleteProject deletes a project by ID
func (r *ProjectRepository) DeleteProject(ctx context.Context, id int) error {
	query := `DELETE FROM projects WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return repository.NewRepositoryError("delete", "project", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return repository.NewRepositoryError("delete", "project", fmt.Errorf("project with id %d not found", id))
	}

	return nil
}