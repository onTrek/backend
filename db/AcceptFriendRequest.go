package db

import "database/sql"

func AcceptFriendRequest(db *sql.DB, userID, friendID string) error {
	// Prepare the SQL statement to accept a friend request
	query := `
		UPDATE friends
		SET pending = FALSE
		WHERE user_id1 = ? AND user_id2 = ? AND pending = TRUE
	`

	// Execute the SQL statement
	result, err := db.Exec(query, friendID, userID)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // No pending friend request found
	}

	return nil
}
