package db

import "database/sql"

func DeleteActivity(db *sql.DB, activityID int) error {
	// Prepare the SQL statement to delete the activity
	stmt, err := db.Prepare("DELETE FROM activities WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(activityID)
	if err != nil {
		return err
	}

	return nil

}
