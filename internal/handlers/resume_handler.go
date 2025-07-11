package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/npmulder/resume-api/internal/repository"
	"github.com/npmulder/resume-api/internal/services"
)

// ResumeHandler handles the HTTP requests for the resume data.
type ResumeHandler struct {
	service services.ResumeService
}

// NewResumeHandler creates a new ResumeHandler.
func NewResumeHandler(service services.ResumeService) *ResumeHandler {
	return &ResumeHandler{service: service}
}

// GetProfile handles the request to get the user's profile.
func (h *ResumeHandler) GetProfile(c *gin.Context) {
	profile, err := h.service.GetProfile(c.Request.Context())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// GetExperiences handles the request to get the user's work experiences.
func (h *ResumeHandler) GetExperiences(c *gin.Context) {
	var filters repository.ExperienceFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	experiences, err := h.service.GetExperiences(c.Request.Context(), filters)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, experiences)
}

// GetSkills handles the request to get the user's skills.
func (h *ResumeHandler) GetSkills(c *gin.Context) {
	var filters repository.SkillFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skills, err := h.service.GetSkills(c.Request.Context(), filters)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, skills)
}

// GetAchievements handles the request to get the user's achievements.
func (h *ResumeHandler) GetAchievements(c *gin.Context) {
	var filters repository.AchievementFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	achievements, err := h.service.GetAchievements(c.Request.Context(), filters)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, achievements)
}

// GetEducation handles the request to get the user's education.
func (h *ResumeHandler) GetEducation(c *gin.Context) {
	var filters repository.EducationFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	education, err := h.service.GetEducation(c.Request.Context(), filters)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, education)
}

// GetProjects handles the request to get the user's projects.
func (h *ResumeHandler) GetProjects(c *gin.Context) {
	var filters repository.ProjectFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projects, err := h.service.GetProjects(c.Request.Context(), filters)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, projects)
}
