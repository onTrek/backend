package functions

import (
	"database/sql"
)

func CheckGroupExistsByIdAndUserId(db *sql.DB, groupId int, userId string) (bool, error) {
	var exists int
	query := "SELECT EXISTS(SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)"
	err := db.QueryRow(query, userId, groupId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func CheckGroupExistsById(db *sql.DB, groupId int) (bool, error) {
	var exists int
	query := "SELECT EXISTS(SELECT 1 FROM groups WHERE id = ?)"
	err := db.QueryRow(query, groupId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
