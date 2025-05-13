package db

import (
	"OnTrek/utils"
	"database/sql"
)

func GetSessionsByUserId(db *sql.DB, userId string) ([]utils.Session, error) {

	var session []utils.Session
	// Get the session for the user
	query := `SELECT s.id, s.created_by, s.created_at, s.closed_at FROM sessions s JOIN session_members sm ON s.id = sm.session_id WHERE sm.user_id = ?`
	rows, err := db.Query(query, userId)
	if err != nil {
		return []utils.Session{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var s utils.Session
		err := rows.Scan(&s.ID, &s.CreatedBy, &s.CreatedAt, &s.ClosedAt)
		if err != nil {
			return []utils.Session{}, err
		}
		session = append(session, s)
	}

	if err := rows.Err(); err != nil {
		return []utils.Session{}, err
	}
	if len(session) == 0 {
		return []utils.Session{}, sql.ErrNoRows
	} else {
		return session, nil
	}
}
