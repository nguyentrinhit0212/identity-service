package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"

	// postgres driver for database migrations
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// file source driver for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func HealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "identity-service",
		})
	})

	router.GET("/health/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ready",
			"service": "identity-service",
		})
	})

	router.GET("/health/migration", func(c *gin.Context) {
		// Create database URL
		dbURL := fmt.Sprintf(
			"postgres://%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_NAME"),
		)

		// Create migration instance
		m, err := migrate.New("file://db/migrations", dbURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to create migration instance",
				"error":   err.Error(),
			})
			return
		}
		defer m.Close()

		// Get migration version
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to get migration version",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "identity-service",
			"version": version,
			"dirty":   dirty,
			"message": "Migration status retrieved successfully",
			"has_run": err != migrate.ErrNilVersion,
		})
	})
}
