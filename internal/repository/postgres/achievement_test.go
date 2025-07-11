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

func TestAchievementRepository(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	repo := NewAchievementRepository(testDB.Pool())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("CreateAchievement", func(t *testing.T) {
		testDB.CleanupTables(t)

		achievement := &models.Achievement{
			Title:        "Performance Optimization",
			Description:  stringPtr("Optimized database queries reducing response time by 40%"),
			Category:     stringPtr(models.AchievementCategoryPerformance),
			ImpactMetric: stringPtr("40% response time reduction"),
			YearAchieved: intPtr(2023),
			OrderIndex:   1,
			IsFeatured:   true,
		}

		err := repo.CreateAchievement(ctx, achievement)
		require.NoError(t, err)
		assert.NotZero(t, achievement.ID)
		assert.NotZero(t, achievement.CreatedAt)
		assert.NotZero(t, achievement.UpdatedAt)
	})

	t.Run("GetAchievements_All", func(t *testing.T) {
		testDB.CleanupTables(t)

		achievements := []*models.Achievement{
			{
				Title:        "Security Implementation",
				Description:  stringPtr("Implemented OAuth2 authentication system"),
				Category:     stringPtr(models.AchievementCategorySecurity),
				ImpactMetric: stringPtr("100% secure user authentication"),
				YearAchieved: intPtr(2024),
				OrderIndex:   1,
				IsFeatured:   true,
			},
			{
				Title:        "Team Leadership",
				Description:  stringPtr("Led team of 5 developers on major project"),
				Category:     stringPtr(models.AchievementCategoryLeadership),
				ImpactMetric: stringPtr("Project delivered 2 weeks early"),
				YearAchieved: intPtr(2023),
				OrderIndex:   2,
				IsFeatured:   false,
			},
			{
				Title:        "Innovation Award",
				Description:  stringPtr("Developed new microservices architecture"),
				Category:     stringPtr(models.AchievementCategoryInnovation),
				ImpactMetric: stringPtr("50% deployment time reduction"),
				YearAchieved: intPtr(2022),
				OrderIndex:   1,
				IsFeatured:   true,
			},
		}

		for _, achievement := range achievements {
			err := repo.CreateAchievement(ctx, achievement)
			require.NoError(t, err)
		}

		// Get all achievements
		filters := repository.AchievementFilters{}
		retrieved, err := repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 3)

		// Should be ordered by year_achieved DESC, order_index
		assert.Equal(t, "Security Implementation", retrieved[0].Title) // 2024
		assert.Equal(t, "Team Leadership", retrieved[1].Title)        // 2023
		assert.Equal(t, "Innovation Award", retrieved[2].Title)       // 2022
	})

	t.Run("GetAchievements_FilterByCategory", func(t *testing.T) {
		testDB.CleanupTables(t)

		achievements := []*models.Achievement{
			{
				Title:    "Performance Boost",
				Category: stringPtr(models.AchievementCategoryPerformance),
			},
			{
				Title:    "Security Audit",
				Category: stringPtr(models.AchievementCategorySecurity),
			},
			{
				Title:    "Performance Optimization",
				Category: stringPtr(models.AchievementCategoryPerformance),
			},
		}

		for _, achievement := range achievements {
			err := repo.CreateAchievement(ctx, achievement)
			require.NoError(t, err)
		}

		// Filter by performance category
		filters := repository.AchievementFilters{
			Category: models.AchievementCategoryPerformance,
		}
		retrieved, err := repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, achievement := range retrieved {
			assert.Equal(t, models.AchievementCategoryPerformance, *achievement.Category)
		}
	})

	t.Run("GetAchievements_FilterByYear", func(t *testing.T) {
		testDB.CleanupTables(t)

		achievements := []*models.Achievement{
			{
				Title:        "Achievement 2023",
				YearAchieved: intPtr(2023),
			},
			{
				Title:        "Achievement 2024-1",
				YearAchieved: intPtr(2024),
			},
			{
				Title:        "Achievement 2024-2",
				YearAchieved: intPtr(2024),
			},
		}

		for _, achievement := range achievements {
			err := repo.CreateAchievement(ctx, achievement)
			require.NoError(t, err)
		}

		// Filter by year 2024
		filters := repository.AchievementFilters{
			Year: intPtr(2024),
		}
		retrieved, err := repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, achievement := range retrieved {
			assert.Equal(t, 2024, *achievement.YearAchieved)
		}
	})

	t.Run("GetAchievements_FilterByFeatured", func(t *testing.T) {
		testDB.CleanupTables(t)

		achievements := []*models.Achievement{
			{
				Title:      "Featured Achievement 1",
				IsFeatured: true,
			},
			{
				Title:      "Regular Achievement",
				IsFeatured: false,
			},
			{
				Title:      "Featured Achievement 2",
				IsFeatured: true,
			},
		}

		for _, achievement := range achievements {
			err := repo.CreateAchievement(ctx, achievement)
			require.NoError(t, err)
		}

		// Filter by featured
		filters := repository.AchievementFilters{
			Featured: boolPtr(true),
		}
		retrieved, err := repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, achievement := range retrieved {
			assert.True(t, achievement.IsFeatured)
		}

		// Filter by non-featured
		filters.Featured = boolPtr(false)
		retrieved, err = repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.False(t, retrieved[0].IsFeatured)
	})

	t.Run("GetAchievements_CombinedFilters", func(t *testing.T) {
		testDB.CleanupTables(t)

		achievements := []*models.Achievement{
			{
				Title:        "Performance 2024 Featured",
				Category:     stringPtr(models.AchievementCategoryPerformance),
				YearAchieved: intPtr(2024),
				IsFeatured:   true,
			},
			{
				Title:        "Performance 2024 Regular",
				Category:     stringPtr(models.AchievementCategoryPerformance),
				YearAchieved: intPtr(2024),
				IsFeatured:   false,
			},
			{
				Title:        "Security 2024 Featured",
				Category:     stringPtr(models.AchievementCategorySecurity),
				YearAchieved: intPtr(2024),
				IsFeatured:   true,
			},
		}

		for _, achievement := range achievements {
			err := repo.CreateAchievement(ctx, achievement)
			require.NoError(t, err)
		}

		// Filter by category, year, and featured
		filters := repository.AchievementFilters{
			Category: models.AchievementCategoryPerformance,
			Year:     intPtr(2024),
			Featured: boolPtr(true),
		}
		retrieved, err := repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, "Performance 2024 Featured", retrieved[0].Title)
	})

	t.Run("GetAchievements_Pagination", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create achievements with different years
		for i := 0; i < 5; i++ {
			achievement := &models.Achievement{
				Title:        "Achievement " + string(rune('A'+i)),
				YearAchieved: intPtr(2024 - i), // 2024, 2023, 2022, 2021, 2020
				OrderIndex:   i,
			}
			err := repo.CreateAchievement(ctx, achievement)
			require.NoError(t, err)
		}

		// Get first page
		filters := repository.AchievementFilters{
			Limit:  2,
			Offset: 0,
		}
		page1, err := repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, page1, 2)
		assert.Equal(t, "Achievement A", page1[0].Title) // 2024 (most recent)
		assert.Equal(t, "Achievement B", page1[1].Title) // 2023

		// Get second page
		filters.Offset = 2
		page2, err := repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, page2, 2)
		assert.Equal(t, "Achievement C", page2[0].Title) // 2022
		assert.Equal(t, "Achievement D", page2[1].Title) // 2021
	})

	t.Run("GetFeaturedAchievements", func(t *testing.T) {
		testDB.CleanupTables(t)

		achievements := []*models.Achievement{
			{Title: "Featured 1", IsFeatured: true},
			{Title: "Regular 1", IsFeatured: false},
			{Title: "Featured 2", IsFeatured: true},
			{Title: "Regular 2", IsFeatured: false},
		}

		for _, achievement := range achievements {
			err := repo.CreateAchievement(ctx, achievement)
			require.NoError(t, err)
		}

		featured, err := repo.GetFeaturedAchievements(ctx)
		require.NoError(t, err)
		assert.Len(t, featured, 2)
		for _, achievement := range featured {
			assert.True(t, achievement.IsFeatured)
		}
	})

	t.Run("UpdateAchievement", func(t *testing.T) {
		testDB.CleanupTables(t)

		achievement := &models.Achievement{
			Title:        "Original Title",
			Description:  stringPtr("Original description"),
			Category:     stringPtr(models.AchievementCategoryPerformance),
			YearAchieved: intPtr(2023),
			IsFeatured:   false,
		}

		err := repo.CreateAchievement(ctx, achievement)
		require.NoError(t, err)
		originalUpdatedAt := achievement.UpdatedAt

		time.Sleep(time.Millisecond * 10)

		// Update achievement
		achievement.Title = "Updated Title"
		achievement.Description = stringPtr("Updated description")
		achievement.Category = stringPtr(models.AchievementCategoryInnovation)
		achievement.ImpactMetric = stringPtr("New impact metric")
		achievement.IsFeatured = true

		err = repo.UpdateAchievement(ctx, achievement)
		require.NoError(t, err)
		assert.True(t, achievement.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		filters := repository.AchievementFilters{}
		retrieved, err := repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		require.Len(t, retrieved, 1)

		updated := retrieved[0]
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, models.AchievementCategoryInnovation, *updated.Category)
		assert.Equal(t, "New impact metric", *updated.ImpactMetric)
		assert.True(t, updated.IsFeatured)
	})

	t.Run("UpdateAchievement_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		achievement := &models.Achievement{
			ID:    999,
			Title: "Non-existent",
		}

		err := repo.UpdateAchievement(ctx, achievement)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "achievement with id 999 not found")
	})

	t.Run("DeleteAchievement", func(t *testing.T) {
		testDB.CleanupTables(t)

		achievement := &models.Achievement{
			Title: "Delete Me",
		}

		err := repo.CreateAchievement(ctx, achievement)
		require.NoError(t, err)

		// Verify it exists
		filters := repository.AchievementFilters{}
		retrieved, err := repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)

		// Delete it
		err = repo.DeleteAchievement(ctx, achievement.ID)
		require.NoError(t, err)

		// Verify it's gone
		retrieved, err = repo.GetAchievements(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 0)
	})

	t.Run("DeleteAchievement_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		err := repo.DeleteAchievement(ctx, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "achievement with id 999 not found")
	})

	t.Run("AchievementCategories_Constants", func(t *testing.T) {
		// Test that our achievement category constants are available
		categories := []string{
			models.AchievementCategoryPerformance,
			models.AchievementCategorySecurity,
			models.AchievementCategoryLeadership,
			models.AchievementCategoryInnovation,
			models.AchievementCategoryEfficiency,
			models.AchievementCategoryTeamwork,
		}

		expectedCategories := []string{
			"performance",
			"security",
			"leadership",
			"innovation",
			"efficiency",
			"teamwork",
		}

		assert.Equal(t, expectedCategories, categories)
	})
}