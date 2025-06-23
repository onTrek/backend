package functions

import "database/sql"

func GetLeaderByGroup(db *sql.DB, groupId int) (string, error) {
	var leaderId string
	query := "SELECT created_by FROM groups WHERE id = ?"
	err := db.QueryRow(query, groupId).Scan(&leaderId)
	if err != nil {
		return "", err
	}
	return leaderId, nil
}
