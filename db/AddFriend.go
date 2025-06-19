package db

import (
	"database/sql"
	"fmt"
)

func AddFriend(db *sql.DB, userID string, friendID string) error {

	// Enable foreign key enforcement
	_, err := db.Exec("PRAGMA foreign_keys = ON") // Enable foreign key enforcement
	if err != nil {
		return fmt.Errorf("error enabling foreign key enforcement: %v", err)
	}

	// Check if the friend already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM friends WHERE user_id1 = ? AND user_id2 = ? OR user_id1 = ? AND user_id2 = ?", userID, friendID, friendID, userID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("Users are already friends") // Friend already exists
	}

	// Add the friend to the database
	_, err = db.Exec("INSERT INTO friends (user_id1, user_id2) VALUES (?, ?)", userID, friendID)
	if err != nil {
		return err
	}

	return nil
}
