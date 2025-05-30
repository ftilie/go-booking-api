package models

import (
	"time"

	"github.com/ftilie/go-booking-api/database"
)

type Event struct {
	Id          int64
	Title       string `binding:"required"`
	Description string
	Location    string
	StartTime   time.Time `binding:"required"`
	EndTime     time.Time `binding:"required"`
	Organizer   int64
	Attendees   []int64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time // Nullable field for soft delete
}

func (e *Event) CreateEvent() error {
	// Save the event to the database
	eventQuery := `
	INSERT INTO events (title, description, location, start_time, end_time, organizer, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	eventStmt, err := database.DB.Prepare(eventQuery)
	if err != nil {
		return err
	}
	defer eventStmt.Close()
	eventResult, err := eventStmt.Exec(e.Title, e.Description, e.Location, e.StartTime, e.EndTime, e.Organizer, e.CreatedAt)
	if err != nil {
		return err
	}

	id, err := eventResult.LastInsertId()
	if err != nil {
		return err
	}

	e.Id = id
	return nil
}

func getAttendees(eventId int64) ([]int64, error) {
	attendeesQuery := `SELECT user_id FROM event_attendees WHERE event_id = ?`
	attendeesStmt, err := database.DB.Prepare(attendeesQuery)
	if err != nil {
		return nil, err
	}
	defer attendeesStmt.Close()
	// Fetch attendees for the event
	var attendees []int64
	attendeesRows, err := attendeesStmt.Query(eventId)
	if err != nil {
		return nil, err
	}
	defer attendeesRows.Close()
	for attendeesRows.Next() {
		var attendee int64
		err := attendeesRows.Scan(&attendee)
		if err != nil {
			return nil, err
		}
		attendees = append(attendees, attendee)
	}
	return attendees, nil
}

func GetEvents() ([]Event, error) {
	eventsQuery := `SELECT * FROM events WHERE deleted_at IS NULL`
	eventsStmt, err := database.DB.Prepare(eventsQuery)
	if err != nil {
		return nil, err
	}
	defer eventsStmt.Close()
	eventsRows, err := eventsStmt.Query()
	if err != nil {
		return nil, err
	}
	defer eventsRows.Close()

	var events []Event
	for eventsRows.Next() {
		var event Event
		err := eventsRows.Scan(&event.Id, &event.Title, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.Organizer, &event.CreatedAt, &event.UpdatedAt, &event.DeletedAt)
		if err != nil {
			return nil, err
		}

		attendees, err := getAttendees(event.Id)
		if err != nil {
			return nil, err
		}
		// Assign attendees to the event
		event.Attendees = attendees
		events = append(events, event)
	}
	return events, nil
}

func GetEvent(eventId int64) (*Event, error) {
	query := `SELECT * FROM events WHERE id = ? AND deleted_at IS NULL`
	row := database.DB.QueryRow(query, eventId)
	var event Event
	err := row.Scan(&event.Id, &event.Title, &event.Description, &event.Location, &event.StartTime, &event.EndTime, &event.Organizer, &event.CreatedAt, &event.UpdatedAt, &event.DeletedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // Event not found
		}
		return nil, err
	}
	attendees, err := getAttendees(eventId)
	if err != nil {
		return nil, err
	}
	event.Attendees = attendees
	return &event, nil
}

func (e *Event) UpdateEvent() error {
	// Update the event in the database
	eventQuery := `
	UPDATE events
	SET title = ?,
		description = ?,
		location = ?,
		start_time = ?,
		end_time = ?,
		created_at = ?,
		updated_at = ?
	WHERE id = ?`
	eventStmt, err := database.DB.Prepare(eventQuery)
	if err != nil {
		return err
	}
	defer eventStmt.Close()
	_, err = eventStmt.Exec(e.Title, e.Description, e.Location, e.StartTime, e.EndTime, e.CreatedAt, e.UpdatedAt, e.Id)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) DeleteEvent() error {
	// Update the event in the database
	eventQuery := `
	UPDATE events
	SET	deleted_at = ?
	WHERE id = ?`
	eventStmt, err := database.DB.Prepare(eventQuery)
	if err != nil {
		return err
	}
	defer eventStmt.Close()
	_, err = eventStmt.Exec(e.DeletedAt, e.Id)
	if err != nil {
		return err
	}
	return nil
}

func (e Event) RegisterForEvent(userId int64) error {
	// Logic to register the user for the event
	attendeeQuery := `
	INSERT INTO event_attendees (event_id, user_id)
	VALUES (?, ?)`
	attendeeStmt, err := database.DB.Prepare(attendeeQuery)
	if err != nil {
		return err
	}
	defer attendeeStmt.Close()
	_, err = attendeeStmt.Exec(e.Id, userId)
	if err != nil {
		return err
	}
	return nil
}

func (e Event) CancelRegistration(userId int64) error {
	// Logic to cancel the user's registration for the event
	attendeeQuery := `
	DELETE FROM event_attendees WHERE event_id = ? AND user_id = ?`
	attendeeStmt, err := database.DB.Prepare(attendeeQuery)
	if err != nil {
		return err
	}
	defer attendeeStmt.Close()
	_, err = attendeeStmt.Exec(e.Id, userId)
	if err != nil {
		return err
	}
	return nil
}
