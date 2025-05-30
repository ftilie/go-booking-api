package models

import (
	"errors"
	"time"

	"github.com/ftilie/go-booking-api/database"
	"github.com/ftilie/go-booking-api/utils"
)

type User struct {
	Id        int64
	Email     string `binding:"required,email"`
	Password  string `binding:"required"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time // Nullable field for soft delete
}

func (u *User) CreateUser() error {
	// Save the user to the database
	userQuery := `
	INSERT INTO users (email, password, created_at)
	VALUES (?, ?, ?)`
	userStmt, err := database.DB.Prepare(userQuery)
	if err != nil {
		return err
	}
	defer userStmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}
	result, err := userStmt.Exec(u.Email, hashedPassword, u.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.Id = id
	return nil
}

func (u *User) Authenticate() (bool, error) {
	if u.Email == "" {
		return false, errors.New("email must be provided")
	}

	// Query for the user using username or email
	query := `
		SELECT id, password FROM users WHERE email = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var storedPassword string
	err = stmt.QueryRow(u.Email).Scan(&u.Id, &storedPassword)
	if err != nil {
		return false, err
	}

	// Compare hashed password
	isValid := utils.CheckPasswordHash(u.Password, storedPassword)
	return isValid, nil
}
