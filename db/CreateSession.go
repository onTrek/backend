package db

import (
	"OnTrek/utils"
	"database/sql"
	"time"
)

func CreateSession(db *sql.DB, userID string, info utils.SessionInfo) (utils.Session, error) {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return utils.Session{}, err
	}

	// Rollback on error unless committed
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Insert into sessions
	now := time.Now().Format(time.RFC3339)
	stmt1, err := tx.Prepare("INSERT INTO sessions (created_by, description, created_at) VALUES (?, ?, ?)")
	if err != nil {
		return utils.Session{}, err
	}
	defer stmt1.Close()

	res, err := stmt1.Exec(userID, info.Description, now)
	if err != nil {
		return utils.Session{}, err
	}

	sessionID, err := res.LastInsertId()
	if err != nil {
		return utils.Session{}, err
	}

	// Insert into session_members
	stmt2, err := tx.Prepare("INSERT INTO session_members (session_id, user_id, latitude, longitude, altitude, accuracy, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return utils.Session{}, err
	}
	defer stmt2.Close()

	_, err = stmt2.Exec(sessionID, userID, info.Latitude, info.Longitude, info.Altitude, info.Accuracy, now)
	if err != nil {
		return utils.Session{}, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return utils.Session{}, err
	}

	return utils.Session{ID: int(sessionID), CreatedBy: userID}, nil
}
