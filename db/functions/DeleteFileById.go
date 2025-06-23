package functions

import (
	"OnTrek/utils"
	"database/sql"
	"fmt"
)

func DeleteFileById(db *sql.DB, fileID int, userID string, gpx utils.Gpx) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
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
		return fmt.Errorf("error enabling foreign key enforcement: %v", err)
	}

	// Prepare the SQL statement to delete the file
	stmt, err := tx.Prepare("DELETE FROM gpx_files WHERE id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(fileID, userID)
	if err != nil {
		return err
	}

	// Delete files from the disk
	err = utils.DeleteFiles(gpx)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
