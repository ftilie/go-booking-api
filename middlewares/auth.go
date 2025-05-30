package middlewares

import (
	"net/http"

	"github.com/ftilie/go-booking-api/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	// This middleware function will check if the user is authenticated
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
		return
	}

	userId, err := utils.VerifyToken(token) // Verify the token to ensure the user is authenticated
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token!"})
		return
	}

	context.Set("userId", userId) // Store the user ID in the context for later use
	context.Next()                // If the token is valid, proceed to the next handler
}
