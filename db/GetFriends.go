package db

import (
	"OnTrek/utils"
	"database/sql"
)

func GetFriends(db *sql.DB, userID string) ([]utils.UserEssentials, error) {
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT id, username FROM users WHERE id IN (SELECT user_id1 FROM friends WHERE user_id2 = ? UNION SELECT user_id2 FROM friends WHERE user_id1 = ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the statement
	rows, err := stmt.Query(userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the friends
	var friends []utils.UserEssentials

	// Iterate through the results
	for rows.Next() {
		var friend utils.UserEssentials
		if err := rows.Scan(&friend.ID, &friend.Username); err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}
