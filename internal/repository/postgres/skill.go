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

// SkillRepository implements repository.SkillRepository for PostgreSQL
type SkillRepository struct {
	db *pgxpool.Pool
}

// NewSkillRepository creates a new PostgreSQL skill repository
func NewSkillRepository(db *pgxpool.Pool) *SkillRepository {
	return &SkillRepository{db: db}
}

// GetSkills retrieves all skills with optional filtering
func (r *SkillRepository) GetSkills(ctx context.Context, filters repository.SkillFilters) ([]*models.Skill, error) {
	query := `
		SELECT id, category, name, level, years_experience, order_index, is_featured, 
		       created_at, updated_at
		FROM skills`
	
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Apply filters
	if filters.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, filters.Category)
		argIndex++
	}

	if filters.Level != "" {
		conditions = append(conditions, fmt.Sprintf("level = $%d", argIndex))
		args = append(args, filters.Level)
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

	query += " ORDER BY category, order_index, name"

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
		return nil, repository.NewRepositoryError("get", "skills", err)
	}
	defer rows.Close()

	var skills []*models.Skill
	for rows.Next() {
		var skill models.Skill
		err := rows.Scan(
			&skill.ID,
			&skill.Category,
			&skill.Name,
			&skill.Level,
			&skill.YearsExperience,
			&skill.OrderIndex,
			&skill.IsFeatured,
			&skill.CreatedAt,
			&skill.UpdatedAt,
		)
		if err != nil {
			return nil, repository.NewRepositoryError("scan", "skill", err)
		}
		skills = append(skills, &skill)
	}

	if err := rows.Err(); err != nil {
		return nil, repository.NewRepositoryError("iterate", "skills", err)
	}

	return skills, nil
}

// GetSkillsByCategory retrieves skills grouped by category
func (r *SkillRepository) GetSkillsByCategory(ctx context.Context, category string) ([]*models.Skill, error) {
	filters := repository.SkillFilters{
		Category: category,
	}
	return r.GetSkills(ctx, filters)
}

// GetFeaturedSkills retrieves only featured skills
func (r *SkillRepository) GetFeaturedSkills(ctx context.Context) ([]*models.Skill, error) {
	featured := true
	filters := repository.SkillFilters{
		Featured: &featured,
	}
	return r.GetSkills(ctx, filters)
}

// CreateSkill creates a new skill entry
func (r *SkillRepository) CreateSkill(ctx context.Context, skill *models.Skill) error {
	query := `
		INSERT INTO skills (category, name, level, years_experience, order_index, is_featured)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		skill.Category,
		skill.Name,
		skill.Level,
		skill.YearsExperience,
		skill.OrderIndex,
		skill.IsFeatured,
	).Scan(&skill.ID, &skill.CreatedAt, &skill.UpdatedAt)

	if err != nil {
		return repository.NewRepositoryError("create", "skill", err)
	}

	return nil
}

// UpdateSkill updates an existing skill
func (r *SkillRepository) UpdateSkill(ctx context.Context, skill *models.Skill) error {
	query := `
		UPDATE skills 
		SET category = $2, name = $3, level = $4, years_experience = $5, 
		    order_index = $6, is_featured = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		skill.ID,
		skill.Category,
		skill.Name,
		skill.Level,
		skill.YearsExperience,
		skill.OrderIndex,
		skill.IsFeatured,
	).Scan(&skill.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return repository.NewRepositoryError("update", "skill", fmt.Errorf("skill with id %d not found", skill.ID))
		}
		return repository.NewRepositoryError("update", "skill", err)
	}

	return nil
}

// DeleteSkill deletes a skill by ID
func (r *SkillRepository) DeleteSkill(ctx context.Context, id int) error {
	query := `DELETE FROM skills WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return repository.NewRepositoryError("delete", "skill", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return repository.NewRepositoryError("delete", "skill", fmt.Errorf("skill with id %d not found", id))
	}

	return nil
}