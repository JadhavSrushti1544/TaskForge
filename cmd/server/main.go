package main

import (
	"log"
	"task-queue/config"
	"task-queue/handlers"
	"task-queue/middleware"
	"task-queue/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router := gin.Default()

	// Initialize repositories and handlers
	userRepo := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo)

	// Public routes
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Protected routes (demo)
	protected := router.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/me", func(c *gin.Context) {
			userID := c.GetInt("user_id")
			user, _ := userRepo.GetUserByID(userID)
			c.JSON(200, user)
		})
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
