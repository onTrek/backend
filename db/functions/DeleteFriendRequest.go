package functions

import "database/sql"

func DeleteFriendRequest(db *sql.DB, userID, friendID string) error {
	// Prepare the SQL statement to delete the friend request
	stmt, err := db.Prepare("DELETE FROM friends WHERE user_id1 = ? AND user_id2 = ? AND pending = TRUE")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(friendID, userID)
	if err != nil {
		return err
	}

	return nil
}
