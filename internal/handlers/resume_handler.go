package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/npmulder/resume-api/internal/repository"
	"github.com/npmulder/resume-api/internal/services"
	"github.com/npmulder/resume-api/internal/utils"
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
// @Summary Get user profile
// @Description Retrieve the user's personal information and summary
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} models.Profile
// @Failure 404 {object} models.APIError "Not found"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /api/v1/profile [get]
func (h *ResumeHandler) GetProfile(c *gin.Context) {
	profile, err := h.service.GetProfile(c.Request.Context())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.NotFound(c, "Profile not found")
			return
		}
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

// GetExperiences handles the request to get the user's work experiences.
// @Summary Get work experiences
// @Description Retrieve the user's work history and professional experiences with optional filtering
// @Tags experiences
// @Accept json
// @Produce json
// @Param company query string false "Filter by company name"
// @Param position query string false "Filter by position title"
// @Param date_from query string false "Filter by start date (ISO format)"
// @Param date_to query string false "Filter by end date (ISO format)"
// @Param is_current query boolean false "Filter for current positions"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.Experience
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /api/v1/experiences [get]
func (h *ResumeHandler) GetExperiences(c *gin.Context) {
	var filters repository.ExperienceFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.ValidationError(c, "Invalid query parameters", err.Error())
		return
	}

	experiences, err := h.service.GetExperiences(c.Request.Context(), filters)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.NotFound(c, "No experiences found matching the criteria")
			return
		}
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, experiences)
}

// GetSkills handles the request to get the user's skills.
// @Summary Get skills
// @Description Retrieve the user's technical and soft skills with optional filtering
// @Tags skills
// @Accept json
// @Produce json
// @Param category query string false "Filter by skill category"
// @Param level query string false "Filter by skill level (beginner, intermediate, advanced, expert)"
// @Param featured query boolean false "Filter for featured skills"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.Skill
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /api/v1/skills [get]
func (h *ResumeHandler) GetSkills(c *gin.Context) {
	var filters repository.SkillFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.ValidationError(c, "Invalid query parameters", err.Error())
		return
	}

	skills, err := h.service.GetSkills(c.Request.Context(), filters)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.NotFound(c, "No skills found matching the criteria")
			return
		}
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, skills)
}

// GetAchievements handles the request to get the user's achievements.
// @Summary Get achievements
// @Description Retrieve the user's key accomplishments and achievements with optional filtering
// @Tags achievements
// @Accept json
// @Produce json
// @Param category query string false "Filter by achievement category"
// @Param year query int false "Filter by year achieved"
// @Param featured query boolean false "Filter for featured achievements"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.Achievement
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /api/v1/achievements [get]
func (h *ResumeHandler) GetAchievements(c *gin.Context) {
	var filters repository.AchievementFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.ValidationError(c, "Invalid query parameters", err.Error())
		return
	}

	achievements, err := h.service.GetAchievements(c.Request.Context(), filters)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.NotFound(c, "No achievements found matching the criteria")
			return
		}
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, achievements)
}

// GetEducation handles the request to get the user's education.
// @Summary Get education
// @Description Retrieve the user's education and certifications with optional filtering
// @Tags education
// @Accept json
// @Produce json
// @Param type query string false "Filter by type (education or certification)"
// @Param institution query string false "Filter by institution name"
// @Param status query string false "Filter by status (completed, in_progress, planned)"
// @Param featured query boolean false "Filter for featured education entries"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.Education
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /api/v1/education [get]
func (h *ResumeHandler) GetEducation(c *gin.Context) {
	var filters repository.EducationFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.ValidationError(c, "Invalid query parameters", err.Error())
		return
	}

	education, err := h.service.GetEducation(c.Request.Context(), filters)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.NotFound(c, "No education records found matching the criteria")
			return
		}
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, education)
}

// GetProjects handles the request to get the user's projects.
// @Summary Get projects
// @Description Retrieve the user's notable projects and implementations with optional filtering
// @Tags projects
// @Accept json
// @Produce json
// @Param status query string false "Filter by status (active, completed, archived, planned)"
// @Param technology query string false "Filter by technology used"
// @Param featured query boolean false "Filter for featured projects"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.Project
// @Failure 400 {object} models.APIError "Bad request"
// @Failure 404 {object} models.APIError "Not found"
// @Failure 500 {object} models.APIError "Internal server error"
// @Router /api/v1/projects [get]
func (h *ResumeHandler) GetProjects(c *gin.Context) {
	var filters repository.ProjectFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.ValidationError(c, "Invalid query parameters", err.Error())
		return
	}

	projects, err := h.service.GetProjects(c.Request.Context(), filters)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.NotFound(c, "No projects found matching the criteria")
			return
		}
		utils.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, projects)
}
