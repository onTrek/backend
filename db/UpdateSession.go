package db

import (
	"OnTrek/utils"
	"database/sql"
	"time"
)

func UpdateSession(db *sql.DB, userId string, session utils.SessionInfo) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("UPDATE session_members SET latitude = ?, longitude = ?, altitude = ?, accuracy = ?, help_request = ?, going_to = ?, timestamp = ? WHERE session_id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement with the provided parameters
	_, err = stmt.Exec(session.Latitude, session.Longitude, session.Altitude, session.Accuracy, session.HelpRequest, session.GoingTo, time.Now().Format(time.RFC3339), session.SessionID, userId)
	if err != nil {
		return err
	}

	return nil
}
