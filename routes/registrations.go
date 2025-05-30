package routes

import (
	"net/http"
	"strconv"

	"github.com/ftilie/go-booking-api/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	// This function will handle attendee registration for an event
	userId := context.GetInt64("userId")
	stringEventId := context.Param("eventId")
	if stringEventId == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Event ID is required!"})
		return
	}
	eventId, err := strconv.ParseInt(stringEventId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse event Id!"})
		return
	}

	event, err := models.GetEvent(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve event from the database!"})
		return
	}

	err = event.RegisterForEvent(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register for the event!"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Successfully registered for the event!"})
}

func cancelRegistration(context *gin.Context) {
	// This function will handle attendee cancellation for an event
	userId := context.GetInt64("userId")
	stringEventId := context.Param("eventId")
	if stringEventId == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Event ID is required!"})
		return
	}
	eventId, err := strconv.ParseInt(stringEventId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse event Id!"})
		return
	}

	event, err := models.GetEvent(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve event from the database!"})
		return
	}

	err = event.CancelRegistration(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to cancel registration for the event!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Successfully cancelled registration for the event!"})
}
