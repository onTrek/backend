package db

import (
	"OnTrek/utils"
	"database/sql"
)

func CheckSessionExistsByIdAndUserId(db *sql.DB, sessionId int, userId string) (utils.Session, error) {
	var session utils.Session
	query := "SELECT session_id FROM sessions_members WHERE user_id = ? AND session_id = ?"
	err := db.QueryRow(query, userId, sessionId).Scan(&session.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return session, nil // Session does not exist
		}
		return session, err // Some other error occurred
	}
	return session, nil // Session exists
}

func CheckSessionExistsById(db *sql.DB, sessionId int) (utils.Session, error) {
	var session utils.Session
	query := "SELECT EXISTS(SELECT 1 FROM sessions WHERE session_id = ?)"
	err := db.QueryRow(query, sessionId).Scan(&session.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return session, nil // Session does not exist
		}
		return session, err // Some other error occurred
	}
	return session, nil // Session exists
}
