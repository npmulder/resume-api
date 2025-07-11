package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		databaseURL = flag.String("database-url", getEnv("DATABASE_URL", "postgres://dev:devpass@localhost:5432/resume_api_dev?sslmode=disable"), "Database URL")
		direction   = flag.String("direction", "up", "Migration direction: up or down")
		steps       = flag.Int("steps", 0, "Number of migration steps (0 means all)")
	)
	flag.Parse()

	if len(flag.Args()) > 0 {
		*direction = flag.Args()[0]
		if len(flag.Args()) > 1 {
			fmt.Sscanf(flag.Args()[1], "%d", steps)
		}
	}

	// Create migration instance
	m, err := migrate.New(
		"file://migrations",
		*databaseURL,
	)
	if err != nil {
		log.Fatal("Failed to create migration instance:", err)
	}
	defer m.Close()

	// Execute migration based on direction
	switch *direction {
	case "up":
		if *steps == 0 {
			err = m.Up()
			if err == migrate.ErrNoChange {
				fmt.Println("No migrations to apply")
				return
			}
		} else {
			err = m.Steps(*steps)
		}
	case "down":
		if *steps == 0 {
			err = m.Down()
		} else {
			err = m.Steps(-*steps)
		}
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get migration version:", err)
		}
		fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)
		return
	case "force":
		if *steps == 0 {
			log.Fatal("Force requires a version number")
		}
		err = m.Force(*steps)
		if err != nil {
			log.Fatal("Failed to force migration version:", err)
		}
		fmt.Printf("Forced migration version to: %d\n", *steps)
		return
	default:
		log.Fatal("Invalid direction. Use 'up', 'down', 'version', or 'force'")
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply")
		} else {
			log.Fatal("Migration failed:", err)
		}
	} else {
		fmt.Printf("Migration %s completed successfully\n", *direction)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}