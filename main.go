package main

import (
	"github.com/ftilie/go-booking-api/database"
	"github.com/ftilie/go-booking-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// This is the entry point of the application.
	database.InitDB() // Initialize the database connection
	server := gin.Default()

	// Register the routes
	routes.RegisterRoutes(server)

	// Start application server
	server.Run(":8080") // Start the server on port 8080
}
