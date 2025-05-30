package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "dummy_secret_key" // This should be injected in a CI pipeline

func GenerateToken(userId int64, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 2).Unix(), // Token valid for 2 hours
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int64, error) {
	// This function will verify the JWT token
	extractedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, errors.New("could not parse token")
	}

	if !extractedToken.Valid {
		return 0, errors.New("token is not valid")
	}

	claims, ok := extractedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("token claims are not valid")
	}
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return 0, errors.New("token has expired")
		}
	}

	// email := claims["email"].(string)
	userId := int64(claims["userId"].(float64))

	return userId, nil
}
