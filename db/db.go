package db

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Connect establishes connection to the database
func Connect() error {
	// Validate required environment variables
	required := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			return fmt.Errorf("missing required environment variable: %s", env)
		}
	}

	// Extract host and port from DB_HOST
	host, port, err := net.SplitHostPort(os.Getenv("DB_HOST"))
	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			// If port is missing, use the default PostgreSQL port
			host = os.Getenv("DB_HOST")
			port = "5432"
		} else {
			return fmt.Errorf("invalid DB_HOST format: %w", err)
		}
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return fmt.Errorf("database connection failed: %w", err)
	}

	log.Printf("Successfully connected to database at %s", os.Getenv("DB_HOST"))
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	if db == nil {
		panic("Database connection not initialized. Call Connect() first")
	}
	return db
}
