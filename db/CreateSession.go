package db

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
	"time"
)

func CreateSession(db *sql.DB, user utils.User, info utils.SessionInfo) (utils.Session, error) {
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

	// Enable foreign key enforcement
	_, err = tx.Exec("PRAGMA foreign_keys = ON") // Enable foreign key enforcement
	if err != nil {
		return utils.Session{}, fmt.Errorf("error enabling foreign key enforcement: %v", err)
	}

	// Check if the file exists
	var fileExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM gpx_files WHERE id = ? AND user_id = ?)", info.FileId, user.ID).Scan(&fileExists)
	if err != nil {
		return utils.Session{}, err
	}

	if !fileExists {
		return utils.Session{}, fmt.Errorf("file with ID %d does not exist for user %s", info.FileId, user.Username)
	}

	// Insert into sessions
	now := time.Now().Format(time.RFC3339)
	stmt1, err := tx.Prepare("INSERT INTO sessions (created_by, description, created_at, file_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return utils.Session{}, err
	}
	defer stmt1.Close()

	res, err := stmt1.Exec(user.ID, info.Description, now, info.FileId)
	if err != nil {
		return utils.Session{}, err
	}

	sessionID, err := res.LastInsertId()
	if err != nil {
		return utils.Session{}, err
	}

	// Insert into session_members
	err = JoinSession(tx, user.ID, int(sessionID))

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return utils.Session{}, err
	}

	return utils.Session{ID: int(sessionID), CreatedBy: user.ID}, nil
}
