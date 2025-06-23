package functions

import (
	"database/sql"
	"fmt"
)

func DeleteUser(db *sql.DB, userID string) error {

	// Enable foreign key enforcement
	_, err := db.Exec("PRAGMA foreign_keys = ON") // Enable foreign key enforcement
	if err != nil {
		return fmt.Errorf("error enabling foreign key enforcement: %v", err)
	}

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
