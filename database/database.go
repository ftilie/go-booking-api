package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "booking.db")

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME,
		deleted_at DATETIME
	);`
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		location TEXT,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		organizer INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME,
		deleted_at DATETIME,
		FOREIGN KEY (organizer) REFERENCES users(id)
	);`
	createAttendeeTable := `
	CREATE TABLE IF NOT EXISTS event_attendees (
		event_id INTEGER,
		user_id INTEGER,
		PRIMARY KEY (event_id, user_id),
		FOREIGN KEY (event_id) REFERENCES events(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, usersTableErr := DB.Exec(createUsersTable)
	if usersTableErr != nil {
		panic("Failed to create users table: " + usersTableErr.Error())
	}
	_, eventsTableErr := DB.Exec(createEventsTable)
	if eventsTableErr != nil {
		panic("Failed to create events table: " + eventsTableErr.Error())
	}
	_, attendeesTableErr := DB.Exec(createAttendeeTable)
	if attendeesTableErr != nil {
		panic("Failed to create event_attendees table: " + attendeesTableErr.Error())
	}
}
