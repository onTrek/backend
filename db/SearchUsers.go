package db

import (
	"OnTrek/utils"
	"database/sql"
)

func SearchUsers(db *sql.DB, query string, userId string) ([]utils.UserEssentials, error) {
	// Prepare the SQL query to search for users by name or email
	sqlQuery := `
		SELECT id, username
		FROM users
		WHERE LOWER(username) LIKE LOWER('%' || ? || '%') AND id != ?
		LIMIT 100;`

	rows, err := db.Query(sqlQuery, query, userId)
	if err != nil {
		return []utils.UserEssentials{}, err
	}
	defer rows.Close()

	var users []utils.UserEssentials
	for rows.Next() {
		var user utils.UserEssentials
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
