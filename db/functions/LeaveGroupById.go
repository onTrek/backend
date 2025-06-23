package functions

import "database/sql"

func LeaveGroupById(db *sql.DB, userId string, groupId int) error {
	// Prepare the SQL statement to leave a group
	stmt, err := db.Prepare("DELETE FROM group_members WHERE user_id = ? AND group_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement with the provided user ID and group ID
	_, err = stmt.Exec(userId, groupId)
	if err != nil {
		return err
	}

	return nil
}
