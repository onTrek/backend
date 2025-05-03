package db

import "database/sql"

func GetUserIdByActivity(db *sql.DB, activityID int) (string, error) {
	var userID string
	query := "SELECT user_id FROM activities WHERE id = ?"
	err := db.QueryRow(query, activityID).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}
