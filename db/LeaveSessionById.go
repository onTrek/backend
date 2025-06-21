package db

import "database/sql"

func LeaveSessionById(db *sql.DB, userId string, sessionId int) error {
	// Prepare the SQL statement to leave a session
	stmt, err := db.Prepare("DELETE FROM session_members WHERE user_id = ? AND session_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement with the provided user ID and session ID
	_, err = stmt.Exec(userId, sessionId)
	if err != nil {
		return err
	}

	return nil
}
