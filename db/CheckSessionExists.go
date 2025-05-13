package db

import (
	"OnTrek/utils"
	"database/sql"
)

func CheckSessionExistsByIdAndUserId(db *sql.DB, sessionId int, userId string) (utils.Session, error) {
	var session utils.Session
	query := "SELECT session_id FROM session_members WHERE user_id = ? AND session_id = ?"
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
	query := "SELECT id FROM sessions WHERE id = ? AND closed_at IS NULL"
	err := db.QueryRow(query, sessionId).Scan(&session.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			session.ID = -1     // Session does not exist
			return session, nil // Session does not exist
		}
		return session, err // Some other error occurred
	}

	return session, nil // Session exists
}
