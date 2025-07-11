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

// AchievementRepository implements repository.AchievementRepository for PostgreSQL
type AchievementRepository struct {
	db *pgxpool.Pool
}

// NewAchievementRepository creates a new PostgreSQL achievement repository
func NewAchievementRepository(db *pgxpool.Pool) *AchievementRepository {
	return &AchievementRepository{db: db}
}

// GetAchievements retrieves all achievements with optional filtering
func (r *AchievementRepository) GetAchievements(ctx context.Context, filters repository.AchievementFilters) ([]*models.Achievement, error) {
	query := `
		SELECT id, title, description, category, impact_metric, year_achieved, 
		       order_index, is_featured, created_at, updated_at
		FROM achievements`
	
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Apply filters
	if filters.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, filters.Category)
		argIndex++
	}

	if filters.Year != nil {
		conditions = append(conditions, fmt.Sprintf("year_achieved = $%d", argIndex))
		args = append(args, *filters.Year)
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

	query += " ORDER BY year_achieved DESC, order_index"

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
		return nil, repository.NewRepositoryError("get", "achievements", err)
	}
	defer rows.Close()

	var achievements []*models.Achievement
	for rows.Next() {
		var achievement models.Achievement
		err := rows.Scan(
			&achievement.ID,
			&achievement.Title,
			&achievement.Description,
			&achievement.Category,
			&achievement.ImpactMetric,
			&achievement.YearAchieved,
			&achievement.OrderIndex,
			&achievement.IsFeatured,
			&achievement.CreatedAt,
			&achievement.UpdatedAt,
		)
		if err != nil {
			return nil, repository.NewRepositoryError("scan", "achievement", err)
		}
		achievements = append(achievements, &achievement)
	}

	if err := rows.Err(); err != nil {
		return nil, repository.NewRepositoryError("iterate", "achievements", err)
	}

	return achievements, nil
}

// GetFeaturedAchievements retrieves only featured achievements
func (r *AchievementRepository) GetFeaturedAchievements(ctx context.Context) ([]*models.Achievement, error) {
	featured := true
	filters := repository.AchievementFilters{
		Featured: &featured,
	}
	return r.GetAchievements(ctx, filters)
}

// CreateAchievement creates a new achievement entry
func (r *AchievementRepository) CreateAchievement(ctx context.Context, achievement *models.Achievement) error {
	query := `
		INSERT INTO achievements (title, description, category, impact_metric, 
		                         year_achieved, order_index, is_featured)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		achievement.Title,
		achievement.Description,
		achievement.Category,
		achievement.ImpactMetric,
		achievement.YearAchieved,
		achievement.OrderIndex,
		achievement.IsFeatured,
	).Scan(&achievement.ID, &achievement.CreatedAt, &achievement.UpdatedAt)

	if err != nil {
		return repository.NewRepositoryError("create", "achievement", err)
	}

	return nil
}

// UpdateAchievement updates an existing achievement
func (r *AchievementRepository) UpdateAchievement(ctx context.Context, achievement *models.Achievement) error {
	query := `
		UPDATE achievements 
		SET title = $2, description = $3, category = $4, impact_metric = $5, 
		    year_achieved = $6, order_index = $7, is_featured = $8, 
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		achievement.ID,
		achievement.Title,
		achievement.Description,
		achievement.Category,
		achievement.ImpactMetric,
		achievement.YearAchieved,
		achievement.OrderIndex,
		achievement.IsFeatured,
	).Scan(&achievement.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return repository.NewRepositoryError("update", "achievement", fmt.Errorf("achievement with id %d not found", achievement.ID))
		}
		return repository.NewRepositoryError("update", "achievement", err)
	}

	return nil
}

// DeleteAchievement deletes an achievement by ID
func (r *AchievementRepository) DeleteAchievement(ctx context.Context, id int) error {
	query := `DELETE FROM achievements WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return repository.NewRepositoryError("delete", "achievement", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return repository.NewRepositoryError("delete", "achievement", fmt.Errorf("achievement with id %d not found", id))
	}

	return nil
}