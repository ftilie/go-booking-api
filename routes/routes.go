package routes

import (
	"github.com/ftilie/go-booking-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/events").Use(middlewares.Authenticate) // Create a group for authenticated routes

	// Register the routes for the events
	authenticated.GET("/", getEvents)
	authenticated.GET("/:eventId", getEvent)
	authenticated.POST("/", createEvent)
	authenticated.PUT("/:eventId", updateEvent)
	authenticated.DELETE("/:eventId", deleteEvent)

	// Register the routes for the users
	server.POST("/signup", signup)
	server.POST("/login", login)

	// Register the routes for the bookings
	authenticated.POST("/:eventId/registration", registerForEvent)
	authenticated.DELETE("/:eventId/registration", cancelRegistration)

}
