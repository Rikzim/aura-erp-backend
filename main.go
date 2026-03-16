package main

import (
	"log"
	"os"

	"aura-erp/backend/config"
	"aura-erp/backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Aura ERP API
// @version 1.0
// @description This is a powerful ERP backend for Aura.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.aura-erp.com/support
// @contact.email support@aura-erp.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /api

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize database connection
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDB()

	// Create Gin router
	router := gin.Default()

	// Setup CORS
	router.Use(corsMiddleware())

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Server running on http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
