package db

import (
	"database/sql"
	"fmt"
)

func DeleteFileById(db *sql.DB, fileID int, userID string) error {

	// Enable foreign key enforcement
	_, err := db.Exec("PRAGMA foreign_keys = ON") // Enable foreign key enforcement
	if err != nil {
		return fmt.Errorf("error enabling foreign key enforcement: %v", err)
	}

	// Prepare the SQL statement to delete the file
	stmt, err := db.Prepare("DELETE FROM gpx_files WHERE id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(fileID, userID)
	if err != nil {
		return err
	}

	return nil
}
