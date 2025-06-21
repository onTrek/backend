package db

import (
	"database/sql"
	"fmt"
	"time"
)

func JoinSession(db *sql.DB, userId string, sessionId int) error {

	// Enable foreign key enforcement
	_, err := db.Exec("PRAGMA foreign_keys = ON") // Enable foreign key enforcement
	if err != nil {
		return fmt.Errorf("error enabling foreign key enforcement: %v", err)
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO session_members (session_id, user_id, timestamp) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement with the provided parameters
	_, err = stmt.Exec(sessionId, userId, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}

	return nil
}
