package db

import "database/sql"

func GetLeaderBySession(db *sql.DB, sessionId int) (string, error) {
	var leaderId string
	query := "SELECT created_by FROM sessions WHERE id = ?"
	err := db.QueryRow(query, sessionId).Scan(&leaderId)
	if err != nil {
		return "", err
	}
	return leaderId, nil
}
