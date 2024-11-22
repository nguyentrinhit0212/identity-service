package main

import (
	"identity-service/config"
	"identity-service/db"
	"identity-service/internal/initializer"
	"identity-service/internal/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	config.LoadOAuthConfig()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	if err := db.Connect(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	repos := initializer.InitRepositories()

    // Step 5: Initialize services
    services := initializer.InitServices(repos)

    // Step 6: Initialize handlers
    handlers := initializer.InitHandlers(services)
	
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // Allow cookies if needed
	}))

	routes.AuthRoutes(router, handlers.AuthHandler)

	log.Println("Server is running at http://localhost:4000")
	if err := router.Run(":4000"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
