package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/repository"
)

func TestSkillRepository(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	repo := NewSkillRepository(testDB.Pool())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("CreateSkill", func(t *testing.T) {
		testDB.CleanupTables(t)

		skill := &models.Skill{
			Category:        "Programming Languages",
			Name:            "Go",
			Level:           stringPtr(models.SkillLevelExpert),
			YearsExperience: intPtr(5),
			OrderIndex:      1,
			IsFeatured:      true,
		}

		err := repo.CreateSkill(ctx, skill)
		require.NoError(t, err)
		assert.NotZero(t, skill.ID)
		assert.NotZero(t, skill.CreatedAt)
		assert.NotZero(t, skill.UpdatedAt)
	})

	t.Run("GetSkills_All", func(t *testing.T) {
		testDB.CleanupTables(t)

		skills := []*models.Skill{
			{
				Category:        "Programming Languages",
				Name:            "Go",
				Level:           stringPtr(models.SkillLevelExpert),
				YearsExperience: intPtr(5),
				OrderIndex:      1,
				IsFeatured:      true,
			},
			{
				Category:        "Programming Languages",
				Name:            "Python",
				Level:           stringPtr(models.SkillLevelAdvanced),
				YearsExperience: intPtr(3),
				OrderIndex:      2,
				IsFeatured:      true,
			},
			{
				Category:        "Databases",
				Name:            "PostgreSQL",
				Level:           stringPtr(models.SkillLevelAdvanced),
				YearsExperience: intPtr(4),
				OrderIndex:      1,
				IsFeatured:      false,
			},
			{
				Category:        "Cloud Platforms",
				Name:            "AWS",
				Level:           stringPtr(models.SkillLevelIntermediate),
				YearsExperience: intPtr(2),
				OrderIndex:      1,
				IsFeatured:      true,
			},
		}

		for _, skill := range skills {
			err := repo.CreateSkill(ctx, skill)
			require.NoError(t, err)
		}

		// Get all skills
		filters := repository.SkillFilters{}
		retrieved, err := repo.GetSkills(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 4)

		// Should be ordered by category, order_index, name
		assert.Equal(t, "AWS", retrieved[0].Name)               // Cloud Platforms
		assert.Equal(t, "PostgreSQL", retrieved[1].Name)       // Databases
		assert.Equal(t, "Go", retrieved[2].Name)               // Programming Languages (order_index 1)
		assert.Equal(t, "Python", retrieved[3].Name)           // Programming Languages (order_index 2)
	})

	t.Run("GetSkills_FilterByCategory", func(t *testing.T) {
		testDB.CleanupTables(t)

		skills := []*models.Skill{
			{
				Category: "Programming Languages",
				Name:     "Go",
			},
			{
				Category: "Programming Languages",
				Name:     "Python",
			},
			{
				Category: "Databases",
				Name:     "PostgreSQL",
			},
		}

		for _, skill := range skills {
			err := repo.CreateSkill(ctx, skill)
			require.NoError(t, err)
		}

		// Filter by category
		filters := repository.SkillFilters{
			Category: "Programming Languages",
		}
		retrieved, err := repo.GetSkills(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, skill := range retrieved {
			assert.Equal(t, "Programming Languages", skill.Category)
		}
	})

	t.Run("GetSkills_FilterByLevel", func(t *testing.T) {
		testDB.CleanupTables(t)

		skills := []*models.Skill{
			{
				Category: "Programming",
				Name:     "Go",
				Level:    stringPtr(models.SkillLevelExpert),
			},
			{
				Category: "Programming",
				Name:     "Python",
				Level:    stringPtr(models.SkillLevelAdvanced),
			},
			{
				Category: "Programming",
				Name:     "JavaScript",
				Level:    stringPtr(models.SkillLevelExpert),
			},
		}

		for _, skill := range skills {
			err := repo.CreateSkill(ctx, skill)
			require.NoError(t, err)
		}

		// Filter by level
		filters := repository.SkillFilters{
			Level: models.SkillLevelExpert,
		}
		retrieved, err := repo.GetSkills(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, skill := range retrieved {
			assert.Equal(t, models.SkillLevelExpert, *skill.Level)
		}
	})

	t.Run("GetSkills_FilterByFeatured", func(t *testing.T) {
		testDB.CleanupTables(t)

		skills := []*models.Skill{
			{
				Category:   "Programming",
				Name:       "Go",
				IsFeatured: true,
			},
			{
				Category:   "Programming",
				Name:       "Python",
				IsFeatured: false,
			},
			{
				Category:   "Programming",
				Name:       "Rust",
				IsFeatured: true,
			},
		}

		for _, skill := range skills {
			err := repo.CreateSkill(ctx, skill)
			require.NoError(t, err)
		}

		// Filter by featured
		filters := repository.SkillFilters{
			Featured: boolPtr(true),
		}
		retrieved, err := repo.GetSkills(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, skill := range retrieved {
			assert.True(t, skill.IsFeatured)
		}

		// Filter by non-featured
		filters.Featured = boolPtr(false)
		retrieved, err = repo.GetSkills(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.False(t, retrieved[0].IsFeatured)
		assert.Equal(t, "Python", retrieved[0].Name)
	})

	t.Run("GetSkills_Pagination", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create 5 skills in same category
		for i := 0; i < 5; i++ {
			skill := &models.Skill{
				Category:   "Programming",
				Name:       "Language " + string(rune('A'+i)),
				OrderIndex: i,
			}
			err := repo.CreateSkill(ctx, skill)
			require.NoError(t, err)
		}

		// Get first page
		filters := repository.SkillFilters{
			Limit:  2,
			Offset: 0,
		}
		page1, err := repo.GetSkills(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, page1, 2)
		assert.Equal(t, "Language A", page1[0].Name)
		assert.Equal(t, "Language B", page1[1].Name)

		// Get second page
		filters.Offset = 2
		page2, err := repo.GetSkills(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, page2, 2)
		assert.Equal(t, "Language C", page2[0].Name)
		assert.Equal(t, "Language D", page2[1].Name)
	})

	t.Run("GetSkillsByCategory", func(t *testing.T) {
		testDB.CleanupTables(t)

		skills := []*models.Skill{
			{Category: "Frontend", Name: "React"},
			{Category: "Frontend", Name: "Vue"},
			{Category: "Backend", Name: "Go"},
			{Category: "Backend", Name: "Node.js"},
		}

		for _, skill := range skills {
			err := repo.CreateSkill(ctx, skill)
			require.NoError(t, err)
		}

		// Get frontend skills
		frontend, err := repo.GetSkillsByCategory(ctx, "Frontend")
		require.NoError(t, err)
		assert.Len(t, frontend, 2)
		for _, skill := range frontend {
			assert.Equal(t, "Frontend", skill.Category)
		}

		// Get backend skills
		backend, err := repo.GetSkillsByCategory(ctx, "Backend")
		require.NoError(t, err)
		assert.Len(t, backend, 2)
		for _, skill := range backend {
			assert.Equal(t, "Backend", skill.Category)
		}
	})

	t.Run("GetFeaturedSkills", func(t *testing.T) {
		testDB.CleanupTables(t)

		skills := []*models.Skill{
			{Category: "Programming", Name: "Go", IsFeatured: true},
			{Category: "Programming", Name: "Python", IsFeatured: false},
			{Category: "Cloud", Name: "AWS", IsFeatured: true},
			{Category: "Database", Name: "PostgreSQL", IsFeatured: false},
		}

		for _, skill := range skills {
			err := repo.CreateSkill(ctx, skill)
			require.NoError(t, err)
		}

		featured, err := repo.GetFeaturedSkills(ctx)
		require.NoError(t, err)
		assert.Len(t, featured, 2)
		for _, skill := range featured {
			assert.True(t, skill.IsFeatured)
		}
	})

	t.Run("UpdateSkill", func(t *testing.T) {
		testDB.CleanupTables(t)

		skill := &models.Skill{
			Category:        "Programming",
			Name:            "Go",
			Level:           stringPtr(models.SkillLevelIntermediate),
			YearsExperience: intPtr(2),
			OrderIndex:      1,
			IsFeatured:      false,
		}

		err := repo.CreateSkill(ctx, skill)
		require.NoError(t, err)
		originalUpdatedAt := skill.UpdatedAt

		time.Sleep(time.Millisecond * 10)

		// Update skill
		skill.Level = stringPtr(models.SkillLevelExpert)
		skill.YearsExperience = intPtr(5)
		skill.IsFeatured = true

		err = repo.UpdateSkill(ctx, skill)
		require.NoError(t, err)
		assert.True(t, skill.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		filters := repository.SkillFilters{Category: "Programming"}
		retrieved, err := repo.GetSkills(ctx, filters)
		require.NoError(t, err)
		require.Len(t, retrieved, 1)

		updated := retrieved[0]
		assert.Equal(t, models.SkillLevelExpert, *updated.Level)
		assert.Equal(t, 5, *updated.YearsExperience)
		assert.True(t, updated.IsFeatured)
	})

	t.Run("UpdateSkill_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		skill := &models.Skill{
			ID:       999,
			Category: "Non-existent",
			Name:     "Nobody",
		}

		err := repo.UpdateSkill(ctx, skill)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "skill with id 999 not found")
	})

	t.Run("DeleteSkill", func(t *testing.T) {
		testDB.CleanupTables(t)

		skill := &models.Skill{
			Category: "Temporary",
			Name:     "Delete Me",
		}

		err := repo.CreateSkill(ctx, skill)
		require.NoError(t, err)

		// Verify it exists
		filters := repository.SkillFilters{Category: "Temporary"}
		retrieved, err := repo.GetSkills(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)

		// Delete it
		err = repo.DeleteSkill(ctx, skill.ID)
		require.NoError(t, err)

		// Verify it's gone
		retrieved, err = repo.GetSkills(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 0)
	})

	t.Run("DeleteSkill_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		err := repo.DeleteSkill(ctx, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "skill with id 999 not found")
	})

	t.Run("SkillLevels_Validation", func(t *testing.T) {
		// Test that our skill level constants are valid
		validLevels := models.ValidSkillLevels()
		assert.Contains(t, validLevels, models.SkillLevelBeginner)
		assert.Contains(t, validLevels, models.SkillLevelIntermediate)
		assert.Contains(t, validLevels, models.SkillLevelAdvanced)
		assert.Contains(t, validLevels, models.SkillLevelExpert)
		assert.Len(t, validLevels, 4)
	})
}