package main

import (
	"log"

	"vocabulary-app-be/internal/auth"
	"vocabulary-app-be/internal/vocab"
	"vocabulary-app-be/pkg/config"
	"vocabulary-app-be/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Gin router
	router := gin.Default()

	// Initialize auth module
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authController := auth.NewController(authService)
	auth.RegisterRoutes(router, authController)

	// Initialize vocab module
	vocabRepo := vocab.NewRepository(db)
	vocabService := vocab.NewService(vocabRepo)
	vocabController := vocab.NewController(vocabService)
	vocab.RegisterRoutes(router, vocabController)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
