package main

import (
	"log"
	"os"

	"Oratio/routes"
	"Oratio/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Supabase DB connection
	services.InitDatabase()

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	routes.RegisterRoutes(r)

	log.Printf("Starting server on port %s...", port)
	r.Run(":" + port)
}
