package routes

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ftilie/go-booking-api/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	// This function will handle retrieving all events
	events, err := models.GetEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve events from the database!"})
		return
	}
	context.JSON(http.StatusOK, events)
}

func getEvent(context *gin.Context) {
	// This function will handle retrieving a specific event by its ID``
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
	if event == nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found!"})
		return
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
	// This function will handle creating a new event
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input could not be parsed!"})
		return
	}

	event.Organizer = context.GetInt64("userId") // Get the user ID from the context set by the authentication middleware
	event.CreatedAt = time.Now()

	err = event.CreateEvent()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create event in the database!"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created successfully!", "event": event})
}

func updateEvent(context *gin.Context) {
	// This function will handle updating an existing event
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

	if event.EndTime.Before(time.Now()) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Cannot update an event that has already ended!"})
		return
	}

	if event.DeletedAt != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Cannot update an event that was deleted!"})
		return
	}

	if context.GetInt64("userId") != event.Organizer {
		context.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to update this event!"})
		return
	}

	// Support for partial patching
	var input map[string]interface{}
	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON!"})
		return
	}
	eventValue := reflect.ValueOf(event).Elem()
	eventType := reflect.TypeOf(event).Elem()

	for key, val := range input {
		for i := 0; i < eventType.NumField(); i++ {
			field := eventType.Field(i)
			fieldName := field.Name

			// Match field names case-insensitively
			if strings.EqualFold(fieldName, key) {
				fieldValue := eventValue.FieldByName(fieldName)
				if !fieldValue.CanSet() {
					continue
				}

				// Type conversion depending on the field type
				switch fieldValue.Kind() {
				case reflect.String:
					if str, ok := val.(string); ok {
						fieldValue.SetString(str)
					}
				case reflect.Int64:
					if f, ok := val.(float64); ok { // JSON numbers are float64
						fieldValue.SetInt(int64(f))
					}
				case reflect.Slice:
					if s, ok := val.([]interface{}); ok {
						var slice []int64
						for _, item := range s {
							if num, ok := item.(float64); ok {
								slice = append(slice, int64(num))
							}
						}
						fieldValue.Set(reflect.ValueOf(slice))
					}
				case reflect.Struct:
					if field.Type == reflect.TypeOf(time.Time{}) {
						if str, ok := val.(string); ok {
							if parsed, err := time.Parse(time.RFC3339, str); err == nil {
								fieldValue.Set(reflect.ValueOf(parsed))
							}
						}
					}
				}
			}
		}
	}

	now := time.Now()
	event.UpdatedAt = &now

	if err := event.UpdateEvent(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update event", "error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfully!", "event": event})
}

func deleteEvent(context *gin.Context) {
	// This function will handle deleting an event
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
	if event == nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found!"})
		return
	}

	if context.GetInt64("userId") != event.Organizer {
		context.JSON(http.StatusForbidden, gin.H{"message": "You are not authorized to delete this event!"})
		return
	}

	now := time.Now()
	event.DeletedAt = &now // Set DeletedAt to current time
	err = event.DeleteEvent()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete event from the database!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully!", "event": event})
}

