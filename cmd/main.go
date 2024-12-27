package main

import (
	"identity-service/config"
	"identity-service/db"
	"identity-service/internal/initializer"
	"identity-service/internal/routes"
	"identity-service/pkg/utils"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load configurations
	config.LoadOAuthConfig()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Initialize database
	if err := db.Connect(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Initialize all components
	repos := initializer.InitRepositories()
	services := initializer.InitServices(repos)

	// Share the key manager with utils package
	utils.SetKeyManager(services.GetKeyManager())

	handlers := initializer.InitHandlers(services)

	// Setup router with middleware
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Register all routes
	routes.AuthRoutes(router, handlers.AuthHandler, handlers.OAuthHandler, services.GetKeyManager(), repos.UserRepo)
	routes.UserRoutes(router, handlers.UserHandler, services.GetKeyManager(), repos.UserRepo)
	routes.TenantRoutes(router, handlers.TenantHandler, services.GetKeyManager(), repos.UserRepo)
	routes.SecurityRoutes(router, handlers.SecurityHandler, services.GetKeyManager(), repos.UserRepo)

	// Start server
	port := ":4000"
	log.Printf("Server is running at http://localhost%s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
