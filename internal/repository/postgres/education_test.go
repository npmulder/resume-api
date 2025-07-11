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

func TestEducationRepository(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	repo := NewEducationRepository(testDB.Pool())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("CreateEducation", func(t *testing.T) {
		testDB.CleanupTables(t)

		education := &models.Education{
			Institution:           "University of Technology",
			DegreeOrCertification: "Bachelor of Science in Computer Science",
			FieldOfStudy:          stringPtr("Computer Science"),
			YearCompleted:         intPtr(2020),
			YearStarted:           intPtr(2016),
			Description:           stringPtr("Focused on software engineering and algorithms"),
			Type:                  models.EducationTypeEducation,
			Status:                models.EducationStatusCompleted,
			OrderIndex:            1,
			IsFeatured:            true,
		}

		err := repo.CreateEducation(ctx, education)
		require.NoError(t, err)
		assert.NotZero(t, education.ID)
		assert.NotZero(t, education.CreatedAt)
		assert.NotZero(t, education.UpdatedAt)
	})

	t.Run("CreateEducation_Certification", func(t *testing.T) {
		testDB.CleanupTables(t)

		expiryDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
		certification := &models.Education{
			Institution:           "AWS",
			DegreeOrCertification: "AWS Certified Solutions Architect",
			FieldOfStudy:          stringPtr("Cloud Computing"),
			YearCompleted:         intPtr(2023),
			Description:           stringPtr("Professional level certification in AWS cloud services"),
			Type:                  models.EducationTypeCertification,
			Status:                models.EducationStatusCompleted,
			CredentialID:          stringPtr("AWS-CSA-123456"),
			CredentialURL:         stringPtr("https://aws.amazon.com/verification/123456"),
			ExpiryDate:            &expiryDate,
			OrderIndex:            1,
			IsFeatured:            true,
		}

		err := repo.CreateEducation(ctx, certification)
		require.NoError(t, err)
		assert.NotZero(t, certification.ID)
		assert.Equal(t, models.EducationTypeCertification, certification.Type)
		assert.NotNil(t, certification.ExpiryDate)
		assert.NotNil(t, certification.CredentialID)
	})

	t.Run("GetEducation_All", func(t *testing.T) {
		testDB.CleanupTables(t)

		educations := []*models.Education{
			{
				Institution:           "University A",
				DegreeOrCertification: "Bachelor's Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusCompleted,
				YearCompleted:         intPtr(2020),
				OrderIndex:            1,
			},
			{
				Institution:           "Certification Authority",
				DegreeOrCertification: "Professional Certification",
				Type:                  models.EducationTypeCertification,
				Status:                models.EducationStatusCompleted,
				YearCompleted:         intPtr(2023),
				OrderIndex:            1,
			},
			{
				Institution:           "University B",
				DegreeOrCertification: "Master's Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusInProgress,
				YearStarted:           intPtr(2024),
				OrderIndex:            2,
			},
		}

		for _, education := range educations {
			err := repo.CreateEducation(ctx, education)
			require.NoError(t, err)
		}

		// Get all education
		filters := repository.EducationFilters{}
		retrieved, err := repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 3)

		// Should be ordered by type, year_completed DESC, order_index
		assert.Equal(t, "Professional Certification", retrieved[0].DegreeOrCertification) // Certification
		assert.Equal(t, "Master's Degree", retrieved[1].DegreeOrCertification)           // Education (no completion year, so nil sorts last)
		assert.Equal(t, "Bachelor's Degree", retrieved[2].DegreeOrCertification)         // Education (2020)
	})

	t.Run("GetEducation_FilterByType", func(t *testing.T) {
		testDB.CleanupTables(t)

		educations := []*models.Education{
			{
				Institution:           "University",
				DegreeOrCertification: "Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusCompleted,
			},
			{
				Institution:           "Cert Authority",
				DegreeOrCertification: "Certification",
				Type:                  models.EducationTypeCertification,
				Status:                models.EducationStatusCompleted,
			},
			{
				Institution:           "Another University",
				DegreeOrCertification: "Another Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusCompleted,
			},
		}

		for _, education := range educations {
			err := repo.CreateEducation(ctx, education)
			require.NoError(t, err)
		}

		// Filter by education type
		filters := repository.EducationFilters{
			Type: models.EducationTypeEducation,
		}
		retrieved, err := repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, education := range retrieved {
			assert.Equal(t, models.EducationTypeEducation, education.Type)
		}

		// Filter by certification type
		filters.Type = models.EducationTypeCertification
		retrieved, err = repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, models.EducationTypeCertification, retrieved[0].Type)
	})

	t.Run("GetEducation_FilterByInstitution", func(t *testing.T) {
		testDB.CleanupTables(t)

		educations := []*models.Education{
			{
				Institution:           "MIT",
				DegreeOrCertification: "Computer Science Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusCompleted,
			},
			{
				Institution:           "Stanford University",
				DegreeOrCertification: "Engineering Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusCompleted,
			},
			{
				Institution:           "MIT OpenCourseWare",
				DegreeOrCertification: "Online Course",
				Type:                  models.EducationTypeCertification,
				Status:                models.EducationStatusCompleted,
			},
		}

		for _, education := range educations {
			err := repo.CreateEducation(ctx, education)
			require.NoError(t, err)
		}

		// Filter by institution (partial match)
		filters := repository.EducationFilters{
			Institution: "MIT",
		}
		retrieved, err := repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, education := range retrieved {
			assert.Contains(t, education.Institution, "MIT")
		}
	})

	t.Run("GetEducation_FilterByStatus", func(t *testing.T) {
		testDB.CleanupTables(t)

		educations := []*models.Education{
			{
				Institution:           "University A",
				DegreeOrCertification: "Completed Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusCompleted,
			},
			{
				Institution:           "University B",
				DegreeOrCertification: "In Progress Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusInProgress,
			},
			{
				Institution:           "University C",
				DegreeOrCertification: "Planned Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusPlanned,
			},
		}

		for _, education := range educations {
			err := repo.CreateEducation(ctx, education)
			require.NoError(t, err)
		}

		// Filter by completed status
		filters := repository.EducationFilters{
			Status: models.EducationStatusCompleted,
		}
		retrieved, err := repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, models.EducationStatusCompleted, retrieved[0].Status)

		// Filter by in progress status
		filters.Status = models.EducationStatusInProgress
		retrieved, err = repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)
		assert.Equal(t, models.EducationStatusInProgress, retrieved[0].Status)
	})

	t.Run("GetEducation_FilterByFeatured", func(t *testing.T) {
		testDB.CleanupTables(t)

		educations := []*models.Education{
			{
				Institution:           "Featured University",
				DegreeOrCertification: "Featured Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusCompleted,
				IsFeatured:            true,
			},
			{
				Institution:           "Regular University",
				DegreeOrCertification: "Regular Degree",
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusCompleted,
				IsFeatured:            false,
			},
			{
				Institution:           "Another Featured",
				DegreeOrCertification: "Another Featured Degree",
				Type:                  models.EducationTypeCertification,
				Status:                models.EducationStatusCompleted,
				IsFeatured:            true,
			},
		}

		for _, education := range educations {
			err := repo.CreateEducation(ctx, education)
			require.NoError(t, err)
		}

		// Filter by featured
		filters := repository.EducationFilters{
			Featured: boolPtr(true),
		}
		retrieved, err := repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 2)
		for _, education := range retrieved {
			assert.True(t, education.IsFeatured)
		}
	})

	t.Run("GetEducation_Pagination", func(t *testing.T) {
		testDB.CleanupTables(t)

		// Create 5 education entries
		for i := 0; i < 5; i++ {
			education := &models.Education{
				Institution:           "University " + string(rune('A'+i)),
				DegreeOrCertification: "Degree " + string(rune('A'+i)),
				Type:                  models.EducationTypeEducation,
				Status:                models.EducationStatusCompleted,
				OrderIndex:            i,
			}
			err := repo.CreateEducation(ctx, education)
			require.NoError(t, err)
		}

		// Get first page
		filters := repository.EducationFilters{
			Limit:  2,
			Offset: 0,
		}
		page1, err := repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, page1, 2)

		// Get second page
		filters.Offset = 2
		page2, err := repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, page2, 2)

		// Ensure different results
		assert.NotEqual(t, page1[0].ID, page2[0].ID)
	})

	t.Run("GetEducationByType", func(t *testing.T) {
		testDB.CleanupTables(t)

		educations := []*models.Education{
			{Institution: "University", DegreeOrCertification: "Degree", Type: models.EducationTypeEducation, Status: models.EducationStatusCompleted},
			{Institution: "Cert Provider", DegreeOrCertification: "Cert", Type: models.EducationTypeCertification, Status: models.EducationStatusCompleted},
		}

		for _, education := range educations {
			err := repo.CreateEducation(ctx, education)
			require.NoError(t, err)
		}

		// Get education type
		educationList, err := repo.GetEducationByType(ctx, models.EducationTypeEducation)
		require.NoError(t, err)
		assert.Len(t, educationList, 1)
		assert.Equal(t, models.EducationTypeEducation, educationList[0].Type)

		// Get certification type
		certList, err := repo.GetEducationByType(ctx, models.EducationTypeCertification)
		require.NoError(t, err)
		assert.Len(t, certList, 1)
		assert.Equal(t, models.EducationTypeCertification, certList[0].Type)
	})

	t.Run("GetFeaturedEducation", func(t *testing.T) {
		testDB.CleanupTables(t)

		educations := []*models.Education{
			{Institution: "Featured Uni", DegreeOrCertification: "Featured Degree", Type: models.EducationTypeEducation, Status: models.EducationStatusCompleted, IsFeatured: true},
			{Institution: "Regular Uni", DegreeOrCertification: "Regular Degree", Type: models.EducationTypeEducation, Status: models.EducationStatusCompleted, IsFeatured: false},
			{Institution: "Another Featured", DegreeOrCertification: "Another Featured", Type: models.EducationTypeCertification, Status: models.EducationStatusCompleted, IsFeatured: true},
		}

		for _, education := range educations {
			err := repo.CreateEducation(ctx, education)
			require.NoError(t, err)
		}

		featured, err := repo.GetFeaturedEducation(ctx)
		require.NoError(t, err)
		assert.Len(t, featured, 2)
		for _, education := range featured {
			assert.True(t, education.IsFeatured)
		}
	})

	t.Run("UpdateEducation", func(t *testing.T) {
		testDB.CleanupTables(t)

		education := &models.Education{
			Institution:           "Original University",
			DegreeOrCertification: "Original Degree",
			Type:                  models.EducationTypeEducation,
			Status:                models.EducationStatusInProgress,
			IsFeatured:            false,
		}

		err := repo.CreateEducation(ctx, education)
		require.NoError(t, err)
		originalUpdatedAt := education.UpdatedAt

		time.Sleep(time.Millisecond * 10)

		// Update education
		education.Institution = "Updated University"
		education.DegreeOrCertification = "Updated Degree"
		education.Status = models.EducationStatusCompleted
		education.YearCompleted = intPtr(2024)
		education.IsFeatured = true

		err = repo.UpdateEducation(ctx, education)
		require.NoError(t, err)
		assert.True(t, education.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		filters := repository.EducationFilters{}
		retrieved, err := repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		require.Len(t, retrieved, 1)

		updated := retrieved[0]
		assert.Equal(t, "Updated University", updated.Institution)
		assert.Equal(t, "Updated Degree", updated.DegreeOrCertification)
		assert.Equal(t, models.EducationStatusCompleted, updated.Status)
		assert.Equal(t, 2024, *updated.YearCompleted)
		assert.True(t, updated.IsFeatured)
	})

	t.Run("UpdateEducation_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		education := &models.Education{
			ID:                    999,
			Institution:           "Non-existent",
			DegreeOrCertification: "Nobody",
		}

		err := repo.UpdateEducation(ctx, education)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "education with id 999 not found")
	})

	t.Run("DeleteEducation", func(t *testing.T) {
		testDB.CleanupTables(t)

		education := &models.Education{
			Institution:           "Delete Me University",
			DegreeOrCertification: "Temporary Degree",
			Type:                  models.EducationTypeEducation,
			Status:                models.EducationStatusCompleted,
		}

		err := repo.CreateEducation(ctx, education)
		require.NoError(t, err)

		// Verify it exists
		filters := repository.EducationFilters{}
		retrieved, err := repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)

		// Delete it
		err = repo.DeleteEducation(ctx, education.ID)
		require.NoError(t, err)

		// Verify it's gone
		retrieved, err = repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		assert.Len(t, retrieved, 0)
	})

	t.Run("DeleteEducation_NotFound", func(t *testing.T) {
		testDB.CleanupTables(t)

		err := repo.DeleteEducation(ctx, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "education with id 999 not found")
	})

	t.Run("EducationConstants_Validation", func(t *testing.T) {
		// Test education types
		types := models.ValidEducationTypes()
		assert.Contains(t, types, models.EducationTypeEducation)
		assert.Contains(t, types, models.EducationTypeCertification)
		assert.Len(t, types, 2)

		// Test education statuses
		statuses := models.ValidEducationStatuses()
		assert.Contains(t, statuses, models.EducationStatusCompleted)
		assert.Contains(t, statuses, models.EducationStatusInProgress)
		assert.Contains(t, statuses, models.EducationStatusPlanned)
		assert.Len(t, statuses, 3)
	})

	t.Run("Education_WithCredentials", func(t *testing.T) {
		testDB.CleanupTables(t)

		expiryDate := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)
		certification := &models.Education{
			Institution:           "Cloud Provider",
			DegreeOrCertification: "Cloud Architect Certification",
			Type:                  models.EducationTypeCertification,
			Status:                models.EducationStatusCompleted,
			CredentialID:          stringPtr("CERT-12345"),
			CredentialURL:         stringPtr("https://verify.provider.com/12345"),
			ExpiryDate:            &expiryDate,
		}

		err := repo.CreateEducation(ctx, certification)
		require.NoError(t, err)

		// Retrieve and verify credentials
		filters := repository.EducationFilters{Type: models.EducationTypeCertification}
		retrieved, err := repo.GetEducation(ctx, filters)
		require.NoError(t, err)
		require.Len(t, retrieved, 1)

		cert := retrieved[0]
		assert.Equal(t, "CERT-12345", *cert.CredentialID)
		assert.Equal(t, "https://verify.provider.com/12345", *cert.CredentialURL)
		assert.Equal(t, expiryDate.Unix(), cert.ExpiryDate.Unix())
	})
}