package functions

import (
	"database/sql"
	"fmt"
)

func DeleteGroupById(db *sql.DB, userId string, groupId int) error {

	// Enable foreign key enforcement
	_, err := db.Exec("PRAGMA foreign_keys = ON") // Enable foreign key enforcement
	if err != nil {
		return fmt.Errorf("error enabling foreign key enforcement: %v", err)
	}

	// Prepare the SQL statement to delete the group
	stmt, err := db.Prepare("DELETE FROM groups WHERE id = ? AND created_by = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement with the provided group ID and user ID
	_, err = stmt.Exec(groupId, userId)
	if err != nil {
		return err
	}

	return nil
}
