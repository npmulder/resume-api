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

	"github.com/npmulder/resume-api/internal/config"
	"github.com/npmulder/resume-api/internal/database"
	"github.com/npmulder/resume-api/internal/handlers"
	"github.com/npmulder/resume-api/internal/middleware"
	"github.com/npmulder/resume-api/internal/repository"
	"github.com/npmulder/resume-api/internal/repository/postgres"
	"github.com/npmulder/resume-api/internal/services"
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

	// Initialize services
	resumeService := services.NewResumeService(repos)

	// Initialize handlers
	resumeHandler := handlers.NewResumeHandler(resumeService)

	// Set up Gin router
	router := gin.New()

	// Register middleware
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.TimeoutMiddleware(cfg.Server.RequestTimeout, logger))

	// Define routes
	router.GET("/health", handlers.HealthCheck)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/profile", resumeHandler.GetProfile)
		apiV1.GET("/experiences", resumeHandler.GetExperiences)
		apiV1.GET("/skills", resumeHandler.GetSkills)
		apiV1.GET("/achievements", resumeHandler.GetAchievements)
		apiV1.GET("/education", resumeHandler.GetEducation)
		apiV1.GET("/projects", resumeHandler.GetProjects)
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
