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
// @Response 200 {object} models.Profile "Example response" {"id":1,"name":"John Doe","title":"Senior Software Engineer","email":"john.doe@example.com","phone":"+1-555-123-4567","location":"San Francisco, CA","linkedin":"https://linkedin.com/in/johndoe","github":"https://github.com/johndoe","summary":"Experienced software engineer with a passion for building scalable applications","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}
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
// @Response 200 {array} models.Experience "Example response" [{"id":1,"company":"Tech Innovations Inc.","position":"Senior Software Engineer","start_date":"2020-01-01T00:00:00Z","end_date":null,"description":"Led development of cloud-native applications","highlights":["Implemented CI/CD pipeline","Reduced deployment time by 50%","Mentored junior developers"],"order_index":1,"is_current":true,"location":"San Francisco, CA","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},{"id":2,"company":"Digital Solutions LLC","position":"Software Developer","start_date":"2017-06-01T00:00:00Z","end_date":"2019-12-31T00:00:00Z","description":"Worked on backend services for e-commerce platform","highlights":["Developed RESTful APIs","Optimized database queries","Implemented payment processing integration"],"order_index":2,"is_current":false,"location":"New York, NY","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}]
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
// @Response 200 {array} models.Skill "Example response" [{"id":1,"category":"Languages","name":"Go","level":"advanced","years_experience":5,"order_index":1,"is_featured":true,"description":"Proficient in Go development including concurrency patterns and standard library","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},{"id":2,"category":"Frameworks","name":"React","level":"intermediate","years_experience":3,"order_index":2,"is_featured":true,"description":"Experience with React and Redux for frontend development","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},{"id":3,"category":"Tools","name":"Docker","level":"expert","years_experience":6,"order_index":3,"is_featured":true,"description":"Expert in containerization and orchestration with Docker and Kubernetes","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}]
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
// @Response 200 {array} models.Achievement "Example response" [{"id":1,"title":"Performance Optimization Award","description":"Recognized for optimizing application performance by 40%","category":"performance","impact_metric":"40% reduction in response time","year_achieved":2022,"order_index":1,"is_featured":true,"created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},{"id":2,"title":"Security Excellence","description":"Identified and fixed critical security vulnerabilities","category":"security","impact_metric":"Prevented potential data breach affecting 10,000+ users","year_achieved":2021,"order_index":2,"is_featured":true,"created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},{"id":3,"title":"Team Leadership Award","description":"Led cross-functional team to successful product launch","category":"leadership","impact_metric":"Delivered project 2 weeks ahead of schedule","year_achieved":2020,"order_index":3,"is_featured":false,"created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}]
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
// @Response 200 {array} models.Education "Example response" [{"id":1,"institution":"Stanford University","degree_or_certification":"Master of Science","field_of_study":"Computer Science","year_completed":2018,"year_started":2016,"description":"Specialized in Artificial Intelligence and Machine Learning","type":"education","status":"completed","order_index":1,"is_featured":true,"degree_title":"Master of Science in Computer Science","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},{"id":2,"institution":"AWS","degree_or_certification":"AWS Certified Solutions Architect","field_of_study":"Cloud Architecture","year_completed":2021,"year_started":2021,"description":"Professional certification for designing distributed systems on AWS","type":"certification","status":"completed","credential_id":"AWS-CSA-123456","credential_url":"https://aws.amazon.com/verification","expiry_date":"2024-01-01T00:00:00Z","order_index":2,"is_featured":true,"degree_title":"AWS Certified Solutions Architect","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},{"id":3,"institution":"University of California, Berkeley","degree_or_certification":"PhD","field_of_study":"Computer Science","year_started":2022,"description":"Research focus on distributed systems and cloud computing","type":"education","status":"in_progress","order_index":3,"is_featured":false,"degree_title":"PhD in Computer Science","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}]
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
// @Response 200 {array} models.Project "Example response" [{"id":1,"name":"Cloud-Native Resume API","description":"RESTful API for resume data with caching and metrics","short_description":"Resume API with advanced features","technologies":["Go","PostgreSQL","Docker","Redis"],"github_url":"https://github.com/username/resume-api","demo_url":"https://api.example.com","start_date":"2022-06-01T00:00:00Z","end_date":null,"status":"active","is_featured":true,"order_index":1,"key_features":["OpenAPI documentation","Redis caching","Prometheus metrics","Distributed tracing"],"created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},{"id":2,"name":"E-commerce Platform","description":"Full-stack e-commerce solution with payment processing","short_description":"Complete e-commerce solution","technologies":["React","Node.js","MongoDB","Stripe"],"github_url":"https://github.com/username/ecommerce","demo_url":"https://shop.example.com","start_date":"2021-01-01T00:00:00Z","end_date":"2021-12-31T00:00:00Z","status":"completed","is_featured":true,"order_index":2,"key_features":["User authentication","Product catalog","Shopping cart","Payment processing","Order tracking"],"created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"},{"id":3,"name":"AI-powered Content Analyzer","description":"Tool for analyzing and categorizing text content using NLP","short_description":"NLP-based content analysis tool","technologies":["Python","TensorFlow","Flask","AWS"],"github_url":null,"demo_url":null,"start_date":"2023-01-01T00:00:00Z","end_date":null,"status":"planned","is_featured":false,"order_index":3,"key_features":["Sentiment analysis","Topic classification","Content summarization","Language detection"],"created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}]
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
