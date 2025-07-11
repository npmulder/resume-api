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

func TestExperienceRepository(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	repo := NewExperienceRepository(testDB.Pool())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("CreateExperience", func(t *testing.T) {
		testDB.CleanupTables(t)

		startDate := time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

		experience := &models.Experience{
			Company:     "TechCorp Inc",
			Position:    "Senior Software Engineer",
			StartDate:   startDate,
			EndDate:     &endDate,
			Description: stringPtr("Led development of microservices architecture"),
			Highlights: []string{
				"Designed and implemented REST APIs",
				"Reduced system latency by 40%",
				"Mentored junior developers",
			},
			OrderIndex: 1,
		}

		err := repo.CreateExperience(ctx, experience)
		require.NoError(t, err)
		assert.NotZero(t, experience.ID)
		assert.NotZero(t, experience.CreatedAt)
		assert.NotZero(t, experience.UpdatedAt)
	})

	t.Run("GetExperienceByID", func(t *testing.T) {
		testDB.CleanupTables(t)

		startDate := time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC)

		experience := &models.Experience{
			Company:     "StartupXYZ",
			Position:    "Full Stack Developer",
			StartDate:   startDate,
			EndDate:     nil, // Current position
			Description: stringPtr("Full stack development using Go and React"),
			Highlights: []string{
				"Built user authentication system",
				"Implemented real-time notifications",
			},
			OrderIndex: 0,
		}

		err := repo.CreateExperience(ctx, experience)
		require.NoError(t, err)

		// Retrieve by ID
		retrieved, err := repo.GetExperienceByID(ctx, experience.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		assert.Equal(t, experience.ID, retrieved.ID)
		assert.Equal(t, experience.Company, retrieved.Company)
		assert.Equal(t, experience.Position, retrieved.Position)
		assert.Equal(t, experience.StartDate.Unix(), retrieved.StartDate.Unix())
		assert.Nil(t, retrieved.EndDate) // Current position
		assert.Equal(t, experience.Description, retrieved.Description)
		assert.Equal(t, experience.Highlights, retrieved.Highlights)
		assert.Equal(t, experience.OrderIndex, retrieved.OrderIndex)
		assert.True(t, retrieved.IsCurrentPosition())
	})

	t.Run("GetExperienceByID_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		experience, err := repo.GetExperienceByID(ctx, 999)
		assert.Error(t, err)
		assert.Nil(t, experience)
		assert.Contains(t, err.Error(), "experience with id 999 not found")
	})

	t.Run("GetExperiences_All", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create multiple experiences
		experiences := []*models.Experience{
			{
				Company:    "Company A",
				Position:   "Engineer",
				StartDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:    timePtr(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)),
				OrderIndex: 0,
			},
			{
				Company:    "Company B", 
				Position:   "Senior Engineer",
				StartDate:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:    nil, // Current
				OrderIndex: 1,
			},
			{
				Company:    "Company C",
				Position:   "Junior Engineer",
				StartDate:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:    timePtr(time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)),
				OrderIndex: 2,
			},
		}

		for _, exp := range experiences {
			err := repo.CreateExperience(ctx, exp)
			require.NoError(t, err)
		}

		// Get all experiences
		filters := repository.ExperienceFilters{}
		retrieved, err := repo.GetExperiences(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 3)

		// Should be ordered by start_date DESC
		assert.Equal(t, "Company B", retrieved[0].Company) // 2024
		assert.Equal(t, "Company A", retrieved[1].Company) // 2023
		assert.Equal(t, "Company C", retrieved[2].Company) // 2022
	})

	t.Run("GetExperiences_FilterByCompany", func(t *testing.T) {
		testDB.CleanupTables(t)

		experiences := []*models.Experience{
			{
				Company:   "Google",
				Position:  "Software Engineer",
				StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Company:   "Microsoft",
				Position:  "Senior Engineer",
				StartDate: time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Company:   "Google",
				Position:  "Staff Engineer",
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}

		for _, exp := range experiences {
			err := repo.CreateExperience(ctx, exp)
			require.NoError(t, err)
		}

		// Filter by company
		filters := repository.ExperienceFilters{
			Company: "Google",
		}
		retrieved, err := repo.GetExperiences(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		assert.Equal(t, "Staff Engineer", retrieved[0].Position) // More recent first
		assert.Equal(t, "Software Engineer", retrieved[1].Position)
	})

	t.Run("GetExperiences_FilterByPosition", func(t *testing.T) {
		testDB.CleanupTables(t)

		experiences := []*models.Experience{
			{
				Company:   "CompanyA",
				Position:  "Software Engineer",
				StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Company:   "CompanyB",
				Position:  "Senior Software Engineer",
				StartDate: time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Company:   "CompanyC",
				Position:  "Data Scientist",
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}

		for _, exp := range experiences {
			err := repo.CreateExperience(ctx, exp)
			require.NoError(t, err)
		}

		// Filter by position (partial match)
		filters := repository.ExperienceFilters{
			Position: "Software",
		}
		retrieved, err := repo.GetExperiences(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, exp := range retrieved {
			assert.Contains(t, exp.Position, "Software")
		}
	})

	t.Run("GetExperiences_FilterByCurrent", func(t *testing.T) {
		testDB.CleanupTables(t)

		experiences := []*models.Experience{
			{
				Company:   "Current Co",
				Position:  "Engineer",
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   nil, // Current
			},
			{
				Company:   "Previous Co",
				Position:  "Engineer",
				StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   timePtr(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
		}

		for _, exp := range experiences {
			err := repo.CreateExperience(ctx, exp)
			require.NoError(t, err)
		}

		// Filter for current positions only
		filters := repository.ExperienceFilters{
			IsCurrent: boolPtr(true),
		}
		retrieved, err := repo.GetExperiences(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, "Current Co", retrieved[0].Company)
		assert.Nil(t, retrieved[0].EndDate)

		// Filter for past positions only
		filters = repository.ExperienceFilters{
			IsCurrent: boolPtr(false),
		}
		retrieved, err = repo.GetExperiences(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, "Previous Co", retrieved[0].Company)
		assert.NotNil(t, retrieved[0].EndDate)
	})

	t.Run("GetExperiences_FilterByDateRange", func(t *testing.T) {
		testDB.CleanupTables(t)

		experiences := []*models.Experience{
			{
				Company:   "Early Co",
				Position:  "Engineer",
				StartDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Company:   "Mid Co",
				Position:  "Engineer",
				StartDate: time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				Company:   "Recent Co",
				Position:  "Engineer",
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}

		for _, exp := range experiences {
			err := repo.CreateExperience(ctx, exp)
			require.NoError(t, err)
		}

		// Filter by date range
		from := "2022-01-01"
		to := "2023-12-31"
		filters := repository.ExperienceFilters{
			DateFrom: &from,
			DateTo:   &to,
		}
		retrieved, err := repo.GetExperiences(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, "Mid Co", retrieved[0].Company)
	})

	t.Run("GetExperiences_Pagination", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create 5 experiences
		for i := 0; i < 5; i++ {
			exp := &models.Experience{
				Company:   "Company " + string(rune('A'+i)),
				Position:  "Engineer",
				StartDate: time.Date(2024-i, 1, 1, 0, 0, 0, 0, time.UTC),
			}
			err := repo.CreateExperience(ctx, exp)
			require.NoError(t, err)
		}

		// Get first page
		filters := repository.ExperienceFilters{
			Limit:  2,
			Offset: 0,
		}
		page1, err := repo.GetExperiences(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, page1, 2)
		assert.Equal(t, "Company A", page1[0].Company) // Most recent

		// Get second page
		filters.Offset = 2
		page2, err := repo.GetExperiences(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, page2, 2)
		assert.Equal(t, "Company C", page2[0].Company)
	})

	t.Run("UpdateExperience", func(t *testing.T) {
		testDB.CleanupTables(t)

		experience := &models.Experience{
			Company:    "Original Company",
			Position:   "Junior Engineer",
			StartDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			OrderIndex: 0,
		}

		err := repo.CreateExperience(ctx, experience)
		require.NoError(t, err)
		originalUpdatedAt := experience.UpdatedAt

		time.Sleep(time.Millisecond * 10)

		// Update experience
		experience.Company = "Updated Company"
		experience.Position = "Senior Engineer"
		experience.EndDate = timePtr(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC))
		experience.Highlights = []string{"Updated achievements"}

		err = repo.UpdateExperience(ctx, experience)
		require.NoError(t, err)
		assert.True(t, experience.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		updated, err := repo.GetExperienceByID(ctx, experience.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Company", updated.Company)
		assert.Equal(t, "Senior Engineer", updated.Position)
		assert.NotNil(t, updated.EndDate)
		assert.Equal(t, []string{"Updated achievements"}, updated.Highlights)
	})

	t.Run("UpdateExperience_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		experience := &models.Experience{
			ID:        999,
			Company:   "Non-existent",
			Position:  "Nobody",
			StartDate: time.Now(),
		}

		err := repo.UpdateExperience(ctx, experience)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "experience with id 999 not found")
	})

	t.Run("DeleteExperience", func(t *testing.T) {
		testDB.CleanupTables(t)

		experience := &models.Experience{
			Company:   "Delete Me Inc",
			Position:  "Temporary",
			StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		err := repo.CreateExperience(ctx, experience)
		require.NoError(t, err)

		// Verify it exists
		retrieved, err := repo.GetExperienceByID(ctx, experience.ID)
		require.NoError(t, err)
		assert.Equal(t, "Delete Me Inc", retrieved.Company)

		// Delete it
		err = repo.DeleteExperience(ctx, experience.ID)
		require.NoError(t, err)

		// Verify it's gone
		_, err = repo.GetExperienceByID(ctx, experience.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("DeleteExperience_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		err := repo.DeleteExperience(ctx, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "experience with id 999 not found")
	})
}