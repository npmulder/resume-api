package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Data structures matching the JSON format
type SeedData struct {
	Profile      Profile      `json:"profile"`
	Experiences  []Experience `json:"experiences"`
	Skills       []Skill      `json:"skills"`
	Achievements []Achievement `json:"achievements"`
	Education    []Education   `json:"education"`
	Projects     []Project     `json:"projects"`
}

type Profile struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	LinkedIn string `json:"linkedin"`
	GitHub   string `json:"github"`
	Summary  string `json:"summary"`
}

type Experience struct {
	Company     string    `json:"company"`
	Position    string    `json:"position"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date"`
	Description string    `json:"description"`
	Highlights  []string  `json:"highlights"`
	Order       int       `json:"order"`
}

type Skill struct {
	Category string `json:"category"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	Order    int    `json:"order"`
	Featured bool   `json:"featured"`
}

type Achievement struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Impact      string `json:"impact"`
	Year        int    `json:"year"`
	Order       int    `json:"order"`
	Featured    bool   `json:"featured"`
}

type Education struct {
	Institution   string  `json:"institution"`
	Degree        string  `json:"degree"`
	Field         string  `json:"field"`
	YearCompleted *int    `json:"year_completed"`
	YearStarted   *int    `json:"year_started"`
	Description   string  `json:"description"`
	Type          string  `json:"type"`
	Status        string  `json:"status"`
	CredentialID  string  `json:"credential_id"`
	CredentialURL string  `json:"credential_url"`
	Order         int     `json:"order"`
	Featured      bool    `json:"featured"`
}

type Project struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ShortDescription string    `json:"short_description"`
	Technologies     []string  `json:"technologies"`
	GitHubURL        string    `json:"github_url"`
	DemoURL          *string   `json:"demo_url"`
	StartDate        string    `json:"start_date"`
	EndDate          *string   `json:"end_date"`
	Status           string    `json:"status"`
	IsFeatured       bool      `json:"is_featured"`
	Order            int       `json:"order"`
	KeyFeatures      []string  `json:"key_features"`
}

func main() {
	// Get database URL from environment or use default
	dbURL := getEnv("DATABASE_URL", "postgres://dev:devpass@localhost:5432/resume_api_dev?sslmode=disable")

	// Determine seed data file path
	seedFile := getEnv("SEED_FILE", "scripts/seed-data.json")
	
	// Check if seed file exists, fall back to example if not
	if _, err := os.Stat(seedFile); os.IsNotExist(err) {
		fmt.Printf("Seed file %s not found, using example data from scripts/seed-data.example.json\n", seedFile)
		seedFile = "scripts/seed-data.example.json"
	}

	// Load seed data from JSON file
	seedData, err := loadSeedData(seedFile)
	if err != nil {
		log.Fatal("Failed to load seed data:", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Printf("Seeding database with data from %s...\n", seedFile)

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}
	defer tx.Rollback()

	// Seed all data
	if err := seedProfile(tx, seedData.Profile); err != nil {
		log.Fatal("Failed to seed profile:", err)
	}

	if err := seedExperiences(tx, seedData.Experiences); err != nil {
		log.Fatal("Failed to seed experiences:", err)
	}

	if err := seedSkills(tx, seedData.Skills); err != nil {
		log.Fatal("Failed to seed skills:", err)
	}

	if err := seedAchievements(tx, seedData.Achievements); err != nil {
		log.Fatal("Failed to seed achievements:", err)
	}

	if err := seedEducation(tx, seedData.Education); err != nil {
		log.Fatal("Failed to seed education:", err)
	}

	if err := seedProjects(tx, seedData.Projects); err != nil {
		log.Fatal("Failed to seed projects:", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}

	fmt.Println("Database seeded successfully!")
}

func loadSeedData(filename string) (*SeedData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read seed file: %w", err)
	}

	var seedData SeedData
	if err := json.Unmarshal(data, &seedData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &seedData, nil
}

func seedProfile(tx *sql.Tx, profile Profile) error {
	query := `
		INSERT INTO profiles (name, title, email, phone, location, linkedin, github, summary)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (email) DO UPDATE SET
			name = EXCLUDED.name,
			title = EXCLUDED.title,
			phone = EXCLUDED.phone,
			location = EXCLUDED.location,
			linkedin = EXCLUDED.linkedin,
			github = EXCLUDED.github,
			summary = EXCLUDED.summary,
			updated_at = CURRENT_TIMESTAMP`

	_, err := tx.Exec(query,
		profile.Name,
		profile.Title,
		profile.Email,
		profile.Phone,
		profile.Location,
		profile.LinkedIn,
		profile.GitHub,
		profile.Summary,
	)
	return err
}

func seedExperiences(tx *sql.Tx, experiences []Experience) error {
	// Clear existing experiences
	if _, err := tx.Exec("DELETE FROM experiences"); err != nil {
		return err
	}

	query := `
		INSERT INTO experiences (company, position, start_date, end_date, description, highlights, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, exp := range experiences {
		var endDate *time.Time
		if exp.EndDate != nil {
			ed, err := time.Parse("2006-01-02", *exp.EndDate)
			if err != nil {
				return fmt.Errorf("failed to parse end date for %s: %w", exp.Company, err)
			}
			endDate = &ed
		}

		startDate, err := time.Parse("2006-01-02", exp.StartDate)
		if err != nil {
			return fmt.Errorf("failed to parse start date for %s: %w", exp.Company, err)
		}

		_, err = tx.Exec(query,
			exp.Company,
			exp.Position,
			startDate,
			endDate,
			exp.Description,
			exp.Highlights,
			exp.Order,
		)
		if err != nil {
			return fmt.Errorf("failed to insert experience for %s: %w", exp.Company, err)
		}
	}

	return nil
}

func seedSkills(tx *sql.Tx, skills []Skill) error {
	// Clear existing skills
	if _, err := tx.Exec("DELETE FROM skills"); err != nil {
		return err
	}

	query := `
		INSERT INTO skills (category, name, level, order_index, is_featured)
		VALUES ($1, $2, $3, $4, $5)`

	for _, skill := range skills {
		_, err := tx.Exec(query,
			skill.Category,
			skill.Name,
			skill.Level,
			skill.Order,
			skill.Featured,
		)
		if err != nil {
			return fmt.Errorf("failed to insert skill %s: %w", skill.Name, err)
		}
	}

	return nil
}

func seedAchievements(tx *sql.Tx, achievements []Achievement) error {
	// Clear existing achievements
	if _, err := tx.Exec("DELETE FROM achievements"); err != nil {
		return err
	}

	query := `
		INSERT INTO achievements (title, description, category, impact_metric, year_achieved, order_index, is_featured)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, achievement := range achievements {
		_, err := tx.Exec(query,
			achievement.Title,
			achievement.Description,
			achievement.Category,
			achievement.Impact,
			achievement.Year,
			achievement.Order,
			achievement.Featured,
		)
		if err != nil {
			return fmt.Errorf("failed to insert achievement %s: %w", achievement.Title, err)
		}
	}

	return nil
}

func seedEducation(tx *sql.Tx, education []Education) error {
	// Clear existing education
	if _, err := tx.Exec("DELETE FROM education"); err != nil {
		return err
	}

	query := `
		INSERT INTO education (institution, degree_or_certification, field_of_study, year_completed, year_started, 
			description, type, status, credential_id, credential_url, order_index, is_featured)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, edu := range education {
		_, err := tx.Exec(query,
			edu.Institution,
			edu.Degree,
			edu.Field,
			edu.YearCompleted,
			edu.YearStarted,
			edu.Description,
			edu.Type,
			edu.Status,
			edu.CredentialID,
			edu.CredentialURL,
			edu.Order,
			edu.Featured,
		)
		if err != nil {
			return fmt.Errorf("failed to insert education record for %s: %w", edu.Institution, err)
		}
	}

	return nil
}

func seedProjects(tx *sql.Tx, projects []Project) error {
	// Clear existing projects
	if _, err := tx.Exec("DELETE FROM projects"); err != nil {
		return err
	}

	query := `
		INSERT INTO projects (name, description, short_description, technologies, github_url, demo_url, 
			start_date, end_date, status, is_featured, order_index, key_features)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, project := range projects {
		// Convert technologies slice to JSON
		techJSON, err := json.Marshal(project.Technologies)
		if err != nil {
			return fmt.Errorf("failed to marshal technologies for %s: %w", project.Name, err)
		}

		var startDate time.Time
		if project.StartDate != "" {
			startDate, err = time.Parse("2006-01-02", project.StartDate)
			if err != nil {
				return fmt.Errorf("failed to parse start date for %s: %w", project.Name, err)
			}
		}

		var endDate *time.Time
		if project.EndDate != nil {
			ed, err := time.Parse("2006-01-02", *project.EndDate)
			if err != nil {
				return fmt.Errorf("failed to parse end date for %s: %w", project.Name, err)
			}
			endDate = &ed
		}

		_, err = tx.Exec(query,
			project.Name,
			project.Description,
			project.ShortDescription,
			string(techJSON),
			project.GitHubURL,
			project.DemoURL,
			startDate,
			endDate,
			project.Status,
			project.IsFeatured,
			project.Order,
			project.KeyFeatures,
		)
		if err != nil {
			return fmt.Errorf("failed to insert project %s: %w", project.Name, err)
		}
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}