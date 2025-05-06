package db

import "database/sql"

func DeleteUser(db *sql.DB, userID string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(userID)
	if err != nil {
		return err
	}

	return nil
}
