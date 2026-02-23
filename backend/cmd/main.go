package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"todo-fullstack/backend/config"
	"todo-fullstack/backend/controllers"
	"todo-fullstack/backend/middleware"
	"todo-fullstack/backend/models"
	"todo-fullstack/backend/repository"
	"todo-fullstack/backend/routes"
	"todo-fullstack/backend/services"
)

func main() {
	// Load environment variables (optional - Docker Compose sets them directly)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	db := config.InitDB()

	// Auto migrate models
	err := db.AutoMigrate(&models.Todo{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	// Initialize repository, service, and controller
	todoRepository := repository.NewTodoRepository(db)
	todoService := services.NewTodoService(todoRepository)
	todoController := controllers.NewTodoController(todoService)

	// Setup Gin router
	r := gin.Default()

	// CORS Middleware
	r.Use(middleware.CORSMiddleware())

	// Logging Middleware
	r.Use(middleware.LoggingMiddleware())

	// Setup routes
	routes.SetupTodoRoutes(r, todoController)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
