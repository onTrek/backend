package db

import "database/sql"

func DeleteFileById(db *sql.DB, fileID int, userID string) error {
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
