package routes

import (
	"net/http"
	"time"

	"github.com/ftilie/go-booking-api/models"
	"github.com/ftilie/go-booking-api/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	// This function will handle user signup
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input could not be parsed!"})
		return
	}

	now := time.Now()
	user.CreatedAt = &now

	err = user.CreateUser()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user in the database!"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}

func login(context *gin.Context) {
	// This function will handle user login
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input could not be parsed!"})
		return
	}

	isAuthenticated, err := user.Authenticate()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed!"})
		return
	}
	if !isAuthenticated {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password!"})
		return
	}

	token, err := utils.GenerateToken(user.Id, user.Email) // Generate a token for the user
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User logged in successfully!", "token": token})
}
