// @title Resume API
// @version 1.0
// @description API for resume data including profile, experiences, skills, achievements, education, and projects
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Import generated docs
	_ "github.com/npmulder/resume-api/docs"
	"github.com/npmulder/resume-api/internal/cache"
	"github.com/npmulder/resume-api/internal/config"
	"github.com/npmulder/resume-api/internal/database"
	"github.com/npmulder/resume-api/internal/handlers"
	"github.com/npmulder/resume-api/internal/middleware"
	"github.com/npmulder/resume-api/internal/repository"
	"github.com/npmulder/resume-api/internal/repository/postgres"
	"github.com/npmulder/resume-api/internal/services"
	"github.com/npmulder/resume-api/internal/tracing"
	"github.com/npmulder/resume-api/internal/versioning"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Set up logger
	logLevel := new(slog.LevelVar)
	if err := logLevel.UnmarshalText([]byte(cfg.Logging.Level)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse log level: %v\n", err)
		os.Exit(1)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// Initialize tracing
	tracer, err := tracing.NewTracer(context.Background(), &cfg.Telemetry, logger)
	if err != nil {
		logger.Error("failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tracer.Shutdown(context.Background()); err != nil {
			logger.Error("failed to shutdown tracer", "error", err)
		}
	}()

	// Establish database connection
	db, err := database.New(context.Background(), &cfg.Database, logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection established")

	// Initialize repositories
	profileRepo := postgres.NewProfileRepository(db.Pool())
	experienceRepo := postgres.NewExperienceRepository(db.Pool())
	skillRepo := postgres.NewSkillRepository(db.Pool())
	achievementRepo := postgres.NewAchievementRepository(db.Pool())
	educationRepo := postgres.NewEducationRepository(db.Pool())
	projectRepo := postgres.NewProjectRepository(db.Pool())

	repos := repository.Repositories{
		Profile:     profileRepo,
		Experience:  experienceRepo,
		Skill:       skillRepo,
		Achievement: achievementRepo,
		Education:   educationRepo,
		Project:     projectRepo,
	}

	// Initialize cache
	cacheClient, err := cache.New(&cfg.Redis)
	if err != nil {
		if cfg.Redis.Enabled {
			logger.Error("failed to initialize Redis cache", "error", err)
			os.Exit(1)
		}
		logger.Info("Redis cache is disabled, using no-op cache")
	} else {
		logger.Info("Redis cache initialized successfully")
		defer cacheClient.Close()
	}

	// Initialize services
	baseResumeService := services.NewResumeService(repos)
	resumeService := services.NewCachedResumeService(baseResumeService, cacheClient, cfg.Redis.TTL)

	// Initialize handlers
	resumeHandler := handlers.NewResumeHandler(resumeService)

	// Set up Gin router
	router := gin.New()

	// Register middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware(logger))
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.TimeoutMiddleware(cfg.Server.RequestTimeout, logger))
	router.Use(middleware.MetricsMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.InputValidationMiddleware())
	router.Use(middleware.RateLimiterMiddleware(middleware.DefaultRateLimiterConfig()))
	router.Use(middleware.TracingMiddleware(tracer))

	// Add version negotiation middleware
	router.Use(versioning.VersionNegotiationMiddleware(versioning.DefaultVersionNegotiationOptions()))

	// Define routes
	router.GET("/health", handlers.HealthCheck)
	router.GET("/metrics", handlers.MetricsHandler())

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create versioned router
	versionedRouter := versioning.NewRouter(router)

	// Register API routes for v1
	v1 := versionedRouter.Group(versioning.V1)
	{
		v1.GET("/profile", resumeHandler.GetProfile)
		v1.GET("/experiences", resumeHandler.GetExperiences)
		v1.GET("/skills", resumeHandler.GetSkills)
		v1.GET("/achievements", resumeHandler.GetAchievements)
		v1.GET("/education", resumeHandler.GetEducation)
		v1.GET("/projects", resumeHandler.GetProjects)
	}

	// Create and start HTTP server
	srv := &http.Server{
		Addr:         cfg.Server.ServerAddress(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Info("starting server", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Implement graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulStop)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server exited gracefully")
}
