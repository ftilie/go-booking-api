package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// This function will hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPasswordHash(password, hash string) bool {
	// This function will check the password against the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
