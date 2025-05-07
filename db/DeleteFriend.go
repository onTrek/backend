package db

import "database/sql"

func DeleteFriend(db *sql.DB, userID, friendID string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("DELETE FROM friends WHERE user_id1 = ? AND user_id2 = ? OR user_id1 = ? AND user_id2 = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	res, err := stmt.Exec(userID, friendID, friendID, userID)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
