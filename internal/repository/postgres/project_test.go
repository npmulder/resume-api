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

func TestProjectRepository(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	repo := NewProjectRepository(testDB.Pool())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("CreateProject", func(t *testing.T) {
		testDB.CleanupTables(t)

		startDate := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC)

		project := &models.Project{
			Name:             "Resume API",
			Description:      stringPtr("REST API for serving resume data built with Go and PostgreSQL"),
			ShortDescription: stringPtr("Go-based resume API with PostgreSQL backend"),
			Technologies:     []string{"Go", "PostgreSQL", "Docker", "Kubernetes"},
			GitHubURL:        stringPtr("https://github.com/user/resume-api"),
			DemoURL:          stringPtr("https://api.example.com"),
			StartDate:        &startDate,
			EndDate:          &endDate,
			Status:           models.ProjectStatusCompleted,
			IsFeatured:       true,
			OrderIndex:       1,
			KeyFeatures: []string{
				"RESTful API design",
				"PostgreSQL database",
				"Docker containerization",
				"Kubernetes deployment",
			},
		}

		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)
		assert.NotZero(t, project.ID)
		assert.NotZero(t, project.CreatedAt)
		assert.NotZero(t, project.UpdatedAt)
	})

	t.Run("CreateProject_OngoingProject", func(t *testing.T) {
		testDB.CleanupTables(t)

		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		project := &models.Project{
			Name:         "Homelab Infrastructure",
			Description:  stringPtr("Personal homelab setup with Kubernetes cluster"),
			Technologies: []string{"Kubernetes", "Docker", "Prometheus", "Grafana"},
			StartDate:    &startDate,
			EndDate:      nil, // Ongoing project
			Status:       models.ProjectStatusActive,
			IsFeatured:   true,
		}

		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)
		assert.NotZero(t, project.ID)
		assert.True(t, project.IsOngoing())
	})

	t.Run("GetProjectByID", func(t *testing.T) {
		testDB.CleanupTables(t)

		project := &models.Project{
			Name:             "Test Project",
			ShortDescription: stringPtr("A test project"),
			Technologies:     []string{"Go", "Testing"},
			Status:           models.ProjectStatusActive,
			KeyFeatures:      []string{"Feature 1", "Feature 2"},
		}

		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)

		// Retrieve by ID
		retrieved, err := repo.GetProjectByID(ctx, project.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		assert.Equal(t, project.ID, retrieved.ID)
		assert.Equal(t, project.Name, retrieved.Name)
		assert.Equal(t, project.ShortDescription, retrieved.ShortDescription)
		assert.Equal(t, project.Technologies, retrieved.Technologies)
		assert.Equal(t, project.Status, retrieved.Status)
		assert.Equal(t, project.KeyFeatures, retrieved.KeyFeatures)
	})

	t.Run("GetProjectByID_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		project, err := repo.GetProjectByID(ctx, 999)
		assert.Error(t, err)
		assert.Nil(t, project)
		assert.Contains(t, err.Error(), "project with id 999 not found")
	})

	t.Run("GetProjects_All", func(t *testing.T) {
		testDB.CleanupTables(t)

		projects := []*models.Project{
			{
				Name:       "Project A",
				StartDate:  timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				Status:     models.ProjectStatusActive,
				OrderIndex: 1,
			},
			{
				Name:       "Project B",
				StartDate:  timePtr(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)),
				Status:     models.ProjectStatusCompleted,
				OrderIndex: 2,
			},
			{
				Name:       "Project C",
				StartDate:  timePtr(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				Status:     models.ProjectStatusArchived,
				OrderIndex: 3,
			},
		}

		for _, project := range projects {
			err := repo.CreateProject(ctx, project)
			require.NoError(t, err)
		}

		// Get all projects
		filters := repository.ProjectFilters{}
		retrieved, err := repo.GetProjects(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 3)

		// Should be ordered by start_date DESC, order_index
		assert.Equal(t, "Project A", retrieved[0].Name) // 2024 (most recent)
		assert.Equal(t, "Project B", retrieved[1].Name) // 2023-06
		assert.Equal(t, "Project C", retrieved[2].Name) // 2023-01
	})

	t.Run("GetProjects_FilterByStatus", func(t *testing.T) {
		testDB.CleanupTables(t)

		projects := []*models.Project{
			{Name: "Active Project", Status: models.ProjectStatusActive},
			{Name: "Completed Project", Status: models.ProjectStatusCompleted},
			{Name: "Another Active", Status: models.ProjectStatusActive},
			{Name: "Archived Project", Status: models.ProjectStatusArchived},
		}

		for _, project := range projects {
			err := repo.CreateProject(ctx, project)
			require.NoError(t, err)
		}

		// Filter by active status
		filters := repository.ProjectFilters{
			Status: models.ProjectStatusActive,
		}
		retrieved, err := repo.GetProjects(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, project := range retrieved {
			assert.Equal(t, models.ProjectStatusActive, project.Status)
		}

		// Filter by completed status
		filters.Status = models.ProjectStatusCompleted
		retrieved, err = repo.GetProjects(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, "Completed Project", retrieved[0].Name)
	})

	t.Run("GetProjects_FilterByTechnology", func(t *testing.T) {
		testDB.CleanupTables(t)

		projects := []*models.Project{
			{
				Name:         "Go Project",
				Technologies: []string{"Go", "PostgreSQL", "Docker"},
				Status:       models.ProjectStatusCompleted,
			},
			{
				Name:         "Python Project",
				Technologies: []string{"Python", "Django", "Redis"},
				Status:       models.ProjectStatusCompleted,
			},
			{
				Name:         "Microservices Project",
				Technologies: []string{"Go", "Kubernetes", "gRPC"},
				Status:       models.ProjectStatusActive,
			},
		}

		for _, project := range projects {
			err := repo.CreateProject(ctx, project)
			require.NoError(t, err)
		}

		// Filter by Go technology
		filters := repository.ProjectFilters{
			Technology: "Go",
		}
		retrieved, err := repo.GetProjects(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, project := range retrieved {
			assert.Contains(t, project.Technologies, "Go")
		}

		// Filter by Django technology
		filters.Technology = "Django"
		retrieved, err = repo.GetProjects(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, "Python Project", retrieved[0].Name)
	})

	t.Run("GetProjects_FilterByFeatured", func(t *testing.T) {
		testDB.CleanupTables(t)

		projects := []*models.Project{
			{Name: "Featured Project 1", Status: models.ProjectStatusActive, IsFeatured: true},
			{Name: "Regular Project", Status: models.ProjectStatusCompleted, IsFeatured: false},
			{Name: "Featured Project 2", Status: models.ProjectStatusActive, IsFeatured: true},
		}

		for _, project := range projects {
			err := repo.CreateProject(ctx, project)
			require.NoError(t, err)
		}

		// Filter by featured
		filters := repository.ProjectFilters{
			Featured: boolPtr(true),
		}
		retrieved, err := repo.GetProjects(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, project := range retrieved {
			assert.True(t, project.IsFeatured)
		}

		// Filter by non-featured
		filters.Featured = boolPtr(false)
		retrieved, err = repo.GetProjects(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.False(t, retrieved[0].IsFeatured)
	})

	t.Run("GetProjects_CombinedFilters", func(t *testing.T) {
		testDB.CleanupTables(t)

		projects := []*models.Project{
			{
				Name:         "Active Go Featured",
				Status:       models.ProjectStatusActive,
				Technologies: []string{"Go", "Docker"},
				IsFeatured:   true,
			},
			{
				Name:         "Active Go Regular",
				Status:       models.ProjectStatusActive,
				Technologies: []string{"Go", "Kubernetes"},
				IsFeatured:   false,
			},
			{
				Name:         "Completed Go Featured",
				Status:       models.ProjectStatusCompleted,
				Technologies: []string{"Go", "PostgreSQL"},
				IsFeatured:   true,
			},
		}

		for _, project := range projects {
			err := repo.CreateProject(ctx, project)
			require.NoError(t, err)
		}

		// Filter by status, technology, and featured
		filters := repository.ProjectFilters{
			Status:     models.ProjectStatusActive,
			Technology: "Go",
			Featured:   boolPtr(true),
		}
		retrieved, err := repo.GetProjects(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, "Active Go Featured", retrieved[0].Name)
	})

	t.Run("GetProjects_Pagination", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create projects with different start dates
		for i := 0; i < 5; i++ {
			project := &models.Project{
				Name:       "Project " + string(rune('A'+i)),
				StartDate:  timePtr(time.Date(2024-i, 1, 1, 0, 0, 0, 0, time.UTC)),
				Status:     models.ProjectStatusCompleted,
				OrderIndex: i,
			}
			err := repo.CreateProject(ctx, project)
			require.NoError(t, err)
		}

		// Get first page
		filters := repository.ProjectFilters{
			Limit:  2,
			Offset: 0,
		}
		page1, err := repo.GetProjects(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, page1, 2)
		assert.Equal(t, "Project A", page1[0].Name) // 2024 (most recent)
		assert.Equal(t, "Project B", page1[1].Name) // 2023

		// Get second page
		filters.Offset = 2
		page2, err := repo.GetProjects(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, page2, 2)
		assert.Equal(t, "Project C", page2[0].Name) // 2022
		assert.Equal(t, "Project D", page2[1].Name) // 2021
	})

	t.Run("GetFeaturedProjects", func(t *testing.T) {
		testDB.CleanupTables(t)

		projects := []*models.Project{
			{Name: "Featured 1", Status: models.ProjectStatusActive, IsFeatured: true},
			{Name: "Regular 1", Status: models.ProjectStatusCompleted, IsFeatured: false},
			{Name: "Featured 2", Status: models.ProjectStatusActive, IsFeatured: true},
			{Name: "Regular 2", Status: models.ProjectStatusCompleted, IsFeatured: false},
		}

		for _, project := range projects {
			err := repo.CreateProject(ctx, project)
			require.NoError(t, err)
		}

		featured, err := repo.GetFeaturedProjects(ctx)
		require.NoError(t, err)
		assert.Len(t, featured, 2)
		for _, project := range featured {
			assert.True(t, project.IsFeatured)
		}
	})

	t.Run("UpdateProject", func(t *testing.T) {
		testDB.CleanupTables(t)

		project := &models.Project{
			Name:         "Original Project",
			Description:  stringPtr("Original description"),
			Technologies: []string{"Go"},
			Status:       models.ProjectStatusActive,
			IsFeatured:   false,
		}

		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)
		originalUpdatedAt := project.UpdatedAt

		time.Sleep(time.Millisecond * 10)

		// Update project
		endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
		project.Name = "Updated Project"
		project.Description = stringPtr("Updated description")
		project.Technologies = []string{"Go", "PostgreSQL", "Docker"}
		project.Status = models.ProjectStatusCompleted
		project.EndDate = &endDate
		project.IsFeatured = true
		project.KeyFeatures = []string{"New feature 1", "New feature 2"}

		err = repo.UpdateProject(ctx, project)
		require.NoError(t, err)
		assert.True(t, project.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		updated, err := repo.GetProjectByID(ctx, project.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Project", updated.Name)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, []string{"Go", "PostgreSQL", "Docker"}, updated.Technologies)
		assert.Equal(t, models.ProjectStatusCompleted, updated.Status)
		assert.NotNil(t, updated.EndDate)
		assert.True(t, updated.IsFeatured)
		assert.Equal(t, []string{"New feature 1", "New feature 2"}, updated.KeyFeatures)
	})

	t.Run("UpdateProject_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		project := &models.Project{
			ID:   999,
			Name: "Non-existent",
		}

		err := repo.UpdateProject(ctx, project)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "project with id 999 not found")
	})

	t.Run("DeleteProject", func(t *testing.T) {
		testDB.CleanupTables(t)

		project := &models.Project{
			Name:   "Delete Me Project",
			Status: models.ProjectStatusCompleted,
		}

		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)

		// Verify it exists
		retrieved, err := repo.GetProjectByID(ctx, project.ID)
		require.NoError(t, err)
		assert.Equal(t, "Delete Me Project", retrieved.Name)

		// Delete it
		err = repo.DeleteProject(ctx, project.ID)
		require.NoError(t, err)

		// Verify it's gone
		_, err = repo.GetProjectByID(ctx, project.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("DeleteProject_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		err := repo.DeleteProject(ctx, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "project with id 999 not found")
	})

	t.Run("ProjectStatuses_Validation", func(t *testing.T) {
		// Test project status constants
		statuses := models.ValidProjectStatuses()
		assert.Contains(t, statuses, models.ProjectStatusActive)
		assert.Contains(t, statuses, models.ProjectStatusCompleted)
		assert.Contains(t, statuses, models.ProjectStatusArchived)
		assert.Contains(t, statuses, models.ProjectStatusPlanned)
		assert.Len(t, statuses, 4)
	})

	t.Run("Project_IsOngoing", func(t *testing.T) {
		// Active project with no end date should be ongoing
		activeProject := &models.Project{
			Status:  models.ProjectStatusActive,
			EndDate: nil,
		}
		assert.True(t, activeProject.IsOngoing())

		// Completed project should not be ongoing
		completedProject := &models.Project{
			Status:  models.ProjectStatusCompleted,
			EndDate: timePtr(time.Now()),
		}
		assert.False(t, completedProject.IsOngoing())

		// Active project with end date should not be ongoing
		endedProject := &models.Project{
			Status:  models.ProjectStatusActive,
			EndDate: timePtr(time.Now()),
		}
		assert.False(t, endedProject.IsOngoing())
	})

	t.Run("Project_JSONBFields", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Test with complex JSONB data
		project := &models.Project{
			Name:   "Complex Project",
			Status: models.ProjectStatusActive,
			Technologies: []string{
				"Go", "PostgreSQL", "Docker", "Kubernetes", 
				"Prometheus", "Grafana", "React", "TypeScript",
			},
			KeyFeatures: []string{
				"Microservices architecture",
				"Event-driven design",
				"Real-time monitoring",
				"Automated deployment",
				"API versioning",
			},
		}

		err := repo.CreateProject(ctx, project)
		require.NoError(t, err)

		// Retrieve and verify JSONB fields
		retrieved, err := repo.GetProjectByID(ctx, project.ID)
		require.NoError(t, err)
		
		assert.Equal(t, project.Technologies, retrieved.Technologies)
		assert.Equal(t, project.KeyFeatures, retrieved.KeyFeatures)
		assert.Len(t, retrieved.Technologies, 8)
		assert.Len(t, retrieved.KeyFeatures, 5)
	})
}