package db

import (
	"OnTrek/utils"
	"database/sql"
	"time"
)

func JoinSession(db *sql.DB, userId string, info utils.SessionInfo) error {

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO session_members (session_id, user_id, latitude, longitude, altitude, accuracy, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement with the provided parameters
	_, err = stmt.Exec(info.SessionID, userId, info.Latitude, info.Longitude, info.Altitude, info.Accuracy, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}

	return nil
}
