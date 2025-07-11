package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npmulder/resume-api/internal/models"
)

func TestProfileRepository(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	repo := NewProfileRepository(testDB.Pool())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("CreateProfile", func(t *testing.T) {
		testDB.CleanupTables(t)

		profile := &models.Profile{
			Name:     "John Doe",
			Title:    "Software Engineer",
			Email:    "john.doe@example.com",
			Phone:    stringPtr("+1-555-0123"),
			Location: stringPtr("San Francisco, CA"),
			LinkedIn: stringPtr("https://linkedin.com/in/johndoe"),
			GitHub:   stringPtr("https://github.com/johndoe"),
			Summary:  stringPtr("Experienced software engineer with 5+ years in Go development"),
		}

		err := repo.CreateProfile(ctx, profile)
		require.NoError(t, err)
		assert.NotZero(t, profile.ID)
		assert.NotZero(t, profile.CreatedAt)
		assert.NotZero(t, profile.UpdatedAt)
		assert.Equal(t, profile.CreatedAt, profile.UpdatedAt)
	})

	t.Run("GetProfile", func(t *testing.T) {
		testDB.CleanupTables(t)

		// First create a profile
		profile := &models.Profile{
			Name:     "Jane Smith",
			Title:    "DevOps Engineer", 
			Email:    "jane.smith@example.com",
			Phone:    stringPtr("+1-555-0456"),
			Location: stringPtr("Austin, TX"),
			LinkedIn: stringPtr("https://linkedin.com/in/janesmith"),
			GitHub:   stringPtr("https://github.com/janesmith"),
			Summary:  stringPtr("DevOps engineer specializing in Kubernetes and cloud infrastructure"),
		}

		err := repo.CreateProfile(ctx, profile)
		require.NoError(t, err)

		// Now retrieve it
		retrieved, err := repo.GetProfile(ctx)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		assert.Equal(t, profile.ID, retrieved.ID)
		assert.Equal(t, profile.Name, retrieved.Name)
		assert.Equal(t, profile.Title, retrieved.Title)
		assert.Equal(t, profile.Email, retrieved.Email)
		assert.Equal(t, profile.Phone, retrieved.Phone)
		assert.Equal(t, profile.Location, retrieved.Location)
		assert.Equal(t, profile.LinkedIn, retrieved.LinkedIn)
		assert.Equal(t, profile.GitHub, retrieved.GitHub)
		assert.Equal(t, profile.Summary, retrieved.Summary)
		assert.Equal(t, profile.CreatedAt.Unix(), retrieved.CreatedAt.Unix())
		assert.Equal(t, profile.UpdatedAt.Unix(), retrieved.UpdatedAt.Unix())
	})

	t.Run("GetProfile_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Try to get profile when none exists
		profile, err := repo.GetProfile(ctx)
		assert.Error(t, err)
		assert.Nil(t, profile)
		assert.Contains(t, err.Error(), "profile not found")
	})

	t.Run("UpdateProfile", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create initial profile
		profile := &models.Profile{
			Name:     "Bob Wilson",
			Title:    "Backend Developer",
			Email:    "bob.wilson@example.com",
			Phone:    stringPtr("+1-555-0789"),
			Location: stringPtr("Seattle, WA"),
			LinkedIn: stringPtr("https://linkedin.com/in/bobwilson"),
			GitHub:   stringPtr("https://github.com/bobwilson"),
			Summary:  stringPtr("Backend developer focused on microservices"),
		}

		err := repo.CreateProfile(ctx, profile)
		require.NoError(t, err)
		originalUpdatedAt := profile.UpdatedAt

		// Wait a moment to ensure updated_at changes
		time.Sleep(time.Millisecond * 10)

		// Update the profile
		profile.Name = "Robert Wilson"
		profile.Title = "Senior Backend Developer"
		profile.Summary = stringPtr("Senior backend developer with expertise in distributed systems")

		err = repo.UpdateProfile(ctx, profile)
		require.NoError(t, err)
		assert.True(t, profile.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be newer")

		// Verify the update
		updated, err := repo.GetProfile(ctx)
		require.NoError(t, err)
		assert.Equal(t, "Robert Wilson", updated.Name)
		assert.Equal(t, "Senior Backend Developer", updated.Title)
		assert.Equal(t, "Senior backend developer with expertise in distributed systems", *updated.Summary)
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("UpdateProfile_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Try to update non-existent profile
		profile := &models.Profile{
			ID:    999,
			Name:  "Non Existent",
			Title: "Nobody",
			Email: "nobody@nowhere.com",
		}

		err := repo.UpdateProfile(ctx, profile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile with id 999 not found")
	})

	t.Run("CreateProfile_DuplicateEmail", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create first profile
		profile1 := &models.Profile{
			Name:  "User One",
			Title: "Developer",
			Email: "same@email.com",
		}
		err := repo.CreateProfile(ctx, profile1)
		require.NoError(t, err)

		// Try to create second profile with same email
		profile2 := &models.Profile{
			Name:  "User Two", 
			Title: "Designer",
			Email: "same@email.com", // Same email
		}
		err = repo.CreateProfile(ctx, profile2)
		assert.Error(t, err, "Should fail due to unique email constraint")
	})

	t.Run("CreateProfile_MinimalData", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create profile with only required fields
		profile := &models.Profile{
			Name:  "Minimal User",
			Title: "Tester",
			Email: "minimal@test.com",
			// All optional fields are nil
		}

		err := repo.CreateProfile(ctx, profile)
		require.NoError(t, err)
		assert.NotZero(t, profile.ID)

		// Verify it can be retrieved
		retrieved, err := repo.GetProfile(ctx)
		require.NoError(t, err)
		assert.Equal(t, "Minimal User", retrieved.Name)
		assert.Equal(t, "Tester", retrieved.Title)
		assert.Equal(t, "minimal@test.com", retrieved.Email)
		assert.Nil(t, retrieved.Phone)
		assert.Nil(t, retrieved.Location)
		assert.Nil(t, retrieved.LinkedIn)
		assert.Nil(t, retrieved.GitHub)
		assert.Nil(t, retrieved.Summary)
	})
}