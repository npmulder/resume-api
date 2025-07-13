// Package services implements the business logic for the Resume API.
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/npmulder/resume-api/internal/cache"
	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/repository"
)

// CachedResumeService is a decorator for ResumeService that adds caching
type CachedResumeService struct {
	service ResumeService
	cache   cache.Cache
	ttl     time.Duration
}

// NewCachedResumeService creates a new cached resume service
func NewCachedResumeService(service ResumeService, cache cache.Cache, ttl time.Duration) ResumeService {
	return &CachedResumeService{
		service: service,
		cache:   cache,
		ttl:     ttl,
	}
}

// GetProfile retrieves the user's profile, with caching
func (s *CachedResumeService) GetProfile(ctx context.Context) (*models.Profile, error) {
	cacheKey := "profile"
	var profile models.Profile

	// Try to get from cache first
	err := s.cache.Get(ctx, cacheKey, &profile)
	if err == nil {
		return &profile, nil
	}

	// If not in cache or error, get from service
	if err != cache.ErrCacheMiss {
		// Log the error but continue to fetch from service
		fmt.Printf("Cache error for profile: %v\n", err)
	}

	// Get from service
	result, err := s.service.GetProfile(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if err := s.cache.Set(ctx, cacheKey, result, s.ttl); err != nil {
		// Log the error but don't fail the request
		fmt.Printf("Failed to cache profile: %v\n", err)
	}

	return result, nil
}

// GetExperiences retrieves work experiences with optional filtering, with caching
func (s *CachedResumeService) GetExperiences(ctx context.Context, filters repository.ExperienceFilters) ([]*models.Experience, error) {
	// Create a cache key based on the filters
	cacheKey := fmt.Sprintf("experiences:%v:%v:%v:%v:%v",
		filters.Company, filters.Position, filters.IsCurrent, filters.Limit, filters.Offset)

	var experiences []*models.Experience

	// Try to get from cache first
	err := s.cache.Get(ctx, cacheKey, &experiences)
	if err == nil {
		return experiences, nil
	}

	// If not in cache or error, get from service
	if err != cache.ErrCacheMiss {
		fmt.Printf("Cache error for experiences: %v\n", err)
	}

	// Get from service
	experiences, err = s.service.GetExperiences(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if err := s.cache.Set(ctx, cacheKey, experiences, s.ttl); err != nil {
		fmt.Printf("Failed to cache experiences: %v\n", err)
	}

	return experiences, nil
}

// GetSkills retrieves skills with optional filtering, with caching
func (s *CachedResumeService) GetSkills(ctx context.Context, filters repository.SkillFilters) ([]*models.Skill, error) {
	// Create a cache key based on the filters
	cacheKey := fmt.Sprintf("skills:%v:%v:%v:%v",
		filters.Category, filters.Featured, filters.Limit, filters.Offset)

	var skills []*models.Skill

	// Try to get from cache first
	err := s.cache.Get(ctx, cacheKey, &skills)
	if err == nil {
		return skills, nil
	}

	// If not in cache or error, get from service
	if err != cache.ErrCacheMiss {
		fmt.Printf("Cache error for skills: %v\n", err)
	}

	// Get from service
	skills, err = s.service.GetSkills(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if err := s.cache.Set(ctx, cacheKey, skills, s.ttl); err != nil {
		fmt.Printf("Failed to cache skills: %v\n", err)
	}

	return skills, nil
}

// GetAchievements retrieves achievements with optional filtering, with caching
func (s *CachedResumeService) GetAchievements(ctx context.Context, filters repository.AchievementFilters) ([]*models.Achievement, error) {
	// Create a cache key based on the filters
	cacheKey := fmt.Sprintf("achievements:%v:%v:%v:%v:%v",
		filters.Year, filters.Category, filters.Featured, filters.Limit, filters.Offset)

	var achievements []*models.Achievement

	// Try to get from cache first
	err := s.cache.Get(ctx, cacheKey, &achievements)
	if err == nil {
		return achievements, nil
	}

	// If not in cache or error, get from service
	if err != cache.ErrCacheMiss {
		fmt.Printf("Cache error for achievements: %v\n", err)
	}

	// Get from service
	achievements, err = s.service.GetAchievements(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if err := s.cache.Set(ctx, cacheKey, achievements, s.ttl); err != nil {
		fmt.Printf("Failed to cache achievements: %v\n", err)
	}

	return achievements, nil
}

// GetEducation retrieves education entries with optional filtering, with caching
func (s *CachedResumeService) GetEducation(ctx context.Context, filters repository.EducationFilters) ([]*models.Education, error) {
	// Create a cache key based on the filters
	cacheKey := fmt.Sprintf("education:%v:%v:%v:%v",
		filters.Type, filters.Status, filters.Limit, filters.Offset)

	var education []*models.Education

	// Try to get from cache first
	err := s.cache.Get(ctx, cacheKey, &education)
	if err == nil {
		return education, nil
	}

	// If not in cache or error, get from service
	if err != cache.ErrCacheMiss {
		fmt.Printf("Cache error for education: %v\n", err)
	}

	// Get from service
	education, err = s.service.GetEducation(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if err := s.cache.Set(ctx, cacheKey, education, s.ttl); err != nil {
		fmt.Printf("Failed to cache education: %v\n", err)
	}

	return education, nil
}

// GetProjects retrieves projects with optional filtering, with caching
func (s *CachedResumeService) GetProjects(ctx context.Context, filters repository.ProjectFilters) ([]*models.Project, error) {
	// Create a cache key based on the filters
	cacheKey := fmt.Sprintf("projects:%v:%v:%v:%v:%v",
		filters.Status, filters.Technology, filters.Featured, filters.Limit, filters.Offset)

	var projects []*models.Project

	// Try to get from cache first
	err := s.cache.Get(ctx, cacheKey, &projects)
	if err == nil {
		return projects, nil
	}

	// If not in cache or error, get from service
	if err != cache.ErrCacheMiss {
		fmt.Printf("Cache error for projects: %v\n", err)
	}

	// Get from service
	projects, err = s.service.GetProjects(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if err := s.cache.Set(ctx, cacheKey, projects, s.ttl); err != nil {
		fmt.Printf("Failed to cache projects: %v\n", err)
	}

	return projects, nil
}
